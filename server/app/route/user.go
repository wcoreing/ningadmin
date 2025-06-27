package route

import (
	"wnadmin/app/controller/user"
	"wnadmin/app/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterUserRoutes 注册用户路由
func RegisterUserRoutes(router *gin.RouterGroup) {
	userController := user.NewUserController()

	userGroup := router.Group("/user")
	{
		// 用户列表
		userGroup.GET("/list", middleware.WithSelect("用户管理", "查询用户列表"), userController.List)
		// 用户详情
		userGroup.GET("/:id", middleware.WithSelect("用户管理", "查询用户详情"), userController.GetInfo)
		// 新增用户
		userGroup.POST("", middleware.WithCreate("用户管理", "新增用户"), userController.Add)
		// 更新用户
		userGroup.PUT("", middleware.WithUpdate("用户管理", "修改用户"), userController.Update)
		// 删除用户
		userGroup.DELETE("/:id", middleware.WithDelete("用户管理", "删除用户"), userController.Delete)
		// 更新用户状态
		userGroup.PUT("/status", middleware.WithUpdate("用户管理", "修改用户状态"), userController.UpdateStatus)
		// 重置用户密码
		userGroup.PUT("/password", middleware.WithUpdate("用户管理", "重置用户密码"), userController.ResetPassword)
		// 更新用户余额
		userGroup.PUT("/balance", middleware.WithUpdate("用户管理", "更新用户余额"), userController.UpdateBalance)
		// 更新用户积分
		userGroup.PUT("/points", middleware.WithUpdate("用户管理", "更新用户积分"), userController.UpdatePoints)
		// 获取用户变动记录列表
		userGroup.GET("/change-log", middleware.WithSelect("用户管理", "查询变动记录"), userController.GetChangeLogList)
		// 获取指定用户的变动记录
		userGroup.GET("/:id/change-log", middleware.WithSelect("用户管理", "查询用户变动记录"), userController.GetUserChangeLog)
	}
}
