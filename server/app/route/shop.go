package route

import (
	"wnadmin/app/controller/shop"

	"github.com/gin-gonic/gin"
)

// RegisterShopRoutes 注册商城相关路由
func RegisterShopRoutes(apiGroup *gin.RouterGroup) {
	// 商品管理路由组
	goodsController := &shop.GoodsController{}
	shopGroup := apiGroup.Group("/shop")
	{
		goodsGroup := shopGroup.Group("/goods")
		{
			goodsGroup.POST("/create", goodsController.Create)       // 创建商品
			goodsGroup.GET("/list", goodsController.List)            // 商品列表
			goodsGroup.GET("/info/:id", goodsController.GetInfo)     // 商品详情
			goodsGroup.PUT("/update", goodsController.Update)        // 更新商品
			goodsGroup.DELETE("/delete/:id", goodsController.Delete) // 删除商品
			goodsGroup.PUT("/status", goodsController.UpdateStatus)  // 更新状态
		}

		// 商品分类
		categoryController := shop.NewCategoryController()
		categoryGroup := shopGroup.Group("/category")
		{
			categoryGroup.GET("/list", categoryController.List)
		}

		// 商品标签
		tagController := shop.NewTagController()
		tagGroup := shopGroup.Group("/tag")
		{
			tagGroup.GET("/list", tagController.List)
		}

		// 商品品牌
		brandController := shop.NewBrandController()
		brandGroup := shopGroup.Group("/brand")
		{
			brandGroup.GET("/list", brandController.List)
		}

		// 商品供应商
		supplierController := shop.NewSupplierController()
		supplierGroup := shopGroup.Group("/supplier")
		{
			supplierGroup.GET("/list", supplierController.List)
		}
	}
}
