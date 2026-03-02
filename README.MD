### 简单的分布式云存储系统


### 安装
开发环境  go 最新，  node >= 20,  mysql/mariadb
克隆仓库到本地，并进入

1. 修改config.yaml 的 db 改为自己的mysql数据库
```yaml
db: cander:123456@tcp(10.1.1.200:3306)/storage?charset=utf8mb4&parseTime=True&loc=Local
```
2. 运行master 节点
```go
go run .\cmd\master\master.go
```
3. 运行 node 节点， 代码写死 至少要启动5个node节点， 默认存储在当前目录的  store_<port>
```go
go run .\cmd\node\node.go 8003  // 8003 是端口号
```
4. 安装前端依赖
进入 frontend 目录
```
npm i -g @quasar/cli
npm install
```

5. 配置文件（.env） 只有一个master 的地址
API_URL=http://localhost:8080

6. 启动
```
quasar dev
```
7. 进入前段页面，注册一个账号， 暂时没做任何安全策略， 任何人都可以注册 
浏览器f12 调试模式里面找到任何一个请求， 查看 请求头的 Authorization, 替换下面的， 将 所有node 节点http://127.0.0.1:8001 加入进去
```
curl -X GET 'http://127.0.0.1:8080/node/add?addr=http://127.0.0.1:8005' -H 'Content-Type: application/json' -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJkaXN0X3N0b3JhZ2UiLCJleHAiOjE3NzI3ODE2NzQsImlhdCI6MTc3MjE3Njg3NCwidXNlcl9pZCI6MSwidXNlcm5hbWUiOiJjYW5kZXIifQ.2kRdGhnWApsIQZQtXCf2yGl-LZXHikOZU4bGBjJq4AM'
```
8. 可以上传试试
