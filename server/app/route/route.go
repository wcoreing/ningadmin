package route

import (
	"wnadmin/app/controller/entry"
	"wnadmin/app/controller/system"

	"github.com/gin-gonic/gin"
)

func Handle(engine *gin.Engine) {
	engine.GET("/", entry.Handle)

	// 注册系统配置路由
	configController := &system.ConfigController{}
	systemGroup := engine.Group("/api/system")
	{
		configGroup := systemGroup.Group("/config")
		{
			configGroup.GET("/list", configController.List)
			configGroup.POST("/add", configController.Add)
			configGroup.PUT("/update", configController.Update)
			configGroup.DELETE("/delete/:id", configController.Delete)
			configGroup.GET("/info/:id", configController.GetInfo)
		}
	}
}
