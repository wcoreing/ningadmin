package route

import (
	"wnadmin/app/controller/auth"
	"wnadmin/app/controller/entry"
	"wnadmin/app/controller/form"

	"github.com/gin-gonic/gin"
)

func RegisterPublicRoutes(engine *gin.Engine) {
	engine.GET("/", entry.Handle)
	engine.POST("/api/auth/login", auth.NewLoginController().Login)
	engine.POST("/api/auth/logout", auth.NewLoginController().Logout)
	engine.POST("/api/form/data/submit", form.NewFormDataController().SubmitForm)
}
