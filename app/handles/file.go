package handles

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"nodestore/app/common"
	"nodestore/app/db"
	"nodestore/app/module"
	"nodestore/app/response"
	"strconv"
	"strings"
	"time"

	"github.com/hyahm/golog"
	"gorm.io/gorm"
)

// ==================== 结构体 ====================
// FileContent 全局文件内容表（存储唯一的文件分片数据）
type FileContent struct {
	ID        int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	MD5       string `gorm:"column:md5;uniqueIndex" json:"md5"`   // 文件MD5（唯一索引）
	ChunkMeta string `gorm:"column:chunk_meta" json:"chunk_meta"` // 分片元数据JSON字符串
	FileSize  int64  `gorm:"column:file_size" json:"file_size"`   // 文件大小
	CreateAt  string `gorm:"column:create_at" json:"create_at"`   // 创建时间
}

// FileChunk 分片表（关联FileContent）
type FileChunk struct {
	ID        int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	ContentID int64  `gorm:"column:content_id;index" json:"content_id"` // 关联FileContent.ID
	Name      string `gorm:"column:name" json:"name"`
	ChunkID   string `gorm:"column:chunk_id" json:"cid"`
	Node      string `gorm:"column:node" json:"node"`
	IsData    bool   `gorm:"column:is_data" json:"data"`
	Index     int    `gorm:"column:idx" json:"idx"`
	Length    int    `gorm:"column:length" json:"len"`
	RealName  string `gorm:"column:real_name" json:"real"`
	CreateAt  string `gorm:"column:create_at" json:"create_at"`
}

func (f *FileChunk) TableName() string {
	return "file_chunks"
}

func (f *FileContent) TableName() string {
	return "file_content"
}

type UploadResponse struct {
	Msg   string `json:"msg"`
	File  string `json:"file"`
	Dir   string `json:"dir"`
	Md5   string `json:"md5"`
	IsNew bool   `json:"is_new"`
}

type File struct{}

// ==================== 文件上传（核心：MD5去重+软链接） ====================
func (File) UploadHandler(w http.ResponseWriter, r *http.Request) {
	golog.Info("开始处理文件上传")
	claims, err := module.GetCurrentUser(r)
	if err != nil {
		http.Error(w, err.Error(), 401)
		return
	}

	// 1. 读取文件
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	defer file.Close()

	// 2. 读取文件内容并计算MD5
	fileData, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "read file failed: "+err.Error(), 400)
		return
	}
	fileMD5 := common.GetFileMD5(fileData)
	fileSize := int64(len(fileData))

	// 3. 获取目录参数
	dir := r.FormValue("dir")
	if dir == "" {
		dir = "/"
	}
	if !strings.HasPrefix(dir, "/") {
		dir = "/" + dir
	}

	// 4. 检查目录是否存在
	// var cnt int
	// err = db.Raw(`SELECT COUNT(*) FROM directories WHERE user_id=? AND path=?`,
	// 	claims.UserID, dir).Scan(&cnt).Error
	// if err != nil || cnt == 0 {
	// 	http.Error(w, "dir not exists", 400)
	// 	return
	// }

	// 5. 查询MD5是否已存在
	var fileContent FileContent
	err = db.DB.Raw(`SELECT id, md5, file_size FROM file_content WHERE md5=?`, fileMD5).Scan(&fileContent).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		http.Error(w, "query file md5 failed: "+err.Error(), 500)
		return
	}

	realName := header.Filename
	encName := common.EncodeFilename(realName)
	now := time.Now().Format(time.DateTime)

	// 6. MD5已存在：创建软链接
	if fileContent.ID > 0 {
		// 检查当前用户是否已存在该文件（避免重复创建软链接）
		var existFile FileInfo
		err = db.DB.Raw(`
		SELECT id FROM file_metadata 
		WHERE user_id=? AND dir=? AND real_name=? AND md5=? AND is_deleted=0`,
			claims.UserID, dir, realName, fileMD5).Scan(&existFile).Error
		if err == nil && existFile.ID > 0 {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"msg":    "file already exists, no need to upload",
				"md5":    fileMD5,
				"file":   realName,
				"dir":    dir,
				"is_new": false,
			})
			return
		}

		// 创建软链接（仅插入元数据，不传输分片）
		err = db.DB.Exec(`
		INSERT INTO file_metadata
		(user_id, real_name, encoded_name, dir, upload_time, file_size, md5, content_id, is_deleted)
		VALUES (?,?,?,?,?,?,?,?,0)`,
			claims.UserID, realName, encName, dir, now, fileContent.FileSize, fileMD5, fileContent.ID).Error
		if err != nil {
			http.Error(w, "create soft link failed: "+err.Error(), 500)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"msg":    "file exists, create soft link success",
			"md5":    fileMD5,
			"file":   realName,
			"dir":    dir,
			"is_new": false,
		})
		return
	}

	// 7. MD5不存在：分片上传并创建全局文件记录
	// 7.1 EC分片
	chunkSize := (len(fileData) + common.K - 1) / common.K
	datas := make([][]byte, common.K)
	for i := 0; i < common.K; i++ {
		s := i * chunkSize
		e := s + chunkSize
		if e > len(fileData) {
			e = len(fileData)
		}
		datas[i] = fileData[s:e]
	}

	parities := make([][]byte, common.M)
	for i := range parities {
		parities[i] = make([]byte, chunkSize)
	}
	for i := 0; i < chunkSize; i++ {
		var xor byte
		for j := 0; j < common.K; j++ {
			if i < len(datas[j]) {
				xor ^= datas[j][i]
			}
		}
		for j := 0; j < common.M; j++ {
			parities[j][i] = xor
		}
	}

	// 7.2 获取节点
	ns, err := pickNodesByUser(claims.UserID, common.K+common.M)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// 7.3 先创建file_content记录（获取content_id）
	chunkMeta := map[string]interface{}{
		"chunk_size": chunkSize,
		"k":          common.K,
		"m":          common.M,
	}
	chunkMetaJSON, _ := json.Marshal(chunkMeta)
	err = db.DB.Exec(`
	INSERT INTO file_content (md5, chunk_meta, file_size, create_at)
	VALUES (?,?,?,?)`, fileMD5, string(chunkMetaJSON), fileSize, now).Error
	if err != nil {
		http.Error(w, "create file content failed: "+err.Error(), 500)
		return
	}

	// 获取新增的content_id
	var contentID int64
	db.DB.Raw(`SELECT id FROM file_content WHERE md5=?`, fileMD5).Scan(&contentID)

	// 7.4 上传数据块
	for i := 0; i < common.K; i++ {
		node := ns[i%len(ns)]
		cid := common.SafeChunkID(contentID, true, i)
		err := uploadChunk(node, cid, datas[i])
		if err != nil {
			golog.Error("上传数据块失败:", err)
			continue
		}
		// 保存分片信息
		chunk := &FileChunk{
			ContentID: contentID,
			Name:      encName,
			ChunkID:   cid,
			Node:      node,
			IsData:    true,
			Index:     i,
			Length:    len(datas[i]),
			RealName:  realName,
			CreateAt:  now,
		}
		db.DB.Create(chunk)
	}

	// 7.5 上传校验块
	for i := 0; i < common.M; i++ {
		node := ns[(common.K+i)%len(ns)]
		cid := common.SafeChunkID(contentID, false, i)
		err := uploadChunk(node, cid, parities[i])
		if err != nil {
			golog.Error("上传校验块失败:", err)
			continue
		}
		// 保存分片信息
		chunk := &FileChunk{
			ContentID: contentID,
			Name:      encName,
			ChunkID:   cid,
			Node:      node,
			IsData:    false,
			Index:     i,
			Length:    len(parities[i]),
			RealName:  realName,
			CreateAt:  now,
		}
		db.DB.Create(chunk)
	}

	// 7.6 创建当前用户的文件元数据
	err = db.DB.Exec(`
	INSERT INTO file_metadata
	(user_id, real_name, encoded_name, dir, upload_time, file_size, md5, content_id, is_deleted)
	VALUES (?,?,?,?,?,?,?,?,0)`,
		claims.UserID, realName, encName, dir, now, fileSize, fileMD5, contentID).Error
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	res := UploadResponse{
		Msg:   "upload ok",
		File:  realName,
		Dir:   dir,
		Md5:   fileMD5,
		IsNew: true,
	}
	w.Header().Set("Content-Type", "application/json")
	response.Success(r, "ok", res)
	// json.NewEncoder(w).Encode(map[string]interface{}{
	// 	"msg":    "upload ok",
	// 	"file":   realName,
	// 	"dir":    dir,
	// 	"md5":    fileMD5,
	// 	"is_new": true,
	// })
}

// uploadChunk 上传分片到节点
func uploadChunk(node, cid string, data []byte) error {
	esc := url.PathEscape(cid)
	u := fmt.Sprintf("%s/chunk?cid=%s", node, esc)
	req, _ := http.NewRequest("PUT", u, bytes.NewReader(data))
	req.Header.Set("Content-Length", fmt.Sprintf("%d", len(data)))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func pickNodesByUser(userID int64, n int) ([]string, error) {
	var ns []string
	err := db.DB.Raw(`SELECT address FROM nodes WHERE user_id=?`, userID).Scan(&ns).Error
	if err != nil {
		return nil, err
	}
	if len(ns) < n {
		return nil, fmt.Errorf("need %d nodes, but only %d available", n, len(ns))
	}
	return ns[:n], nil
}

func (File) DownloadByMD5Handler(w http.ResponseWriter, r *http.Request) {
	claims, err := module.GetCurrentUser(r)
	if err != nil {
		http.Error(w, err.Error(), 401)
		return
	}

	// 1. 获取MD5参数
	fileMD5 := r.URL.Query().Get("md5")
	if fileMD5 == "" {
		http.Error(w, "md5 parameter is required", 400)
		return
	}

	// 2. 查询用户是否有权限（未软删除的元数据）
	var fileInfo FileInfo
	err = db.DB.Raw(`
	SELECT id, content_id, real_name, file_size 
	FROM file_metadata 
	WHERE user_id=? AND md5=? AND is_deleted=0`,
		claims.UserID, fileMD5).Scan(&fileInfo).Error
	if err != nil || fileInfo.ID == 0 {
		http.Error(w, "file not found or no permission", 404)
		return
	}

	// 3. 查询分片信息
	var chunks []*FileChunk
	err = db.DB.Raw(`
	SELECT chunk_id, node, is_data, idx, length 
	FROM file_chunks 
	WHERE content_id=?`, fileInfo.ContentID).Scan(&chunks).Error
	if err != nil {
		http.Error(w, "get chunks failed: "+err.Error(), 500)
		return
	}
	if len(chunks) == 0 {
		http.Error(w, "file chunks lost", 500)
		return
	}

	// 4. 恢复完整文件
	fullFileData, err := restoreFileByChunks(chunks)
	if err != nil {
		http.Error(w, "restore file failed: "+err.Error(), 500)
		return
	}

	// 5. 设置下载响应头
	w.Header().Set("Content-Disposition", "attachment; filename="+url.PathEscape(fileInfo.RealName))
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", strconv.FormatInt(int64(len(fullFileData)), 10))
	w.Write(fullFileData)
}

// ==================== 核心改造2：在线播放接口（支持Range请求） ====================
func (File) StreamHandler(w http.ResponseWriter, r *http.Request) {
	claims, err := module.GetCurrentUser(r)
	if err != nil {
		http.Error(w, err.Error(), 401)
		return
	}

	// 1. 获取参数
	fileMD5 := r.URL.Query().Get("md5")
	if fileMD5 == "" {
		http.Error(w, "md5 parameter is required", 400)
		return
	}

	// 2. 查询用户权限和文件信息
	var fileInfo FileInfo
	err = db.DB.Raw(`
	SELECT id, content_id, real_name, file_size 
	FROM file_metadata 
	WHERE user_id=? AND md5=? AND is_deleted=0`,
		claims.UserID, fileMD5).Scan(&fileInfo).Error
	if err != nil || fileInfo.ID == 0 {
		http.Error(w, "file not found or no permission", 404)
		return
	}

	// 3. 验证是否为支持的媒体类型
	fileExt := common.GetFileExt(fileInfo.RealName)
	mediaType := common.GetMediaType(fileInfo.RealName)
	if mediaType == "application/octet-stream" {
		http.Error(w, "unsupported media type: "+fileExt, 400)
		return
	}

	// 4. 查询分片并恢复完整文件（实际生产中可优化为流式恢复，此处简化）
	var chunks []*FileChunk
	err = db.DB.Raw(`
	SELECT chunk_id, node, is_data, idx, length 
	FROM file_chunks 
	WHERE content_id=?`, fileInfo.ContentID).Scan(&chunks).Error
	if err != nil {
		http.Error(w, "get chunks failed: "+err.Error(), 500)
		return
	}
	if len(chunks) == 0 {
		http.Error(w, "file chunks lost", 500)
		return
	}

	// 5. 恢复完整文件
	fullFileData, err := restoreFileByChunks(chunks)
	if err != nil {
		http.Error(w, "restore file failed: "+err.Error(), 500)
		return
	}
	fileSize := int64(len(fullFileData))

	// 6. 解析Range请求
	rangeHeader := r.Header.Get("Range")
	start, end, err := common.ParseRangeHeader(rangeHeader, fileSize)
	if err != nil {
		http.Error(w, "invalid Range: "+err.Error(), 416)
		return
	}

	// 7. 设置流式播放响应头
	w.Header().Set("Content-Type", mediaType)
	w.Header().Set("Accept-Ranges", "bytes")
	w.Header().Set("Content-Length", strconv.FormatInt(end-start+1, 10))

	// 8. 处理Range请求
	if rangeHeader != "" {
		w.WriteHeader(http.StatusPartialContent)
		w.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, fileSize))
		// 写入指定范围的字节
		w.Write(fullFileData[start : end+1])
	} else {
		// 无Range请求，返回完整文件
		w.Write(fullFileData)
	}
}

// ==================== 文件删除（软删除） ====================
func (File) DeleteFileHandler(w http.ResponseWriter, r *http.Request) {
	claims, err := module.GetCurrentUser(r)
	if err != nil {
		http.Error(w, err.Error(), 401)
		return
	}

	// 支持两种删除方式：按文件名+目录 或 按MD5
	fname := r.URL.Query().Get("filename")
	fileMD5 := r.URL.Query().Get("md5")
	dir := r.URL.Query().Get("dir")
	if dir == "" {
		dir = "/"
	}

	var errMsg string
	var affected int64

	// 按MD5删除
	if fileMD5 != "" {
		result := db.DB.Exec(`
		UPDATE file_metadata 
		SET is_deleted=1 
		WHERE user_id=? AND md5=? AND is_deleted=0`,
			claims.UserID, fileMD5)
		if result.Error != nil {
			errMsg = "delete by md5 failed: " + result.Error.Error()
		} else {
			affected = result.RowsAffected
		}
	} else if fname != "" {
		// 按文件名+目录删除
		result := db.DB.Exec(`
		UPDATE file_metadata 
		SET is_deleted=1 
		WHERE user_id=? AND dir=? AND real_name=? AND is_deleted=0`,
			claims.UserID, dir, fname)
		if result.Error != nil {
			errMsg = "delete by filename failed: " + result.Error.Error()
		} else {
			affected = result.RowsAffected
		}
	} else {
		http.Error(w, "filename or md5 parameter is required", 400)
		return
	}

	if errMsg != "" {
		http.Error(w, errMsg, 500)
		return
	}

	if affected == 0 {
		http.Error(w, "file not found or already deleted", 404)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response.Success(r, fmt.Sprintf("file: %s is delete", fname), nil)
	// json.NewEncoder(w).Encode(map[string]string{
	// 	"msg":  "file deleted (soft delete)",
	// 	"md5":  fileMD5,
	// 	"file": fname,
	// 	"dir":  dir,
	// })
}
