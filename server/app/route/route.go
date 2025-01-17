package route

import (
	"wnadmin/app/controller/entry"
	"wnadmin/app/controller/system"

	"github.com/gin-gonic/gin"
)

func Handle(engine *gin.Engine) {
	engine.GET("/", entry.Handle)

	// 注册系统配置路由
	configController := &system.ConfigController{}
	menuController := &system.MenuController{}
	roleController := &system.RoleController{}
	userController := &system.UserController{}

	systemGroup := engine.Group("/api/system")
	{
		// 配置管理
		configGroup := systemGroup.Group("/config")
		{
			configGroup.GET("/list", configController.List)
			configGroup.POST("/add", configController.Add)
			configGroup.PUT("/update", configController.Update)
			configGroup.DELETE("/delete/:id", configController.Delete)
			configGroup.GET("/info/:id", configController.GetInfo)
		}

		// 菜单管理
		menuGroup := systemGroup.Group("/menu")
		{
			menuGroup.GET("/list", menuController.List)
			menuGroup.POST("/add", menuController.Add)
			menuGroup.PUT("/update", menuController.Update)
			menuGroup.DELETE("/delete/:id", menuController.Delete)
			menuGroup.GET("/info/:id", menuController.GetInfo)
			menuGroup.GET("/tree", menuController.GetTree)
			menuGroup.GET("/check-permission", menuController.CheckPermission)
		}

		// 角色管理
		roleGroup := systemGroup.Group("/role")
		{
			roleGroup.GET("/list", roleController.List)
			roleGroup.GET("/:roleId", roleController.GetInfo)
			roleGroup.POST("", roleController.Add)
			roleGroup.PUT("", roleController.Update)
			roleGroup.DELETE("/:roleId", roleController.Delete)
			roleGroup.PUT("/status", roleController.UpdateStatus)
			roleGroup.GET("/permission/:roleId", roleController.GetRolePermissions)
			roleGroup.PUT("/permission", roleController.UpdateRolePermissions)
		}

		// 用户管理
		userGroup := systemGroup.Group("/user")
		{
			userGroup.GET("/list", userController.List)
			userGroup.GET("/info/:userId", userController.GetInfo)
			userGroup.POST("", userController.Add)
			userGroup.PUT("", userController.Update)
			userGroup.DELETE("/delete/:userId", userController.Delete)
			userGroup.PUT("/status", userController.UpdateStatus)
			userGroup.PUT("/password", userController.ResetPassword)
		}
	}
}
