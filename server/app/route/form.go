package route

import (
	"wnadmin/app/controller/form"
	"wnadmin/app/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterFormRoutes(r *gin.RouterGroup) {
	formGroup := r.Group("/form")
	formController := form.NewFormController()
	formListController := form.NewFormListController()
	formDataController := form.NewFormDataController()

	// 表单模板管理
	formTemplateGroup := formGroup.Group("/template")
	{
		formTemplateGroup.GET("/list", middleware.WithSelect("表单模板", "查询模板列表"), formController.GetTemplateList)
		formTemplateGroup.POST("", middleware.WithCreate("表单模板", "新增模板"), formController.CreateTemplate)
		formTemplateGroup.PUT("", middleware.WithUpdate("表单模板", "修改模板"), formController.UpdateTemplate)
		formTemplateGroup.DELETE("/:id", middleware.WithDelete("表单模板", "删除模板"), formController.DeleteTemplate)
		formTemplateGroup.GET("/:id", middleware.WithSelect("表单模板", "查询模板详情"), formController.GetTemplateInfo)
	}

	// 表单分享管理
	formShareGroup := formGroup.Group("/share")
	{
		formShareGroup.POST("/create", middleware.WithCreate("表单分享", "创建分享"), formController.CreateShare)
	}

	// 表单管理
	formFormGroup := formGroup.Group("/form")
	{
		formFormGroup.GET("/list", middleware.WithSelect("表单管理", "查询表单列表"), formListController.GetFormList)
		formFormGroup.GET("/:id", middleware.WithSelect("表单管理", "查询表单详情"), formListController.GetFormDetail)
		formFormGroup.POST("", middleware.WithCreate("表单管理", "新增表单"), formListController.CreateForm)
		formFormGroup.PUT("", middleware.WithUpdate("表单管理", "修改表单"), formListController.UpdateForm)
		formFormGroup.DELETE("/:id", middleware.WithDelete("表单管理", "删除表单"), formListController.DeleteForm)
	}

	// 表单数据管理
	formDataGroup := formGroup.Group("/data")
	{
		formDataGroup.GET("/list", middleware.WithSelect("表单数据", "查询表单数据"), formDataController.GetFormDataList)
		formDataGroup.GET("/:id", middleware.WithSelect("表单数据", "查询数据详情"), formDataController.GetFormDataDetail)
		formDataGroup.DELETE("/:id", middleware.WithDelete("表单数据", "删除表单数据"), formDataController.DeleteFormData)
	}
}
