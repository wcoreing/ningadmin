package shop

import (
	"strconv"
	"wnadmin/app/common/response"
	shopDto "wnadmin/app/model/dto/shop"
	"wnadmin/app/service/shop/feedback"

	"github.com/gin-gonic/gin"
)

// FeedbackController 反馈控制器
type FeedbackController struct {
	feedbackService *feedback.Service
}

// NewFeedbackController 创建反馈控制器
func NewFeedbackController() *FeedbackController {
	return &FeedbackController{
		feedbackService: feedback.NewService(),
	}
}

// List 获取反馈列表
// @Summary 获取反馈列表
// @Description 获取用户反馈列表，支持分页和条件查询
// @Tags 反馈管理
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Param pageNum query int false "页码"
// @Param pageSize query int false "每页数量"
// @Param keyword query string false "关键字"
// @Param userId query int false "用户ID"
// @Param status query int false "状态(0:待处理 1:已处理 2:已关闭)"
// @Success 200 {object} response.Response "成功"
// @Router /shop/feedback/list [get]
func (c *FeedbackController) List(ctx *gin.Context) {
	var query shopDto.FeedbackQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		response.FailWithMessage(201, "参数错误: "+err.Error(), ctx)
		return
	}

	list, total, err := c.feedbackService.GetFeedbackList(&query)
	if err != nil {
		response.FailWithMessage(201, "获取反馈列表失败: "+err.Error(), ctx)
		return
	}

	response.OkWithDetailed(gin.H{
		"list":  list,
		"total": total,
	}, "获取成功", ctx)
}

// Add 新增反馈
// @Summary 新增反馈
// @Description 创建用户反馈
// @Tags 反馈管理
// @Accept json
// @Produce json
// @Param data body shopDto.CreateFeedbackReq true "反馈信息"
// @Success 200 {object} response.Response "成功"
// @Router /shop/feedback [post]
func (c *FeedbackController) Add(ctx *gin.Context) {
	var req shopDto.CreateFeedbackReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(201, "参数错误: "+err.Error(), ctx)
		return
	}

	if err := c.feedbackService.CreateFeedback(&req); err != nil {
		response.FailWithMessage(201, "创建反馈失败: "+err.Error(), ctx)
		return
	}

	response.OkWithMessage("创建成功", ctx)
}

// Update 更新反馈
// @Summary 更新反馈
// @Description 更新反馈信息
// @Tags 反馈管理
// @Accept json
// @Produce json
// @Param data body shopDto.UpdateFeedbackReq true "反馈信息"
// @Success 200 {object} response.Response "成功"
// @Router /shop/feedback [put]
func (c *FeedbackController) Update(ctx *gin.Context) {
	var req shopDto.UpdateFeedbackReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(201, "参数错误: "+err.Error(), ctx)
		return
	}

	if err := c.feedbackService.UpdateFeedback(&req); err != nil {
		response.FailWithMessage(201, "更新反馈失败: "+err.Error(), ctx)
		return
	}

	response.OkWithMessage("更新成功", ctx)
}

// Delete 删除反馈
// @Summary 删除反馈
// @Description 删除反馈
// @Tags 反馈管理
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Param id path int true "反馈ID"
// @Success 200 {object} response.Response "成功"
// @Router /shop/feedback/{id} [delete]
func (c *FeedbackController) Delete(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.FailWithMessage(201, "无效的ID", ctx)
		return
	}

	if err := c.feedbackService.DeleteFeedback(id); err != nil {
		response.FailWithMessage(201, "删除反馈失败: "+err.Error(), ctx)
		return
	}

	response.OkWithMessage("删除成功", ctx)
}

// GetInfo 获取反馈详情
// @Summary 获取反馈详情
// @Description 获取反馈详细信息
// @Tags 反馈管理
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Param id path int true "反馈ID"
// @Success 200 {object} response.Response "成功"
// @Router /shop/feedback/{id} [get]
func (c *FeedbackController) GetInfo(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.FailWithMessage(201, "无效的ID", ctx)
		return
	}

	feedbackInfo, err := c.feedbackService.GetFeedbackInfo(id)
	if err != nil {
		response.FailWithMessage(201, "获取反馈详情失败: "+err.Error(), ctx)
		return
	}

	response.OkWithDetailed(feedbackInfo, "获取成功", ctx)
}

// Reply 回复反馈
// @Summary 回复反馈
// @Description 管理员回复用户反馈
// @Tags 反馈管理
// @Accept json
// @Produce json
// @Param id path int true "反馈ID"
// @Param data body gin.H true "回复内容"
// @Success 200 {object} response.Response "成功"
// @Router /shop/feedback/{id}/reply [post]
func (c *FeedbackController) Reply(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.FailWithMessage(201, "无效的ID", ctx)
		return
	}

	var body struct {
		Reply string `json:"reply" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		response.FailWithMessage(201, "参数错误: "+err.Error(), ctx)
		return
	}

	if err := c.feedbackService.ReplyFeedback(id, body.Reply); err != nil {
		response.FailWithMessage(201, "回复反馈失败: "+err.Error(), ctx)
		return
	}

	response.OkWithMessage("回复成功", ctx)
}

// UpdateStatus 更新反馈状态
// @Summary 更新反馈状态
// @Description 更新反馈处理状态
// @Tags 反馈管理
// @Accept json
// @Produce json
// @Param id path int true "反馈ID"
// @Param data body gin.H true "状态信息"
// @Success 200 {object} response.Response "成功"
// @Router /shop/feedback/{id}/status [put]
func (c *FeedbackController) UpdateStatus(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.FailWithMessage(201, "无效的ID", ctx)
		return
	}

	var body struct {
		Status int `json:"status" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		response.FailWithMessage(201, "参数错误: "+err.Error(), ctx)
		return
	}

	if err := c.feedbackService.UpdateFeedbackStatus(id, body.Status); err != nil {
		response.FailWithMessage(201, "更新状态失败: "+err.Error(), ctx)
		return
	}

	response.OkWithMessage("更新成功", ctx)
}
