package form

import (
	"strconv"
	"wnadmin/app/common/response"
	"wnadmin/app/common/utils"
	"wnadmin/app/model/dto"
	"wnadmin/app/service/form"

	"github.com/gin-gonic/gin"
)

// FormDataController 表单数据控制器
type FormDataController struct {
	formDataService *form.FormDataService
}

// NewFormDataController 创建表单数据控制器实例
func NewFormDataController() *FormDataController {
	return &FormDataController{
		formDataService: form.NewFormDataService(),
	}
}

// GetFormDataList 获取表单数据列表
func (c *FormDataController) GetFormDataList(ctx *gin.Context) {
	var query dto.FormDataQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		response.FailWithMessage(response.ERROR, err.Error(), ctx)
		return
	}

	list, total, err := c.formDataService.GetFormDataList(&query)
	if err != nil {
		response.FailWithMessage(response.ERROR, err.Error(), ctx)
		return
	}

	response.OkWithData(gin.H{
		"total": total,
		"list":  list,
	}, ctx)
}

// GetFormDataDetail 获取表单数据详情
func (c *FormDataController) GetFormDataDetail(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		response.FailWithMessage(response.ERROR, "无效的ID", ctx)
		return
	}

	data, err := c.formDataService.GetFormDataDetail(id)
	if err != nil {
		response.FailWithMessage(response.ERROR, err.Error(), ctx)
		return
	}

	response.OkWithData(data, ctx)
}

// DeleteFormData 删除表单数据
func (c *FormDataController) DeleteFormData(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		response.FailWithMessage(response.ERROR, "无效的ID", ctx)
		return
	}

	err = c.formDataService.DeleteFormData(id)
	if err != nil {
		response.FailWithMessage(response.ERROR, err.Error(), ctx)
		return
	}

	response.Ok(ctx)
}

// SubmitForm 提交表单
func (c *FormDataController) SubmitForm(ctx *gin.Context) {
	var req dto.FormSubmitRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(response.ERROR, err.Error(), ctx)
		return
	}

	// 获取IP地址
	ip := ctx.ClientIP()
	// 获取当前用户
	username := utils.GetCurrentUsername(ctx)

	err := c.formDataService.SubmitForm(&req, ip, username)
	if err != nil {
		response.FailWithMessage(response.ERROR, err.Error(), ctx)
		return
	}

	response.Ok(ctx)
}
