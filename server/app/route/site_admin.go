package route

import (
	"wnadmin/app/controller/site"
	"wnadmin/app/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterSiteAdminRoutes 官网资料和询盘
func RegisterSiteAdminRoutes(apiGroup *gin.RouterGroup) {
	controller := site.NewController()
	group := apiGroup.Group("/site-admin")
	{
		group.GET("/profile", middleware.WithSelect("企业资料", "查询企业资料"), controller.Profile)
		group.PUT("/profile", middleware.WithUpdate("企业资料", "保存企业资料"), controller.SaveProfile)
		group.GET("/inquiries", middleware.WithSelect("询盘", "查询询盘"), controller.Inquiries)
		group.PUT("/inquiries/:id/status", middleware.WithUpdate("询盘", "更新询盘状态"), controller.UpdateInquiryStatus)
	}
}