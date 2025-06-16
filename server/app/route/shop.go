package route

import (
	"wnadmin/app/controller/shop"

	"github.com/gin-gonic/gin"
)

// RegisterShopRoutes 注册商城相关路由
func RegisterShopRoutes(apiGroup *gin.RouterGroup) {
	// 商品管理路由组
	goodsController := shop.NewGoodsController()
	shopGroup := apiGroup.Group("/shop")
	{
		goodsGroup := shopGroup.Group("/goods")
		{
			goodsGroup.GET("/list", goodsController.List)
			goodsGroup.GET("/info/:id", goodsController.GetInfo)
			goodsGroup.POST("create", goodsController.Create)
			goodsGroup.PUT("", goodsController.Update)
			goodsGroup.DELETE("/:id", goodsController.Delete)
			goodsGroup.PUT("/status", goodsController.UpdateStatus)
			goodsGroup.POST("/copy/:id", goodsController.CopyGoods)
			goodsGroup.PUT("/move", goodsController.BatchMove)
			goodsGroup.PUT("/update", goodsController.Update)
			goodsGroup.DELETE("/batch", goodsController.BatchDelete)

			// 套餐相关路由
			comboGroup := goodsGroup.Group("/combo")
			{
				comboGroup.GET("/config/:id", goodsController.GetComboConfig)            // 获取套餐配置
				comboGroup.POST("/calculate", goodsController.CalculateComboPrice)       // 计算套餐价格
				comboGroup.POST("/validate", goodsController.ValidateComboSelection)     // 验证套餐选择
				comboGroup.GET("/default/:id", goodsController.GetDefaultComboSelection) // 获取默认套餐选择
			}
		}

		// 订单管理
		orderController := shop.NewOrderController()
		orderGroup := shopGroup.Group("/order")
		{
			orderGroup.GET("/list", orderController.GetOrderList)
			orderGroup.PUT("/ship", orderController.ShipOrder)
			orderGroup.POST("/logistics", orderController.GetLogisticsTrack)
		}

		// 商品分类
		categoryController := shop.NewCategoryController()
		categoryGroup := shopGroup.Group("/categories")
		{
			categoryGroup.GET("", categoryController.List)
			categoryGroup.POST("", categoryController.Create)
			categoryGroup.PUT("/:id", categoryController.Update)
			categoryGroup.DELETE("/:id", categoryController.Delete)
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

		// Banner管理
		bannerController := shop.NewBannerController()
		bannerGroup := shopGroup.Group("/banner")
		{
			bannerGroup.GET("/list", bannerController.GetBannerList)            // 获取Banner列表
			bannerGroup.GET("/position", bannerController.GetBannersByPosition) // 获取指定位置的Banner列表
			bannerGroup.GET("/info/:id", bannerController.GetBannerInfo)
			bannerGroup.POST("/create", bannerController.CreateBanner)          // 创建Banner
			bannerGroup.PUT("/update", bannerController.UpdateBanner)           // 更新Banner
			bannerGroup.DELETE("/delete/:id", bannerController.DeleteBanner)    // 删除Banner
			bannerGroup.PUT("/status/:id", bannerController.ChangeBannerStatus) // 修改Banner状态
		}

		// 用户反馈管理
		feedbackController := shop.NewFeedbackController()
		feedbackGroup := shopGroup.Group("/feedback")
		{
			feedbackGroup.GET("/list", feedbackController.List)               // 获取反馈列表
			feedbackGroup.GET("/:id", feedbackController.GetInfo)             // 获取反馈详情
			feedbackGroup.POST("", feedbackController.Add)                    // 新增反馈
			feedbackGroup.PUT("", feedbackController.Update)                  // 更新反馈
			feedbackGroup.DELETE("/:id", feedbackController.Delete)           // 删除反馈
			feedbackGroup.POST("/:id/reply", feedbackController.Reply)        // 回复反馈
			feedbackGroup.PUT("/:id/status", feedbackController.UpdateStatus) // 更新反馈状态
		}

	}
}
