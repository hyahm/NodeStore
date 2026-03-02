package handles

import (
	"net/http"
	"net/url"
	"nodestore/app/db"
	"nodestore/app/module"
	"nodestore/app/response"
	"path"
	"strings"
	"time"

	"github.com/hyahm/golog"
	"github.com/hyahm/xmux"
)

type CreateDir struct {
	Path string `json:"path"`
	Prev string `json:"prev"`
}

func NewCreateDir() *CreateDir {
	return &CreateDir{}
}

// ==================== 目录管理 ====================
func (cd *CreateDir) CreateDirHandler(w http.ResponseWriter, r *http.Request) {
	claims, err := module.GetCurrentUser(r)
	if err != nil {
		http.Error(w, err.Error(), 401)
		return
	}
	cd = xmux.GetInstance(r).Data.(*CreateDir)
	// path := r.URL.Query().Get("path")
	// if path == "" {
	// 	http.Error(w, "path empty", 400)
	// 	return
	// }

	dir := path.Join(cd.Prev, cd.Path)
	// path = filepath.Clean(path)
	if !strings.HasPrefix(dir, "/") {
		dir = "/" + dir
	}

	now := time.Now().Format(time.DateTime)
	err = db.DB.Exec(`INSERT INTO directories (user_id,path,create_time) VALUES (?,?,?)`,
		claims.UserID, dir, now).Error
	if err != nil {
		http.Error(w, err.Error(), 401)
		return
	}
	response.Success(r, "dir created", nil)
	// w.Write([]byte("dir created"))
}

// FileInfo 文件元数据表（用户级，软链接）
type FileInfo struct {
	ID          int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID      int64  `gorm:"column:user_id;index" json:"user_id"`
	RealName    string `gorm:"column:real_name" json:"filename"`
	EncodedName string `gorm:"column:encoded_name" json:"encoded_name"`
	Dir         string `gorm:"column:dir" json:"dir"`
	UploadTime  string `gorm:"column:upload_time" json:"upload_time"`
	FileSize    int64  `gorm:"column:file_size" json:"file_size"`
	MD5         string `gorm:"column:md5;index" json:"md5"`                   // 文件MD5
	ContentID   int64  `gorm:"column:content_id;index" json:"content_id"`     // 关联FileContent.ID
	IsDeleted   int    `gorm:"column:is_deleted;default:0" json:"is_deleted"` // 软删除标记（0-未删，1-已删）
}

// FSListResponse 文件系统列表返回结构体
type FSListResponse struct {
	Dirs   []string   `json:"dirs"`    // 子目录列表
	Files  []FileInfo `json:"files"`   // 文件列表（过滤软删除）
	Dir    string     `json:"dir"`     // 当前查询的目录
	UserID int64      `json:"user_id"` // 当前用户ID
}

// ==================== 文件系统列表（合并目录+文件） ====================
func (cd *CreateDir) FsListHandler(w http.ResponseWriter, r *http.Request) {

	claims, err := module.GetCurrentUser(r)
	if err != nil {
		http.Error(w, err.Error(), 401)
		return
	}

	dir := r.URL.Query().Get("dir")
	if dir == "" {
		dir = "/"
	}
	dir, err = url.PathUnescape(dir)
	if err != nil {
		http.Error(w, err.Error(), 401)
		return
	}
	// dir = filepath.Clean(dir)
	// golog.Info(dir)
	// if !strings.HasPrefix(dir, "/") {
	// 	dir = "/" + dir
	// }
	// 查询子目录
	subDirs := make([]string, 0)
	parentDir := dir
	if parentDir != "/" {
		parentDir += "/"
	}
	err = db.DB.Raw(`
	SELECT path FROM directories 
	WHERE user_id=? AND path LIKE ?`,
		claims.UserID, parentDir+"%").Scan(&subDirs).Error
	if err != nil {
		http.Error(w, "failed to query dirs: "+err.Error(), 500)
		return
	}
	newSubDir := make([]string, 0, len(subDirs))
	for i := range subDirs {
		path := strings.Replace(subDirs[i], dir+"/", "", 1)
		if dir == "/" {
			path = strings.Replace(subDirs[i], dir, "", 1)
		}

		golog.Info(path)
		if path == "" {
			continue
		}
		if !strings.Contains(path, "/") {
			newSubDir = append(newSubDir, subDirs[i])
		}
	}
	// 查询文件（过滤软删除）
	files := make([]FileInfo, 0)
	err = db.DB.Raw(`
	SELECT id, user_id, real_name, encoded_name, dir, upload_time, file_size, md5, content_id, is_deleted
	FROM file_metadata
	WHERE user_id=? AND dir = ? AND is_deleted = 0
	ORDER BY upload_time DESC`, claims.UserID, dir).Scan(&files).Error
	if err != nil {
		http.Error(w, "failed to query files: "+err.Error(), 500)
		return
	}
	golog.Info(newSubDir)
	resp := FSListResponse{
		Dirs:   newSubDir,
		Files:  files,
		Dir:    dir,
		UserID: claims.UserID,
	}

	w.Header().Set("Content-Type", "application/json")
	response.Success(r, "ok", resp)
	// json.NewEncoder(w).Encode(response)
}
