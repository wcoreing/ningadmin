package route

import (
	"wnadmin/app/controller/auth"
	"wnadmin/app/controller/system"

	"github.com/gin-gonic/gin"
)

func RegisterSystemRoutes(r *gin.RouterGroup) {
	authController := auth.NewLoginController()
	r.GET("/auth/info", authController.GetUserInfo)

	systemGroup := r.Group("/system")
	{
		// 配置管理
		configController := system.NewConfigController()
		configGroup := systemGroup.Group("/config")
		{
			configGroup.GET("/list", configController.List)
			configGroup.GET("/tree", configController.GetTree)
			configGroup.POST("/add", configController.Add)
			configGroup.PUT("/update", configController.Update)
			configGroup.DELETE("/delete/:id", configController.Delete)
			configGroup.GET("/info/:id", configController.GetInfo)
		}

		// 菜单管理
		menuController := system.NewMenuController()
		menuGroup := systemGroup.Group("/menu")
		{
			menuGroup.GET("/list", menuController.List)
			menuGroup.POST("/add", menuController.Add)
			menuGroup.PUT("/update", menuController.Update)
			menuGroup.DELETE("/delete/:id", menuController.Delete)
			menuGroup.GET("/info/:id", menuController.GetInfo)
			menuGroup.GET("/tree", menuController.GetTree)
			menuGroup.GET("/check-permission", menuController.CheckPermission)
			menuGroup.PUT("/sort", menuController.UpdateSort)
		}

		// 角色管理
		roleController := &system.RoleController{}
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
		userController := &system.UserController{}
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

		// 弹窗管理
		popupController := system.NewPopupController()
		popupGroup := systemGroup.Group("/popup")
		{
			popupGroup.GET("/list", popupController.List)
			popupGroup.GET("/:id", popupController.GetInfo)
			popupGroup.POST("", popupController.Add)
			popupGroup.PUT("", popupController.Update)
			popupGroup.DELETE("/:id", popupController.Delete)
			popupGroup.PUT("/status", popupController.UpdateStatus)
			popupGroup.PUT("/stats", popupController.UpdateStats)
			popupGroup.GET("/active", popupController.GetActivePopups)
		}
	}
}
