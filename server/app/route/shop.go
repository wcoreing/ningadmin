package route

import (
	"wnadmin/app/controller/shop"

	"github.com/gin-gonic/gin"
)

// InitShopRoutes 初始化商品相关路由
func InitShopRoutes(apiGroup *gin.RouterGroup) {
	// 商品相关路由
	goodsController := shop.NewGoodsController()
	shopGroup := apiGroup.Group("/shop")
	{
		// 商品管理
		goodsGroup := shopGroup.Group("/goods")
		{
			goodsGroup.GET("/list", goodsController.List)
			goodsGroup.DELETE("/delete/:id", goodsController.Delete)
			goodsGroup.PUT("/status", goodsController.UpdateStatus)
		}
	}
}
