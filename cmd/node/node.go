package main

import (
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/hyahm/golog"
)

// 全局存储根目录（可通过命令行参数指定）
var rootStoreDir string

// ==================== 核心工具函数 ====================
// safeLocalPath 生成安全的本地存储路径
// 处理特殊字符、用户隔离、目录隔离，避免路径遍历攻击
func safeLocalPath(chunkID string) string {
	// 1. 解码URL转义的chunkID
	decodedCID, err := url.PathUnescape(chunkID)
	if err != nil {
		decodedCID = chunkID
	}

	// 2. 替换非法字符，确保路径安全
	// 移除 ../、./ 等路径遍历字符
	safeCID := filepath.Clean(decodedCID)
	safeCID = strings.ReplaceAll(safeCID, "..", "_")
	safeCID = strings.ReplaceAll(safeCID, "/", "_")
	safeCID = strings.ReplaceAll(safeCID, "\\", "_")
	safeCID = strings.ReplaceAll(safeCID, ":", "_")
	safeCID = strings.ReplaceAll(safeCID, "*", "_")
	safeCID = strings.ReplaceAll(safeCID, "?", "_")
	safeCID = strings.ReplaceAll(safeCID, "\"", "_")
	safeCID = strings.ReplaceAll(safeCID, "<", "_")
	safeCID = strings.ReplaceAll(safeCID, ">", "_")
	safeCID = strings.ReplaceAll(safeCID, "|", "_")

	// 3. 生成最终存储路径（根目录/节点端口/分片ID）
	return filepath.Join(rootStoreDir, safeCID)
}

// ==================== HTTP 处理器 ====================
// chunkHandler 处理分片的上传（PUT）和下载（GET）
func chunkHandler(w http.ResponseWriter, r *http.Request) {
	golog.Info("3333")
	// 1. 获取分片ID（cid参数）
	rawCID := r.URL.Query().Get("cid")
	if rawCID == "" {
		http.Error(w, "chunk ID (cid) is required", http.StatusBadRequest)
		return
	}

	// 2. 生成安全的本地存储路径
	localPath := safeLocalPath(rawCID)
	// 确保存储目录存在
	os.MkdirAll(filepath.Dir(localPath), 0755)

	// 3. 处理PUT请求（上传分片）
	if r.Method == http.MethodPut {
		// 读取请求体数据
		chunkData, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "failed to read chunk data: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		// 写入本地文件（覆盖已有文件）
		err = os.WriteFile(localPath, chunkData, 0644)
		if err != nil {
			http.Error(w, "failed to save chunk: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// 返回成功响应
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("chunk saved successfully: " + rawCID))
		return
	}

	// 4. 处理GET请求（下载分片）
	if r.Method == http.MethodGet {
		// 检查文件是否存在
		_, err := os.Stat(localPath)
		if os.IsNotExist(err) {
			http.Error(w, "chunk not found: "+rawCID, http.StatusNotFound)
			return
		} else if err != nil {
			http.Error(w, "failed to stat chunk: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// 读取文件并返回
		chunkData, err := os.ReadFile(localPath)
		if err != nil {
			http.Error(w, "failed to read chunk: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// 设置响应头（修复Content-Length类型错误）
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Length", strconv.Itoa(len(chunkData))) // 改为正确的int转string
		w.WriteHeader(http.StatusOK)
		w.Write(chunkData)
		return
	}

	// 5. 不支持的请求方法
	http.Error(w, "method not allowed (only PUT/GET)", http.StatusMethodNotAllowed)
}

// healthHandler 健康检查接口（可选，用于监控）
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"healthy","store_dir":"` + rootStoreDir + `"}`))
}

// ==================== 主函数 ====================
func main() {
	defer golog.Sync()
	// 1. 解析命令行参数
	if len(os.Args) < 2 {
		println("Usage: go run node.go <port> [store_dir]")
		println("Example: go run node.go 8081 ./store_8081")
		os.Exit(1)
	}

	// 2. 获取端口和存储目录
	port := os.Args[1]
	if len(os.Args) >= 3 {
		rootStoreDir = os.Args[2]
	} else {
		// 默认存储目录：./store_<port>
		rootStoreDir = filepath.Join(".", "store_"+port)
	}

	// 3. 初始化存储目录
	err := os.MkdirAll(rootStoreDir, 0755)
	if err != nil {
		println("Failed to create store directory:", err.Error())
		os.Exit(1)
	}

	// 4. 注册HTTP路由
	http.HandleFunc("/chunk", chunkHandler)
	http.HandleFunc("/health", healthHandler) // 健康检查

	// 5. 启动HTTP服务
	listenAddr := ":" + port
	println("====================================")
	println("Storage Node Starting...")
	println("Listen Address: " + listenAddr)
	println("Store Directory: " + rootStoreDir)
	println("====================================")

	err = http.ListenAndServe(listenAddr, nil)
	if err != nil {
		println("Failed to start node:", err.Error())
		os.Exit(1)
	}
}
