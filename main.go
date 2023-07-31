package main

import (
	"github.com/RaymondCode/simple-demo/service"
	"github.com/gin-gonic/gin"
)

func main() {
	//监听9090端口,前端会定时轮询服务端接口查询消息记录.
	go service.RunMessageServer()
	//初始化一个gin,采用默认的中间件服务
	r := gin.Default()
	//初始化路由
	initRouter(r)

	//默认是8080端口
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
