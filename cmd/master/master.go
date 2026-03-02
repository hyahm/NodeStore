package main

import (
	"nodestore/app"
)

// ==================== 全局配置 ====================

type Share struct {
	ID         int64  `json:"id"`
	UserID     int64  `json:"user_id"`
	FileID     int64  `json:"file_id"`
	ShareKey   string `json:"share_key"`
	ExpireAt   string `json:"expire_at,omitempty"`
	CreateTime string `json:"create_time"`
}

// ==================== 核心改造1：按MD5下载文件 ====================

// ==================== 主函数（新增路由） ====================
func main() {

	app.Run()
}
