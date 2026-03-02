package db

import (
	"log"
	"nodestore/app/config"
	"sync"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB
	mu sync.Mutex // 数据库操作锁
)

// ==================== 数据库初始化 ====================
func InitDB() error {
	var err error
	DB, err = gorm.Open(mysql.Open(config.Cfg.DB), &gorm.Config{})
	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}

	// 1. 用户表
	DB.Exec(`
	CREATE TABLE IF NOT EXISTS users (
		id bigint PRIMARY KEY auto_increment,
		username varchar(50) NOT NULL default '',
		password varchar(100) NOT NULL default '',
		create_time datetime 
	);`)

	// 2. 节点表
	DB.Exec(`
	CREATE TABLE IF NOT EXISTS nodes (
		id bigint PRIMARY KEY auto_increment,
		user_id INTEGER NOT NULL,
		address varchar(255) not null default '',
		create_time datetime,
		UNIQUE(user_id, address)
	);`)

	// 3. 目录表
	DB.Exec(`
	CREATE TABLE IF NOT EXISTS directories (
		id bigint PRIMARY KEY auto_increment,
		user_id INTEGER NOT NULL,
		path varchar(60) NOT NULL,
		create_time datetime,
		UNIQUE(user_id, path)
	);`)

	// 4. 全局文件内容表（新增）
	DB.Exec(`
	CREATE TABLE IF NOT EXISTS file_content (
		id bigint PRIMARY KEY auto_increment,
		md5 varchar(32) NOT NULL UNIQUE,
		chunk_meta TEXT NOT NULL,
		file_size bigint NOT NULL default 0,
		create_at datetime NOT NULL
	);`)

	// 5. 分片表（修改：关联content_id）
	DB.Exec(`
	CREATE TABLE IF NOT EXISTS file_chunks (
		id bigint PRIMARY KEY auto_increment,
		content_id bigint NOT NULL,
		name varchar(255) NOT NULL default '',
		chunk_id varchar(255) NOT NULL default '',
		node varchar(255) NOT NULL default '',
		is_data tinyint(1) NOT NULL default 0,
		idx int NOT NULL default 0,
		length int NOT NULL default 0,
		real_name varchar(255) NOT NULL default '',
		create_at datetime NOT NULL,
		INDEX idx_content_id (content_id)
	);`)

	// 6. 文件元数据表（修改：新增md5、content_id、is_deleted）
	DB.Exec(`
	CREATE TABLE IF NOT EXISTS file_metadata (
		id bigint PRIMARY KEY auto_increment,
		user_id bigint NOT NULL,
		real_name TEXT NOT NULL,
		encoded_name TEXT NOT NULL,
		dir varchar(60) NOT NULL DEFAULT '/',
		upload_time datetime,
		file_size bigint NOT NULL default 0,
		md5 varchar(32) NOT NULL,
		content_id bigint NOT NULL,
		is_deleted tinyint(1) NOT NULL default 0,
		INDEX idx_user_id_dir (user_id, dir),
		INDEX idx_md5 (md5),
		INDEX idx_content_id (content_id)
	);`)

	// 7. 分享表
	DB.Exec(`
	CREATE TABLE IF NOT EXISTS shares (
		id bigint PRIMARY KEY auto_increment,
		user_id bigint NOT NULL,
		file_id bigint NOT NULL,
		share_key varchar(255) not null default '',
		expire_at datetime,
		create_time datetime
	);`)

	return nil
}
