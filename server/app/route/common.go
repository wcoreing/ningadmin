package route

import (
	"wnadmin/app/controller/common"

	"github.com/gin-gonic/gin"
)

// RegisterCommonRoutes 注册通用路由
func RegisterCommonRoutes(apiGroup *gin.RouterGroup) {
	// 通用分类路由
	categoryController := common.NewCategoryController()
	commonGroup := apiGroup.Group("/common")
	{
		categoryGroup := commonGroup.Group("/category")
		{
			categoryGroup.GET("/list", categoryController.GetCategoryList)
			categoryGroup.GET("/tree", categoryController.GetCategoryTree)
			categoryGroup.POST("", categoryController.CreateCategory)
			categoryGroup.PUT("/:id", categoryController.UpdateCategory)
			categoryGroup.DELETE("/:id", categoryController.DeleteCategory)
		}
	}
}
