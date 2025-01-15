package route

import (
	"github.com/gin-gonic/gin"
	"wnadmin/app/controller/entry"
)

func Handle(engine *gin.Engine) {
	engine.GET("/", entry.Handle)
}
