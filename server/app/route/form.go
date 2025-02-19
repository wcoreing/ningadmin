package route

import (
	"wnadmin/app/controller/form"

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
	formFormGroup := formGroup.Group("/form")
	{
		formFormGroup.GET("/list", formListController.GetFormList)
		formFormGroup.GET("/:id", formListController.GetFormDetail)
		formFormGroup.POST("", formListController.CreateForm)
		formFormGroup.PUT("", formListController.UpdateForm)
		formFormGroup.DELETE("/:id", formListController.DeleteForm)
	}

	// 表单数据管理
	formDataGroup := formGroup.Group("/data")
	{
		formDataGroup.GET("/list", formDataController.GetFormDataList)
		formDataGroup.GET("/:id", formDataController.GetFormDataDetail)
		formDataGroup.DELETE("/:id", formDataController.DeleteFormData)
	}
}
