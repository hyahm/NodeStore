package app

import (
	"fmt"
	"nodestore/app/handles"
	"nodestore/app/response"

	"github.com/hyahm/xmux"
)

func Run(port int) {

	res := &response.Response{}
	router := xmux.NewRouter()
	router.SetHeader("Access-Control-Allow-Origin", "*")
	router.SetHeader("Access-Control-Allow-Headers", "token,content-type,authorization")
	// 用户
	router.BindResponse(res)
	router.Post("/user/register", handles.NewLogin().RegisterHandler).BindJson(handles.NewLogin())
	router.Post("/user/login", handles.NewLogin().LoginHandler).BindJson(handles.NewLogin())

	// 节点
	node := handles.Node{}

	router.Get("/node/add", node.AddNodeHandler)
	router.Get("/node/remove", node.RemoveNodeHandler)
	router.Get("/node/list", node.ListNodesHandler)

	// 目录
	router.Post("/dir/create", handles.NewCreateDir().CreateDirHandler).BindJson(handles.NewCreateDir())

	// 文件系统列表
	router.Get("/fs/list", handles.NewCreateDir().FsListHandler)

	// 文件操作
	file := handles.File{}
	router.Post("/file/upload", file.UploadHandler)
	router.Get("/file/download_by_md5", file.DownloadByMD5Handler).BindResponse(nil) // 按MD5下载
	router.Get("/file/stream", file.StreamHandler).BindResponse(nil)                 // 在线播放
	router.Get("/file/delete", file.DeleteFileHandler)                               // 软删除（支持MD5）

	// 分享
	share := handles.Share{}
	router.Get("/share/create", share.ShareCreateHandler).BindResponse(nil)
	router.Get("/share/download", share.ShareDownloadHandler).BindResponse(nil)

	fmt.Println("master :8080 - MD5 download + online stream (mp3/mp4)")
	router.Run(fmt.Sprintf(":%d", port))
}
