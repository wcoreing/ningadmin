package route

import (
	"wnadmin/app/controller/shop"

	"github.com/gin-gonic/gin"
)

// InitShopRoutes 初始化商品相关路由
func InitShopRoutes(apiGroup *gin.RouterGroup) {
	shopGroup := apiGroup.Group("/shop")
	{
		// 商品管理
		goodsController := shop.NewGoodsController()
		goodsGroup := shopGroup.Group("/goods")
		{
			goodsGroup.GET("/list", goodsController.List)
			goodsGroup.DELETE("/delete/:id", goodsController.Delete)
			goodsGroup.PUT("/status", goodsController.UpdateStatus)
			goodsGroup.GET("/info/:id", goodsController.GetInfo)
			goodsGroup.PUT("/update", goodsController.Update)
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
