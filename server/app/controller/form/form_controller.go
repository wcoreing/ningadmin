package form

import (
	"wnadmin/app/common/response"
	"wnadmin/app/common/utils"
	formdo "wnadmin/app/model/do/form"
	"wnadmin/app/model/dto"
	formsvc "wnadmin/app/service/form"

	"github.com/gin-gonic/gin"
)

type FormController struct {
	formService *formsvc.FormService
}

func NewFormController() *FormController {
	return &FormController{
		formService: formsvc.NewFormService(),
	}
}

// GetTemplateList 获取表单模板列表
func (c *FormController) GetTemplateList(ctx *gin.Context) {
	var query dto.FormTemplateQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		response.FailWithMessage(response.ERROR, err.Error(), ctx)
		return
	}

	templates, total, err := c.formService.GetTemplateList(&query)
	if err != nil {
		response.FailWithMessage(response.ERROR, err.Error(), ctx)
		return
	}

	response.OkWithData(gin.H{
		"total": total,
		"list":  templates,
	}, ctx)
}

// SaveTemplate 保存表单模板
func (c *FormController) SaveTemplate(ctx *gin.Context) {
	var template formdo.FormTemplate
	if err := ctx.ShouldBindJSON(&template); err != nil {
		response.FailWithMessage(response.ERROR, err.Error(), ctx)
		return
	}

	// 设置创建者/更新者
	if template.ID > 0 {
		template.UpdateBy = utils.GetCurrentUsername(ctx)
	} else {
		template.CreateBy = utils.GetCurrentUsername(ctx)
	}

	err := c.formService.SaveTemplate(&template)
	if err != nil {
		response.FailWithMessage(response.ERROR, err.Error(), ctx)
		return
	}

	response.OkWithData(gin.H{
		"id": template.ID,
	}, ctx)
}

// DeleteTemplate 删除表单模板
func (c *FormController) DeleteTemplate(ctx *gin.Context) {
	id := utils.ParseParamID(ctx)
	if id <= 0 {
		response.FailWithMessage(response.ERROR, "无效的参数", ctx)
		return
	}

	err := c.formService.DeleteTemplate(id)
	if err != nil {
		response.FailWithMessage(response.ERROR, err.Error(), ctx)
		return
	}

	response.Ok(ctx)
}

// CreateShare 创建表单分享
func (c *FormController) CreateShare(ctx *gin.Context) {
	var req dto.FormShareCreateDTO
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(response.ERROR, err.Error(), ctx)
		return
	}

	share := &formdo.FormShare{
		TemplateID: req.TemplateID,
		ExpireTime: req.ExpireTime,
		CreateBy:   utils.GetCurrentUsername(ctx),
	}

	err := c.formService.CreateShare(share)
	if err != nil {
		response.FailWithMessage(response.ERROR, err.Error(), ctx)
		return
	}

	response.OkWithData(&dto.FormShareDTO{
		ShareCode: share.ShareCode,
		ShareURL:  utils.GetBaseURL(ctx) + "/form/share/" + share.ShareCode,
	}, ctx)
}

// GetShareInfo 获取分享表单信息
func (c *FormController) GetShareInfo(ctx *gin.Context) {
	shareCode := ctx.Param("shareCode")
	if shareCode == "" {
		response.FailWithMessage(response.ERROR, "无效的参数", ctx)
		return
	}

	template, err := c.formService.GetShareInfo(shareCode)
	if err != nil {
		response.FailWithMessage(response.ERROR, err.Error(), ctx)
		return
	}

	response.OkWithData(gin.H{
		"templateId": template.ID,
		"formItems":  template.FormItems,
		"config":     template.Config,
	}, ctx)
}

// SubmitForm 提交表单
func (c *FormController) SubmitForm(ctx *gin.Context) {
	var req dto.FormSubmissionDTO
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(response.ERROR, err.Error(), ctx)
		return
	}

	submission := &formdo.FormSubmission{
		TemplateID: req.TemplateID,
		FormData:   req.FormData,
		SubmitBy:   utils.GetCurrentUsername(ctx),
	}

	err := c.formService.SubmitForm(submission)
	if err != nil {
		response.FailWithMessage(response.ERROR, err.Error(), ctx)
		return
	}

	response.OkWithData(gin.H{
		"id": submission.ID,
	}, ctx)
}

// GetSubmissionList 获取表单提交记录
func (c *FormController) GetSubmissionList(ctx *gin.Context) {
	var query dto.FormSubmissionQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		response.FailWithMessage(response.ERROR, err.Error(), ctx)
		return
	}

	submissions, total, err := c.formService.GetSubmissionList(&query)
	if err != nil {
		response.FailWithMessage(response.ERROR, err.Error(), ctx)
		return
	}

	response.OkWithData(gin.H{
		"total": total,
		"list":  submissions,
	}, ctx)
}
