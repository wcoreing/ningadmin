package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"runtime"
	"wnadmin/app/config"
	"wnadmin/app/middleware"
	"wnadmin/app/provider"
	"wnadmin/app/route"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	provider.Register()
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(middleware.CORSMiddleware())
	route.Handle(r)
	r.Run(fmt.Sprintf(":%d", config.App.Port)) // listen and serve on 0.0.0.0:8080
}
