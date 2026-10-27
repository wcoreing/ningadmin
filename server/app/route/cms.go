package route

import (
	"wnadmin/app/controller/cms"
	"wnadmin/app/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterCmsRoutes 内容发布
func RegisterCmsRoutes(apiGroup *gin.RouterGroup) {
	controller := cms.NewController()
	group := apiGroup.Group("/cms")
	{
		group.GET("/content/list", middleware.WithSelect("内容发布", "查询内容"), controller.ContentList)
		group.GET("/content/info/:id", middleware.WithSelect("内容发布", "查询内容详情"), controller.ContentInfo)
		group.POST("/content", middleware.WithCreate("内容发布", "保存内容"), controller.SaveContent)
		group.PUT("/content", middleware.WithUpdate("内容发布", "保存内容"), controller.SaveContent)
		group.DELETE("/content/:id", middleware.WithDelete("内容发布", "删除内容"), controller.DeleteContent)
	}
}
