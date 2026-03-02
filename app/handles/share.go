package handles

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"nodestore/app/common"
	"nodestore/app/db"
	"nodestore/app/module"
	"nodestore/app/response"
	"sort"
	"time"

	"github.com/hyahm/golog"
)

type Share struct{}

type ShareResponse struct {
	ShareKey  string `json:"share_key"`
	ExpireAt  string `json:"expire_at"`
	ShareLink string `json:"share_link"`
}

// ==================== 分享文件（适配MD5和软删除） ====================
func (Share) ShareCreateHandler(w http.ResponseWriter, r *http.Request) {
	claims, err := module.GetCurrentUser(r)
	if err != nil {
		http.Error(w, err.Error(), 401)
		return
	}

	// 支持按ID或MD5创建分享
	fileID := r.URL.Query().Get("file_id")
	fileMD5 := r.URL.Query().Get("md5")
	days := r.URL.Query().Get("days")
	if days == "" {
		days = "7"
	}

	var fid int64
	if fileMD5 != "" {
		// 按MD5查文件ID
		err = db.DB.Raw(`
		SELECT id FROM file_metadata 
		WHERE user_id=? AND md5=? AND is_deleted=0`,
			claims.UserID, fileMD5).Scan(&fid).Error
		if err != nil || fid == 0 {
			http.Error(w, "file not found by md5", 403)
			return
		}
	} else {
		fmt.Sscan(fileID, &fid)
		// 验证文件存在且未删除
		var cnt int
		err = db.DB.Raw(`SELECT COUNT(*) FROM file_metadata WHERE id=? AND user_id=? AND is_deleted=0`,
			fid, claims.UserID).Scan(&cnt).Error
		if err != nil || cnt == 0 {
			http.Error(w, "file not found or deleted", 403)
			return
		}
	}

	shareKey := common.RandString(16)
	now := time.Now()
	exp := now.Add(7 * 24 * time.Hour).Format(time.DateTime)
	if days == "0" {
		exp = ""
	}

	err = db.DB.Exec(`
	INSERT INTO shares (user_id, file_id, share_key, expire_at, create_time)
	VALUES (?,?,?,?,?)`,
		claims.UserID, fid, shareKey, exp, now.Format(time.DateTime)).Error
	if err != nil {
		http.Error(w, err.Error(), 401)
		return
	}
	sr := ShareResponse{
		ShareKey:  shareKey,
		ExpireAt:  exp,
		ShareLink: fmt.Sprintf("/share/download?key=%s", shareKey),
	}
	response.Success(r, "ok", sr)
	// json.NewEncoder(w).Encode(map[string]interface{}{
	// 	"share_key":  shareKey,
	// 	"expire_at":  exp,
	// 	"share_link": fmt.Sprintf("/share/download?key=%s", shareKey),
	// })
}

func (Share) ShareDownloadHandler(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	if key == "" {
		http.Error(w, "key empty", 400)
		return
	}

	// 查分享
	type ShareDB struct {
		UserID   int64  `json:"user_id"`
		FileID   int64  `json:"file_id"`
		ExpireAt string `json:"expire_at"`
	}
	var s ShareDB
	err := db.DB.Raw(`
	SELECT user_id, file_id, expire_at FROM shares WHERE share_key=?`, key).Scan(&s).Error
	if err != nil {
		http.Error(w, "share not found", 404)
		return
	}

	// 过期判断
	if s.ExpireAt != "" {
		exp, _ := time.Parse(time.DateTime, s.ExpireAt)
		if time.Now().After(exp) {
			http.Error(w, "share expired", 410)
			return
		}
	}

	// 查文件（过滤软删除）
	var fm struct {
		realName  string
		contentID int64
	}
	err = db.DB.Raw(`
	SELECT real_name, content_id
	FROM file_metadata WHERE id=? AND is_deleted=0`, s.FileID).Scan(&fm).Error
	if err != nil {
		http.Error(w, "file not found or deleted", 404)
		return
	}

	// 查分片
	var chunks []*FileChunk
	err = db.DB.Raw(`
	SELECT chunk_id, node, is_data, idx, length 
	FROM file_chunks 
	WHERE content_id=?`, fm.contentID).Scan(&chunks).Error
	if err != nil || len(chunks) == 0 {
		http.Error(w, "file lost", 500)
		return
	}

	// 恢复文件
	fullFileData, err := restoreFileByChunks(chunks)
	if err != nil {
		http.Error(w, "restore file failed: "+err.Error(), 500)
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename="+url.PathEscape(fm.realName))
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(fullFileData)
}

// restoreFileByChunks 从分片恢复完整文件（复用原有逻辑，封装为函数）
func restoreFileByChunks(chunks []*FileChunk) ([]byte, error) {
	good := make(map[int][]byte)
	lost := []int{}
	for _, c := range chunks {
		if !c.IsData {
			continue
		}
		b, err := fetch(c.Node, c.ChunkID)
		if err == nil && len(b) == c.Length {
			good[c.Index] = b
		} else {
			lost = append(lost, c.Index)
			golog.Error("读取数据块失败:", c.ChunkID, err)
		}
	}

	// 容错恢复
	if len(lost) > 0 && len(lost) <= common.M {
		parity := make(map[int][]byte)
		for _, c := range chunks {
			if c.IsData {
				continue
			}
			b, err := fetch(c.Node, c.ChunkID)
			if err == nil && len(b) == c.Length {
				parity[c.Index] = b
			} else {
				golog.Error("读取校验块失败:", c.ChunkID, err)
			}
		}
		if p0, ok := parity[0]; ok {
			for _, idx := range lost {
				fix := make([]byte, len(p0))
				for i := range fix {
					fix[i] = p0[i]
					for j, b := range good {
						if j == idx {
							continue
						}
						if i < len(b) {
							fix[i] ^= b[i]
						}
					}
				}
				good[idx] = fix
			}
		}
	}

	// 拼接文件
	keys := make([]int, 0, len(good))
	for k := range good {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	var full []byte
	for _, k := range keys {
		full = append(full, good[k]...)
	}

	return full, nil
}

// fetch 从节点获取分片
func fetch(node, cid string) ([]byte, error) {
	esc := url.PathEscape(cid)
	u := fmt.Sprintf("%s/chunk?cid=%s", node, esc)
	resp, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}
