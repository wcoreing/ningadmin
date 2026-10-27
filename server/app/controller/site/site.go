package site

import (
	"strconv"
	"wnadmin/app/common/response"
	"wnadmin/app/model/dto"
	siteservice "wnadmin/app/service/site"

	"github.com/gin-gonic/gin"
)

// Controller 官网接口
type Controller struct {
	service *siteservice.Service
}

// NewController 创建官网控制器
func NewController() *Controller {
	return &Controller{service: &siteservice.Service{}}
}

// Home 首页
func (c *Controller) Home(ctx *gin.Context) {
	data, err := c.service.Home()
	if err != nil {
		response.FailWithMessage(500, err.Error(), ctx)
		return
	}
	response.OkWithData(data, ctx)
}

// Profile 企业资料
func (c *Controller) Profile(ctx *gin.Context) {
	data, err := c.service.Profile()
	if err != nil {
		response.FailWithMessage(500, err.Error(), ctx)
		return
	}
	response.OkWithData(data, ctx)
}

// Page 栏目
func (c *Controller) Page(ctx *gin.Context) {
	data, err := c.service.Page(ctx.Param("code"))
	if err != nil {
		response.FailWithMessage(404, err.Error(), ctx)
		return
	}
	response.OkWithData(data, ctx)
}

// Categories 分类
func (c *Controller) Categories(ctx *gin.Context) {
	data, err := c.service.Categories()
	if err != nil {
		response.FailWithMessage(500, err.Error(), ctx)
		return
	}
	response.OkWithData(data, ctx)
}

// Goods 产品列表
func (c *Controller) Goods(ctx *gin.Context) {
	var query dto.SiteGoodsQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		response.FailWithMessage(400, err.Error(), ctx)
		return
	}
	data, err := c.service.Goods(&query)
	if err != nil {
		response.FailWithMessage(500, err.Error(), ctx)
		return
	}
	response.OkWithData(data, ctx)
}

// GoodsDetail 产品详情
func (c *Controller) GoodsDetail(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		response.FailWithMessage(400, "产品编号不正确", ctx)
		return
	}
	data, err := c.service.GoodsDetail(id)
	if err != nil {
		response.FailWithMessage(404, err.Error(), ctx)
		return
	}
	response.OkWithData(data, ctx)
}

// Articles 资讯列表
func (c *Controller) Articles(ctx *gin.Context) {
	var query dto.SiteArticleQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		response.FailWithMessage(400, err.Error(), ctx)
		return
	}
	data, err := c.service.Articles(&query)
	if err != nil {
		response.FailWithMessage(500, err.Error(), ctx)
		return
	}
	response.OkWithData(data, ctx)
}

// Article 资讯详情
func (c *Controller) Article(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		response.FailWithMessage(400, "资讯编号不正确", ctx)
		return
	}
	data, err := c.service.Article(id)
	if err != nil {
		response.FailWithMessage(404, err.Error(), ctx)
		return
	}
	response.OkWithData(data, ctx)
}

// CreateInquiry 提交询盘
func (c *Controller) CreateInquiry(ctx *gin.Context) {
	var req dto.SiteInquiryReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(400, "请完整填写询盘信息", ctx)
		return
	}
	if err := c.service.CreateInquiry(&req); err != nil {
		response.FailWithMessage(400, err.Error(), ctx)
		return
	}
	response.OkWithMessage("询盘已提交", ctx)
}

// Inquiries 询盘列表
func (c *Controller) Inquiries(ctx *gin.Context) {
	var query dto.SiteInquiryQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		response.FailWithMessage(400, err.Error(), ctx)
		return
	}
	data, err := c.service.Inquiries(&query)
	if err != nil {
		response.FailWithMessage(500, err.Error(), ctx)
		return
	}
	response.OkWithData(data, ctx)
}

// UpdateInquiryStatus 更新询盘状态
func (c *Controller) UpdateInquiryStatus(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		response.FailWithMessage(400, "询盘编号不正确", ctx)
		return
	}
	var req dto.SiteInquiryStatusReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(400, "状态不正确", ctx)
		return
	}
	if err := c.service.UpdateInquiryStatus(id, req.Status); err != nil {
		response.FailWithMessage(400, err.Error(), ctx)
		return
	}
	response.OkWithMessage("已更新", ctx)
}

// SaveProfile 保存企业资料
func (c *Controller) SaveProfile(ctx *gin.Context) {
	var req dto.SiteProfile
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(400, "请完整填写企业资料", ctx)
		return
	}
	if err := c.service.SaveProfile(&req); err != nil {
		response.FailWithMessage(400, err.Error(), ctx)
		return
	}
	response.OkWithMessage("已保存", ctx)
}
