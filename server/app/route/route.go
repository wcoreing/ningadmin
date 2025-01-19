package route

import (
	"wnadmin/app/controller/auth"
	"wnadmin/app/controller/entry"
	"wnadmin/app/controller/form"
	"wnadmin/app/controller/system"
	"wnadmin/app/middleware"

	"github.com/gin-gonic/gin"
)

func Handle(engine *gin.Engine) {
	engine.GET("/", entry.Handle)

	// 注册系统配置路由
	configController := system.NewConfigController()
	menuController := &system.MenuController{}
	roleController := &system.RoleController{}
	userController := &system.UserController{}
	loginController := auth.NewLoginController()
	formController := form.NewFormController()

	// 登录相关路由（无需认证）
	engine.POST("/api/auth/login", loginController.Login)
	engine.POST("/api/auth/logout", loginController.Logout)
	formPublicGroup := engine.Group("/form")
	{
		formPublicGroup.GET("/share/:shareCode", formController.GetShareInfo)
		formPublicGroup.POST("/submit", formController.SubmitForm)
	}
	// 需要认证的路由组
	apiGroup := engine.Group("/api")
	apiGroup.Use(middleware.JWT())
	{
		// 获取当前用户信息
		apiGroup.GET("/auth/info", loginController.GetUserInfo)

		// 系统管理路由组
		systemGroup := apiGroup.Group("/system")
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
			formGroup := apiGroup.Group("/form")
			// 表单管理
			formTemplateGroup := formGroup.Group("/form-template")
			{
				formTemplateGroup.GET("/list", formController.GetTemplateList)
				formTemplateGroup.POST("/save", formController.SaveTemplate)
				formTemplateGroup.DELETE("/:id", formController.DeleteTemplate)
			}

			// 表单分享管理
			formShareGroup := formGroup.Group("/form-share")
			{
				formShareGroup.POST("/create", formController.CreateShare)
			}

			// 表单提交记录管理
			formSubmissionGroup := formGroup.Group("/form-submission")
			{
				formSubmissionGroup.GET("/list", formController.GetSubmissionList)
			}
			// 公开的表单路由(不需要认证)

		}

	}
}
