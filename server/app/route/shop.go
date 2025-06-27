package route

import (
	"wnadmin/app/controller/shop"
	"wnadmin/app/middleware"

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
			goodsGroup.GET("/list", middleware.WithSelect("商品管理", "查询商品列表"), goodsController.List)
			goodsGroup.GET("/info/:id", middleware.WithSelect("商品管理", "查询商品详情"), goodsController.GetInfo)
			goodsGroup.POST("create", middleware.WithCreate("商品管理", "新增商品"), goodsController.Create)
			goodsGroup.PUT("", middleware.WithUpdate("商品管理", "修改商品"), goodsController.Update)
			goodsGroup.DELETE("/:id", middleware.WithDelete("商品管理", "删除商品"), goodsController.Delete)
			goodsGroup.PUT("/status", middleware.WithUpdate("商品管理", "修改商品状态"), goodsController.UpdateStatus)
			goodsGroup.POST("/copy/:id", middleware.WithCreate("商品管理", "复制商品"), goodsController.CopyGoods)
			goodsGroup.PUT("/move", middleware.WithUpdate("商品管理", "批量移动商品"), goodsController.BatchMove)
			goodsGroup.PUT("/update", middleware.WithUpdate("商品管理", "更新商品"), goodsController.Update)
			goodsGroup.DELETE("/batch", middleware.WithDelete("商品管理", "批量删除商品"), goodsController.BatchDelete)

			// 套餐相关路由
			comboGroup := goodsGroup.Group("/combo")
			{
				comboGroup.GET("/config/:id", middleware.WithSelect("商品管理", "获取套餐配置"), goodsController.GetComboConfig)              // 获取套餐配置
				comboGroup.POST("/calculate", middleware.WithSelect("商品管理", "计算套餐价格"), goodsController.CalculateComboPrice)         // 计算套餐价格
				comboGroup.POST("/validate", middleware.WithSelect("商品管理", "验证套餐选择"), goodsController.ValidateComboSelection)       // 验证套餐选择
				comboGroup.GET("/default/:id", middleware.WithSelect("商品管理", "获取默认套餐选择"), goodsController.GetDefaultComboSelection) // 获取默认套餐选择
			}
		}

		// 订单管理
		orderController := shop.NewOrderController()
		orderGroup := shopGroup.Group("/order")
		{
			orderGroup.GET("/list", middleware.WithSelect("订单管理", "查询订单列表"), orderController.GetOrderList)
			orderGroup.PUT("/ship", middleware.WithUpdate("订单管理", "发货"), orderController.ShipOrder)
			orderGroup.POST("/logistics", middleware.WithSelect("订单管理", "查询物流轨迹"), orderController.GetLogisticsTrack)
		}

		// 商品分类
		categoryController := shop.NewCategoryController()
		categoryGroup := shopGroup.Group("/categories")
		{
			categoryGroup.GET("", middleware.WithSelect("商品分类", "查询分类列表"), categoryController.List)
			categoryGroup.POST("", middleware.WithCreate("商品分类", "新增分类"), categoryController.Create)
			categoryGroup.PUT("/:id", middleware.WithUpdate("商品分类", "修改分类"), categoryController.Update)
			categoryGroup.DELETE("/:id", middleware.WithDelete("商品分类", "删除分类"), categoryController.Delete)
		}

		// 商品标签
		tagController := shop.NewTagController()
		tagGroup := shopGroup.Group("/tag")
		{
			tagGroup.GET("/list", middleware.WithSelect("商品标签", "查询标签列表"), tagController.List)
		}

		// 商品品牌
		brandController := shop.NewBrandController()
		brandGroup := shopGroup.Group("/brand")
		{
			brandGroup.GET("/list", middleware.WithSelect("商品品牌", "查询品牌列表"), brandController.List)
		}

		// 商品供应商
		supplierController := shop.NewSupplierController()
		supplierGroup := shopGroup.Group("/supplier")
		{
			supplierGroup.GET("/list", middleware.WithSelect("商品供应商", "查询供应商列表"), supplierController.List)
		}

		// Banner管理
		bannerController := shop.NewBannerController()
		bannerGroup := shopGroup.Group("/banner")
		{
			bannerGroup.GET("/list", middleware.WithSelect("Banner管理", "查询Banner列表"), bannerController.GetBannerList)              // 获取Banner列表
			bannerGroup.GET("/position", middleware.WithSelect("Banner管理", "查询指定位置Banner"), bannerController.GetBannersByPosition) // 获取指定位置的Banner列表
			bannerGroup.GET("/info/:id", middleware.WithSelect("Banner管理", "查询Banner详情"), bannerController.GetBannerInfo)
			bannerGroup.POST("/create", middleware.WithCreate("Banner管理", "新增Banner"), bannerController.CreateBanner)            // 创建Banner
			bannerGroup.PUT("/update", middleware.WithUpdate("Banner管理", "修改Banner"), bannerController.UpdateBanner)             // 更新Banner
			bannerGroup.DELETE("/delete/:id", middleware.WithDelete("Banner管理", "删除Banner"), bannerController.DeleteBanner)      // 删除Banner
			bannerGroup.PUT("/status/:id", middleware.WithUpdate("Banner管理", "修改Banner状态"), bannerController.ChangeBannerStatus) // 修改Banner状态
		}

		// 用户反馈管理
		feedbackController := shop.NewFeedbackController()
		feedbackGroup := shopGroup.Group("/feedback")
		{
			feedbackGroup.GET("/list", middleware.WithSelect("用户反馈", "查询反馈列表"), feedbackController.List)               // 获取反馈列表
			feedbackGroup.GET("/:id", middleware.WithSelect("用户反馈", "查询反馈详情"), feedbackController.GetInfo)             // 获取反馈详情
			feedbackGroup.POST("", middleware.WithCreate("用户反馈", "新增反馈"), feedbackController.Add)                      // 新增反馈
			feedbackGroup.PUT("", middleware.WithUpdate("用户反馈", "修改反馈"), feedbackController.Update)                    // 更新反馈
			feedbackGroup.DELETE("/:id", middleware.WithDelete("用户反馈", "删除反馈"), feedbackController.Delete)             // 删除反馈
			feedbackGroup.POST("/:id/reply", middleware.WithCreate("用户反馈", "回复反馈"), feedbackController.Reply)          // 回复反馈
			feedbackGroup.PUT("/:id/status", middleware.WithUpdate("用户反馈", "修改反馈状态"), feedbackController.UpdateStatus) // 更新反馈状态
		}

	}
}
