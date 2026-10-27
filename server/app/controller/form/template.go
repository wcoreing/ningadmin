package form

import (
	"encoding/json"
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
	var req dto.FormTemplateDTO
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(response.ERROR, err.Error(), ctx)
		return
	}

	// 转换为DO对象
	template := &formdo.FormTemplate{
		ID:             req.ID,
		Name:           req.Name,
		Description:    req.Description,
		SubmitURL:      req.SubmitURL,
		SubmitMethod:   req.SubmitMethod,
		SuccessMessage: req.SuccessMessage,
		FailMessage:    req.FailMessage,
		FormItems:      req.FormItems,
		Config:         req.Config,
	}

	// 获取当前用户
	username := utils.GetCurrentUsername(ctx)

	// 调用 service 层处理业务逻辑
	savedTemplate, err := c.formService.SaveTemplate(template, username)
	if err != nil {
		response.FailWithMessage(response.ERROR, err.Error(), ctx)
		return
	}

	response.OkWithData(gin.H{
		"id": savedTemplate.ID,
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

// CreateTemplate 创建表单模板
func (c *FormController) CreateTemplate(ctx *gin.Context) {
	var req dto.FormTemplateCreateDTO
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(response.ERROR, err.Error(), ctx)
		return
	}

	template := &formdo.FormTemplate{
		Name:           req.Name,
		Description:    req.Description,
		SubmitURL:      req.SubmitURL,
		SubmitMethod:   req.SubmitMethod,
		SuccessMessage: req.SuccessMessage,
		FailMessage:    req.FailMessage,
		FormItems:      req.FormItems,
		Config:         req.Config,
	}

	username := utils.GetCurrentUsername(ctx)
	savedTemplate, err := c.formService.CreateTemplate(template, username)
	if err != nil {
		response.FailWithMessage(response.ERROR, err.Error(), ctx)
		return
	}

	response.OkWithData(gin.H{
		"id": savedTemplate.ID,
	}, ctx)
}

// UpdateTemplate 更新表单模板
func (c *FormController) UpdateTemplate(ctx *gin.Context) {
	var req dto.FormTemplateDTO
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(response.ERROR, err.Error(), ctx)
		return
	}

	if req.ID == 0 {
		response.FailWithMessage(response.ERROR, "ID不能为空", ctx)
		return
	}

	template := &formdo.FormTemplate{
		ID:             req.ID,
		Name:           req.Name,
		Description:    req.Description,
		SubmitURL:      req.SubmitURL,
		SubmitMethod:   req.SubmitMethod,
		SuccessMessage: req.SuccessMessage,
		FailMessage:    req.FailMessage,
		FormItems:      req.FormItems,
		Config:         req.Config,
	}

	username := utils.GetCurrentUsername(ctx)
	savedTemplate, err := c.formService.UpdateTemplate(template, username)
	if err != nil {
		response.FailWithMessage(response.ERROR, err.Error(), ctx)
		return
	}

	response.OkWithData(gin.H{
		"id": savedTemplate.ID,
	}, ctx)
}

// GetTemplateInfo 获取表单模板详情
func (c *FormController) GetTemplateInfo(ctx *gin.Context) {
	id := utils.ParseParamID(ctx)
	if id <= 0 {
		response.FailWithMessage(response.ERROR, "无效的参数", ctx)
		return
	}

	template, err := c.formService.GetTemplateInfo(id)
	if err != nil {
		response.FailWithMessage(response.ERROR, err.Error(), ctx)
		return
	}

	// 解析 form_items 和 config
	var formItems interface{}
	var config interface{}
	if template.FormItems != nil {
		if err := json.Unmarshal(template.FormItems, &formItems); err != nil {
			response.FailWithMessage(response.ERROR, "解析表单项失败", ctx)
			return
		}
	}
	if template.Config != nil {
		if err := json.Unmarshal(template.Config, &config); err != nil {
			response.FailWithMessage(response.ERROR, "解析配置失败", ctx)
			return
		}
	}

	response.OkWithData(gin.H{
		"id":             template.ID,
		"name":           template.Name,
		"description":    template.Description,
		"submitUrl":      template.SubmitURL,
		"submitMethod":   template.SubmitMethod,
		"successMessage": template.SuccessMessage,
		"failMessage":    template.FailMessage,
		"formItems":      formItems,
		"config":         config,
		"status":         template.Status,
		"createBy":       template.CreateBy,
		"createTime":     template.CreateTime,
		"updateBy":       template.UpdateBy,
		"updateTime":     template.UpdateTime,
	}, ctx)
}
