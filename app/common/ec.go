package common

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"path/filepath"
	"strings"
)

const (
	K = 3
	M = 2
)

// ==================== 媒体类型映射（用于在线播放） ====================
var mediaTypeMap = map[string]string{
	".mp3":  "audio/mpeg",
	".mp4":  "video/mp4",
	".avi":  "video/x-msvideo",
	".mov":  "video/quicktime",
	".m4a":  "audio/mp4",
	".wav":  "audio/wav",
	".flac": "audio/flac",
	".webm": "video/webm",
}

// ==================== 工具函数 ====================
// EncodeFilename 编码文件名
func EncodeFilename(filename string) string {
	return base64.URLEncoding.EncodeToString([]byte(filename))
}

// DecodeFilename 解码文件名
func DecodeFilename(encoded string) (string, error) {
	data, err := base64.URLEncoding.DecodeString(encoded)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// SafeChunkID 生成唯一分片ID（关联content_id）
func SafeChunkID(contentID int64, isData bool, idx int) string {
	prefix := "d"
	if !isData {
		prefix = "p"
	}
	return fmt.Sprintf("content_%d_%s_%d", contentID, prefix, idx)
}

// GetFileMD5 计算文件MD5
func GetFileMD5(data []byte) string {
	h := md5.New()
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}

// RandString 生成随机字符串
func RandString(n int) string {
	b := make([]byte, n)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)[:n]
}

// GetFileExt 获取文件扩展名（小写）
func GetFileExt(filename string) string {
	return strings.ToLower(filepath.Ext(filename))
}

// GetMediaType 根据扩展名获取媒体类型
func GetMediaType(filename string) string {
	ext := GetFileExt(filename)
	if t, ok := mediaTypeMap[ext]; ok {
		return t
	}
	return "application/octet-stream"
}

// ==================== 全局变量 ====================
