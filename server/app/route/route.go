package route

import (
	"wnadmin/app/middleware"

	"github.com/gin-gonic/gin"
)

func Handle(engine *gin.Engine) {
	// 注册公共路由(无需认证)
	RegisterPublicRoutes(engine)
	RegisterResourceRoutes(engine)
	// 需要认证的路由组
	apiGroup := engine.Group("/api")
	apiGroup.Use(middleware.JWT())
	{
		RegisterSystemRoutes(apiGroup)
		RegisterShopRoutes(apiGroup)
		RegisterDistributionRoutes(apiGroup)
		RegisterShareRoutes(apiGroup)
		RegisterFormRoutes(apiGroup)
		RegisterFileRoutes(apiGroup)
		RegisterCommonRoutes(apiGroup)
		RegisterUserRoutes(apiGroup)
	}

}
