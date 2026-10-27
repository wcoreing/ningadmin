package cms

import (
	"strconv"
	"wnadmin/app/common/response"
	"wnadmin/app/model/dto"
	cmsservice "wnadmin/app/service/cms"

	"github.com/gin-gonic/gin"
)

// Controller 内容发布
type Controller struct {
	service *cmsservice.Service
}

// NewController 创建控制器
func NewController() *Controller {
	return &Controller{service: &cmsservice.Service{}}
}

// ContentList 内容列表
func (c *Controller) ContentList(ctx *gin.Context) {
	var query dto.CmsContentQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		response.FailWithMessage(201, "参数错误", ctx)
		return
	}
	list, total, err := c.service.ContentList(&query)
	if err != nil {
		response.FailWithMessage(201, err.Error(), ctx)
		return
	}
	response.OkWithData(gin.H{"list": list, "total": total}, ctx)
}

// ContentInfo 内容详情
func (c *Controller) ContentInfo(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		response.FailWithMessage(201, "内容编号不正确", ctx)
		return
	}
	data, err := c.service.ContentInfo(id)
	if err != nil {
		response.FailWithMessage(201, err.Error(), ctx)
		return
	}
	response.OkWithData(data, ctx)
}

// SaveContent 保存内容
func (c *Controller) SaveContent(ctx *gin.Context) {
	var req dto.CmsContentReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(201, "请选择栏目并填写标题", ctx)
		return
	}
	if err := c.service.SaveContent(&req); err != nil {
		response.FailWithMessage(201, err.Error(), ctx)
		return
	}
	response.OkWithMessage("保存成功", ctx)
}

// DeleteContent 删除内容
func (c *Controller) DeleteContent(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		response.FailWithMessage(201, "内容编号不正确", ctx)
		return
	}
	if err := c.service.DeleteContent(id); err != nil {
		response.FailWithMessage(201, err.Error(), ctx)
		return
	}
	response.OkWithMessage("删除成功", ctx)
}
