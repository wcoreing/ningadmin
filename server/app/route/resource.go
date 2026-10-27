package route

import (
	"net/http"
	"os"
	"path"

	"github.com/gin-gonic/gin"
)

func RegisterResourceRoutes(engine *gin.Engine) {
	// 静态资源路由
	uploadPath := "./uploads"
	engine.Static("/uploads", uploadPath)

	// 添加一个处理不存在文件的中间件
	engine.Use(func(c *gin.Context) {
		if c.Request.Method == http.MethodGet {
			if path.Ext(c.Request.URL.Path) != "" { // 如果是文件请求
				// 检查文件是否存在
				if _, err := os.Stat(uploadPath + c.Request.URL.Path); os.IsNotExist(err) {
					c.File("./uploads/default.png") // 返回默认图片
					c.Abort()
					return
				}
			}
		}
		c.Next()
	})
}
