package route

import (
	"wnadmin/app/controller/auth"
	"wnadmin/app/controller/common"
	"wnadmin/app/controller/entry"
	"wnadmin/app/controller/form"
	"wnadmin/app/controller/system"
	"wnadmin/app/middleware"

	"net/http"
	"os"
	"path"

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
	formListController := form.NewFormListController()

	// 登录相关路由（无需认证）
	engine.POST("/api/auth/login", loginController.Login)
	engine.POST("/api/auth/logout", loginController.Logout)
	engine.POST("/api/form/data/submit", form.NewFormDataController().SubmitForm)

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
				menuGroup.PUT("/sort", menuController.UpdateSort)
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
			formTemplateGroup := formGroup.Group("/template")
			{
				formTemplateGroup.GET("/list", formController.GetTemplateList)
				formTemplateGroup.POST("", formController.CreateTemplate)
				formTemplateGroup.PUT("", formController.UpdateTemplate)
				formTemplateGroup.DELETE("/:id", formController.DeleteTemplate)
				formTemplateGroup.GET("/:id", formController.GetTemplateInfo)
			}

			// 表单分享管理
			formShareGroup := formGroup.Group("/share")
			{
				formShareGroup.POST("/create", formController.CreateShare)
			}

			// 表单管理
			formFormGroup := formGroup.Group("form")
			{
				formFormGroup.GET("list", formListController.GetFormList)
				formFormGroup.GET("/:id", formListController.GetFormDetail)
				formFormGroup.POST("", formListController.CreateForm)
				formFormGroup.PUT("", formListController.UpdateForm)
				formFormGroup.DELETE("/:id", formListController.DeleteForm)
			}

			// 表单数据管理
			formDataController := form.NewFormDataController()
			formDataGroup := formGroup.Group("/data")
			{
				formDataGroup.GET("/list", formDataController.GetFormDataList)
				formDataGroup.GET("/:id", formDataController.GetFormDataDetail)
				formDataGroup.DELETE("/:id", formDataController.DeleteFormData)
			}

			// 添加文件管理路由
			fileController := common.NewFileController()
			fileGroup := apiGroup.Group("/files")
			{
				fileGroup.POST("/upload", fileController.Upload)
				fileGroup.GET("", fileController.List)
				fileGroup.DELETE("/:id", fileController.Delete)
				fileGroup.GET("/categories", fileController.CategoryTree)
				fileGroup.POST("/move", fileController.MoveFiles)
			}

			// 注册商城相关路由
			RegisterShopRoutes(apiGroup)
		}
	}

	// 添加静态资源路由，只需要一个就够了
	uploadPath := "./uploads" // 文件上传目录
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
