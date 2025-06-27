package route

import (
	"wnadmin/app/controller/common"
	"wnadmin/app/middleware"

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
			categoryGroup.GET("/list", middleware.WithSelect("通用分类", "查询分类列表"), categoryController.GetCategoryList)
			categoryGroup.GET("/tree", middleware.WithSelect("通用分类", "查询分类树"), categoryController.GetCategoryTree)
			categoryGroup.POST("", middleware.WithCreate("通用分类", "新增分类"), categoryController.CreateCategory)
			categoryGroup.PUT("/:id", middleware.WithUpdate("通用分类", "修改分类"), categoryController.UpdateCategory)
			categoryGroup.DELETE("/:id", middleware.WithDelete("通用分类", "删除分类"), categoryController.DeleteCategory)
		}
	}
}
