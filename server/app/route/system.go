package route

import (
	"wnadmin/app/controller/auth"
	"wnadmin/app/controller/system"
	"wnadmin/app/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterSystemRoutes(r *gin.RouterGroup) {
	authController := auth.NewLoginController()
	r.GET("/auth/info", middleware.WithSelect("用户认证", "获取用户信息"), authController.GetUserInfo)

	systemGroup := r.Group("/system")
	{
		// 配置管理
		configController := system.NewConfigController()
		configGroup := systemGroup.Group("/config")
		{
			configGroup.GET("/list", middleware.WithSelect("系统配置", "查询配置列表"), configController.List)
			configGroup.GET("/tree", middleware.WithSelect("系统配置", "查询配置树"), configController.GetTree)
			configGroup.POST("/add", middleware.WithCreate("系统配置", "新增配置"), configController.Add)
			configGroup.PUT("/update", middleware.WithUpdate("系统配置", "修改配置"), configController.Update)
			configGroup.DELETE("/delete/:id", middleware.WithDelete("系统配置", "删除配置"), configController.Delete)
			configGroup.GET("/info/:id", middleware.WithSelect("系统配置", "查询配置详情"), configController.GetInfo)
		}

		// 菜单管理
		menuController := system.NewMenuController()
		menuGroup := systemGroup.Group("/menu")
		{
			menuGroup.GET("/list", middleware.WithSelect("菜单管理", "查询菜单列表"), menuController.List)
			menuGroup.POST("/add", middleware.WithCreate("菜单管理", "新增菜单"), menuController.Add)
			menuGroup.PUT("/update", middleware.WithUpdate("菜单管理", "修改菜单"), menuController.Update)
			menuGroup.DELETE("/delete/:id", middleware.WithDelete("菜单管理", "删除菜单"), menuController.Delete)
			menuGroup.GET("/info/:id", middleware.WithSelect("菜单管理", "查询菜单详情"), menuController.GetInfo)
			menuGroup.GET("/tree", middleware.WithSelect("菜单管理", "查询菜单树"), menuController.GetTree)
			menuGroup.GET("/check-permission", middleware.WithSelect("菜单管理", "检查菜单权限"), menuController.CheckPermission)
			menuGroup.PUT("/sort", middleware.WithUpdate("菜单管理", "菜单排序"), menuController.UpdateSort)
			menuGroup.PUT("/status", middleware.WithUpdate("菜单管理", "修改菜单状态"), menuController.UpdateStatus)
		}

		// 角色管理
		roleController := &system.RoleController{}
		roleGroup := systemGroup.Group("/role")
		{
			roleGroup.GET("/list", middleware.WithSelect("角色管理", "查询角色列表"), roleController.List)
			roleGroup.GET("/:roleId", middleware.WithSelect("角色管理", "查询角色详情"), roleController.GetInfo)
			roleGroup.POST("", middleware.WithCreate("角色管理", "新增角色"), roleController.Add)
			roleGroup.PUT("", middleware.WithUpdate("角色管理", "修改角色"), roleController.Update)
			roleGroup.DELETE("/:roleId", middleware.WithDelete("角色管理", "删除角色"), roleController.Delete)
			roleGroup.PUT("/status", middleware.WithUpdate("角色管理", "修改角色状态"), roleController.UpdateStatus)
			roleGroup.GET("/permission/:roleId", middleware.WithSelect("角色管理", "查询角色权限"), roleController.GetRolePermissions)
			roleGroup.PUT("/permission", middleware.WithUpdate("角色管理", "修改角色权限"), roleController.UpdateRolePermissions)
		}

		// 用户管理
		userController := &system.UserController{}
		userGroup := systemGroup.Group("/user")
		{
			userGroup.GET("/list", middleware.WithSelect("用户管理", "查询用户列表"), userController.List)
			userGroup.GET("/info/:userId", middleware.WithSelect("用户管理", "查询用户详情"), userController.GetInfo)
			userGroup.POST("", middleware.WithCreate("用户管理", "新增用户"), userController.Add)
			userGroup.PUT("", middleware.WithUpdate("用户管理", "修改用户"), userController.Update)
			userGroup.DELETE("/delete/:userId", middleware.WithDelete("用户管理", "删除用户"), userController.Delete)
			userGroup.PUT("/status", middleware.WithUpdate("用户管理", "修改用户状态"), userController.UpdateStatus)
			userGroup.PUT("/password", middleware.WithUpdate("用户管理", "重置用户密码"), userController.ResetPassword)
		}

		// 部门管理
		deptController := system.NewDeptController()
		deptGroup := systemGroup.Group("/dept")
		{
			deptGroup.GET("/list", middleware.WithSelect("部门管理", "查询部门列表"), deptController.List)
			deptGroup.GET("/tree", middleware.WithSelect("部门管理", "查询部门树"), deptController.GetTree)
			deptGroup.POST("", middleware.WithCreate("部门管理", "新增部门"), deptController.Add)
			deptGroup.PUT("", middleware.WithUpdate("部门管理", "修改部门"), deptController.Update)
			deptGroup.DELETE("/:deptId", middleware.WithDelete("部门管理", "删除部门"), deptController.Delete)
			deptGroup.GET("/:deptId", middleware.WithSelect("部门管理", "查询部门详情"), deptController.GetInfo)
			deptGroup.PUT("/status", middleware.WithUpdate("部门管理", "修改部门状态"), deptController.UpdateStatus)
			deptGroup.GET("/options", middleware.WithSelect("部门管理", "查询部门选项"), deptController.GetDeptOptions)
		}

		// 弹窗管理
		popupController := system.NewPopupController()
		popupGroup := systemGroup.Group("/popup")
		{
			popupGroup.GET("/list", middleware.WithSelect("弹窗管理", "查询弹窗列表"), popupController.List)
			popupGroup.GET("/:id", middleware.WithSelect("弹窗管理", "查询弹窗详情"), popupController.GetInfo)
			popupGroup.POST("", middleware.WithCreate("弹窗管理", "新增弹窗"), popupController.Add)
			popupGroup.PUT("", middleware.WithUpdate("弹窗管理", "修改弹窗"), popupController.Update)
			popupGroup.DELETE("/:id", middleware.WithDelete("弹窗管理", "删除弹窗"), popupController.Delete)
			popupGroup.PUT("/status", middleware.WithUpdate("弹窗管理", "修改弹窗状态"), popupController.UpdateStatus)
			popupGroup.PUT("/stats", middleware.WithUpdate("弹窗管理", "更新弹窗统计"), popupController.UpdateStats)
			popupGroup.GET("/active", middleware.WithSelect("弹窗管理", "查询活跃弹窗"), popupController.GetActivePopups)
		}

		// 操作日志管理
		operationLogController := system.NewOperationLogController()
		operationLogGroup := systemGroup.Group("/operation-log")
		{
			operationLogGroup.GET("/list", middleware.WithSelect("操作日志", "查询操作日志"), operationLogController.List)
			operationLogGroup.DELETE("/delete/:id", middleware.WithDelete("操作日志", "删除操作日志"), operationLogController.Delete)
			operationLogGroup.DELETE("/batch-delete", middleware.WithDelete("操作日志", "批量删除操作日志"), operationLogController.BatchDelete)
			operationLogGroup.DELETE("/clear", middleware.WithDelete("操作日志", "清空操作日志"), operationLogController.Clear)
			operationLogGroup.GET("/export", middleware.WithSelect("操作日志", "导出操作日志"), operationLogController.Export)
		}
	}
}
