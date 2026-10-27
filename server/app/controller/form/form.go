package form

import (
	"wnadmin/app/common/response"
	"wnadmin/app/common/utils"
	formdo "wnadmin/app/model/do/form"
	"wnadmin/app/model/dto"
	formsvc "wnadmin/app/service/form"

	"github.com/gin-gonic/gin"
)

type FormListController struct {
	formListService *formsvc.FormListService
}

func NewFormListController() *FormListController {
	return &FormListController{
		formListService: formsvc.NewFormListService(),
	}
}

// GetFormList 获取表单列表
func (c *FormListController) GetFormList(ctx *gin.Context) {
	var query dto.FormListQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		response.FailWithMessage(response.ERROR, err.Error(), ctx)
		return
	}

	forms, total, err := c.formListService.GetFormList(&query)
	if err != nil {
		response.FailWithMessage(response.ERROR, err.Error(), ctx)
		return
	}

	response.OkWithData(gin.H{
		"total": total,
		"list":  forms,
	}, ctx)
}

// GetFormDetail 获取表单详情
func (c *FormListController) GetFormDetail(ctx *gin.Context) {
	id := utils.ParseParamID(ctx)
	if id <= 0 {
		response.FailWithMessage(response.ERROR, "无效的参数", ctx)
		return
	}

	form, err := c.formListService.GetFormDetail(id)
	if err != nil {
		response.FailWithMessage(response.ERROR, err.Error(), ctx)
		return
	}

	response.OkWithData(form, ctx)
}

// CreateForm 创建表单
func (c *FormListController) CreateForm(ctx *gin.Context) {
	var req dto.FormCreateDTO
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(response.ERROR, err.Error(), ctx)
		return
	}

	form := &formdo.Form{
		TemplateID:      req.TemplateID,
		FormName:        req.FormName,
		FormDescription: req.FormDescription,
		FormAttributes:  req.FormAttributes,
		FormStatus:      req.FormStatus,
	}

	username := utils.GetCurrentUsername(ctx)
	savedForm, err := c.formListService.CreateForm(form, username)
	if err != nil {
		response.FailWithMessage(response.ERROR, err.Error(), ctx)
		return
	}

	response.OkWithData(gin.H{
		"id": savedForm.ID,
	}, ctx)
}

// UpdateForm 更新表单
func (c *FormListController) UpdateForm(ctx *gin.Context) {
	var req dto.FormUpdateDTO
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(response.ERROR, err.Error(), ctx)
		return
	}

	if req.ID == 0 {
		response.FailWithMessage(response.ERROR, "ID不能为空", ctx)
		return
	}

	form := &formdo.Form{
		ID:              req.ID,
		TemplateID:      req.TemplateID,
		FormName:        req.FormName,
		FormDescription: req.FormDescription,
		FormAttributes:  req.FormAttributes,
		FormStatus:      req.FormStatus,
	}

	username := utils.GetCurrentUsername(ctx)
	savedForm, err := c.formListService.UpdateForm(form, username)
	if err != nil {
		response.FailWithMessage(response.ERROR, err.Error(), ctx)
		return
	}

	response.OkWithData(gin.H{
		"id": savedForm.ID,
	}, ctx)
}

// DeleteForm 删除表单
func (c *FormListController) DeleteForm(ctx *gin.Context) {
	id := utils.ParseParamID(ctx)
	if id <= 0 {
		response.FailWithMessage(response.ERROR, "无效的参数", ctx)
		return
	}

	username := utils.GetCurrentUsername(ctx)
	err := c.formListService.DeleteForm(id, username)
	if err != nil {
		response.FailWithMessage(response.ERROR, err.Error(), ctx)
		return
	}

	response.Ok(ctx)
}
