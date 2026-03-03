package main

import (
	"flag"
	"log"
	"nodestore/app"
	"nodestore/app/config"
	"nodestore/app/db"

	"github.com/hyahm/golog"
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

func main() {
	port := flag.Int("p", 8080, "服务端口")
	// 字符串型：-env，默认 dev，说明 "运行环境"
	configfile := flag.String("c", "config.yaml", "配置文件")

	// 2. 解析命令行参数（必须调用，否则参数值为默认值）
	flag.Parse()
	defer golog.Sync()
	err := config.InitConfig(*configfile)
	if err != nil {
		log.Fatal(err)
	}
	err = db.InitDB()
	if err != nil {
		golog.Error("数据库初始化失败:", err)
		return
	}

	app.Run(*port)
}
