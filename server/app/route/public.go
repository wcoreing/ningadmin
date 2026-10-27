package route

import (
	"wnadmin/app/controller/auth"
	"wnadmin/app/controller/entry"
	"wnadmin/app/controller/form"
	"wnadmin/app/controller/site"

	"github.com/gin-gonic/gin"
)

func RegisterPublicRoutes(engine *gin.Engine) {
	engine.GET("/", entry.Handle)
	engine.POST("/api/auth/login", auth.NewLoginController().Login)
	engine.POST("/api/auth/logout", auth.NewLoginController().Logout)
	engine.POST("/api/form/data/submit", form.NewFormDataController().SubmitForm)

	siteController := site.NewController()
	siteGroup := engine.Group("/api/site")
	{
		siteGroup.GET("/home", siteController.Home)
		siteGroup.GET("/profile", siteController.Profile)
		siteGroup.GET("/pages/:code", siteController.Page)
		siteGroup.GET("/categories", siteController.Categories)
		siteGroup.GET("/goods", siteController.Goods)
		siteGroup.GET("/goods/:id", siteController.GoodsDetail)
		siteGroup.GET("/articles", siteController.Articles)
		siteGroup.GET("/articles/:id", siteController.Article)
		siteGroup.POST("/inquiries", siteController.CreateInquiry)
	}
}
