package route

import (
	"wnadmin/app/controller/common"

	"github.com/gin-gonic/gin"
)

func RegisterFileRoutes(r *gin.RouterGroup) {
	fileController := common.NewFileController()
	fileGroup := r.Group("/files")
	{
		fileGroup.POST("/upload", fileController.Upload)
		fileGroup.GET("", fileController.List)
		fileGroup.DELETE("/:id", fileController.Delete)
		fileGroup.GET("/categories", fileController.CategoryTree)
		fileGroup.POST("/move", fileController.MoveFiles)
	}
}
