package main

import (
	"fmt"
	"runtime"
	"wnadmin/app/config"
	"wnadmin/app/middleware"
	"wnadmin/app/provider"
	"wnadmin/app/provider/mylog"
	"wnadmin/app/route"

	"github.com/gin-gonic/gin"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	provider.Register()

	// 设置Gin模式
	gin.SetMode(gin.ReleaseMode)

	// 创建Gin引擎
	r := gin.New() // 不使用默认的Logger和Recovery中间件

	// 使用自定义日志和恢复中间件
	r.Use(middleware.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.CORSMiddleware())

	// 注册路由
	route.Handle(r)

	// 启动服务
	serverAddr := fmt.Sprintf(":%d", config.App.Port)
	mylog.Log.Info("服务启动，监听地址: " + serverAddr)
	if err := r.Run(serverAddr); err != nil {
		mylog.Log.Fatalf("服务启动失败: %v", err)
	}
}
