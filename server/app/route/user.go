package route

import (
	"wnadmin/app/controller/user"

	"github.com/gin-gonic/gin"
)

// RegisterUserRoutes 注册用户路由
func RegisterUserRoutes(router *gin.RouterGroup) {
	userController := user.NewUserController()

	userGroup := router.Group("/user")
	{
		// 用户列表
		userGroup.GET("/list", userController.List)
		// 用户详情
		userGroup.GET("/:id", userController.GetInfo)
		// 新增用户
		userGroup.POST("", userController.Add)
		// 更新用户
		userGroup.PUT("", userController.Update)
		// 删除用户
		userGroup.DELETE("/:id", userController.Delete)
		// 更新用户状态
		userGroup.PUT("/status", userController.UpdateStatus)
		// 重置用户密码
		userGroup.PUT("/password", userController.ResetPassword)
		// 更新用户余额
		userGroup.PUT("/balance", userController.UpdateBalance)
		// 更新用户积分
		userGroup.PUT("/points", userController.UpdatePoints)
		// 获取用户变动记录列表
		userGroup.GET("/change-log", userController.GetChangeLogList)
		// 获取指定用户的变动记录
		userGroup.GET("/:id/change-log", userController.GetUserChangeLog)
	}
}
