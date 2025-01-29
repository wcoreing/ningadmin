package main

import (
	"fmt"
	"runtime"
	"wnadmin/app/config"
	"wnadmin/app/middleware"
	"wnadmin/app/provider"
	"wnadmin/app/route"

	"github.com/gin-gonic/gin"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	provider.Register()
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(middleware.CORSMiddleware())
	route.Handle(r)
	r.Run(fmt.Sprintf(":%d", config.App.Port))
}
