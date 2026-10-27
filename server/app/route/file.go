package route

import (
	"wnadmin/app/controller/common"
	"wnadmin/app/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterFileRoutes(r *gin.RouterGroup) {
	fileController := common.NewFileController()
	fileGroup := r.Group("/files")
	{
		fileGroup.POST("/upload", middleware.WithCreate("文件管理", "上传文件"), fileController.Upload)
		fileGroup.GET("", middleware.WithSelect("文件管理", "查询文件列表"), fileController.List)
		fileGroup.DELETE("/:id", middleware.WithDelete("文件管理", "删除文件"), fileController.Delete)
		fileGroup.GET("/categories", middleware.WithSelect("文件管理", "查询文件分类"), fileController.CategoryTree)
		fileGroup.POST("/move", middleware.WithUpdate("文件管理", "移动文件"), fileController.MoveFiles)
	}
}
