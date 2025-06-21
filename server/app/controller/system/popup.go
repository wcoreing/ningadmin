package system

import (
	"strconv"
	"wnadmin/app/common/response"
	"wnadmin/app/model/dto"
	"wnadmin/app/service/system"

	"github.com/gin-gonic/gin"
)

// PopupController 弹窗控制器
type PopupController struct {
	popupService *system.PopupService
}

// NewPopupController 创建弹窗控制器
func NewPopupController() *PopupController {
	return &PopupController{
		popupService: &system.PopupService{},
	}
}

// List 获取弹窗列表
// @Summary 获取弹窗列表
// @Description 获取弹窗列表，支持分页和条件查询
// @Tags 弹窗管理
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Param title query string false "标题"
// @Param trigger_type query string false "触发类型"
// @Param status query int false "状态"
// @Success 200 {object} response.Response "成功"
// @Router /system/popup/list [get]
func (c *PopupController) List(ctx *gin.Context) {
	var query dto.PopupQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		response.FailWithMessage(201, "参数错误: "+err.Error(), ctx)
		return
	}

	// 设置分页默认值
	if query.Page <= 0 {
		query.Page = 1
	}
	if query.PageSize <= 0 {
		query.PageSize = 10
	}

	list, total, err := c.popupService.GetPopupList(&query)
	if err != nil {
		response.FailWithMessage(201, "获取弹窗列表失败: "+err.Error(), ctx)
		return
	}

	response.OkWithDetailed(gin.H{
		"list":  list,
		"total": total,
	}, "获取成功", ctx)
}

// Add 新增弹窗
// @Summary 新增弹窗
// @Description 新增弹窗
// @Tags 弹窗管理
// @Accept json
// @Produce json
// @Param data body dto.CreatePopupReq true "弹窗信息"
// @Success 200 {object} response.Response "成功"
// @Router /system/popup [post]
func (c *PopupController) Add(ctx *gin.Context) {
	var req dto.CreatePopupReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(201, "参数错误: "+err.Error(), ctx)
		return
	}

	if err := c.popupService.CreatePopup(&req); err != nil {
		response.FailWithMessage(201, "创建弹窗失败: "+err.Error(), ctx)
		return
	}

	response.OkWithMessage("创建成功", ctx)
}

// Update 更新弹窗
// @Summary 更新弹窗
// @Description 更新弹窗
// @Tags 弹窗管理
// @Accept json
// @Produce json
// @Param data body dto.UpdatePopupReq true "弹窗信息"
// @Success 200 {object} response.Response "成功"
// @Router /system/popup [put]
func (c *PopupController) Update(ctx *gin.Context) {
	var req dto.UpdatePopupReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(201, "参数错误: "+err.Error(), ctx)
		return
	}

	if err := c.popupService.UpdatePopup(&req); err != nil {
		response.FailWithMessage(201, "更新弹窗失败: "+err.Error(), ctx)
		return
	}

	response.OkWithMessage("更新成功", ctx)
}

// Delete 删除弹窗
// @Summary 删除弹窗
// @Description 删除弹窗
// @Tags 弹窗管理
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Param id path int true "弹窗ID"
// @Success 200 {object} response.Response "成功"
// @Router /system/popup/{id} [delete]
func (c *PopupController) Delete(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.FailWithMessage(201, "无效的ID", ctx)
		return
	}

	if err := c.popupService.DeletePopup(id); err != nil {
		response.FailWithMessage(201, "删除弹窗失败: "+err.Error(), ctx)
		return
	}

	response.OkWithMessage("删除成功", ctx)
}

// GetInfo 获取弹窗详情
// @Summary 获取弹窗详情
// @Description 获取弹窗详情
// @Tags 弹窗管理
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Param id path int true "弹窗ID"
// @Success 200 {object} response.Response "成功"
// @Router /system/popup/{id} [get]
func (c *PopupController) GetInfo(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.FailWithMessage(201, "无效的ID", ctx)
		return
	}

	popup, err := c.popupService.GetPopupInfo(id)
	if err != nil {
		response.FailWithMessage(201, "获取弹窗详情失败: "+err.Error(), ctx)
		return
	}

	response.OkWithDetailed(popup, "获取成功", ctx)
}

// UpdateStatus 更新弹窗状态
// @Summary 更新弹窗状态
// @Description 更新弹窗状态
// @Tags 弹窗管理
// @Accept json
// @Produce json
// @Param data body dto.PopupStatusReq true "状态信息"
// @Success 200 {object} response.Response "成功"
// @Router /system/popup/status [put]
func (c *PopupController) UpdateStatus(ctx *gin.Context) {
	var req dto.PopupStatusReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(201, "参数错误: "+err.Error(), ctx)
		return
	}

	if err := c.popupService.UpdatePopupStatus(&req); err != nil {
		response.FailWithMessage(201, "更新状态失败: "+err.Error(), ctx)
		return
	}

	response.OkWithMessage("更新成功", ctx)
}

// UpdateStats 更新弹窗统计
// @Summary 更新弹窗统计
// @Description 更新弹窗统计（显示次数、点击次数）
// @Tags 弹窗管理
// @Accept json
// @Produce json
// @Param data body dto.PopupStatsReq true "统计信息"
// @Success 200 {object} response.Response "成功"
// @Router /system/popup/stats [put]
func (c *PopupController) UpdateStats(ctx *gin.Context) {
	var req dto.PopupStatsReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(201, "参数错误: "+err.Error(), ctx)
		return
	}

	if err := c.popupService.UpdatePopupStats(&req); err != nil {
		response.FailWithMessage(201, "更新统计失败: "+err.Error(), ctx)
		return
	}

	response.OkWithMessage("更新成功", ctx)
}

// GetActivePopups 获取活跃弹窗
// @Summary 获取活跃弹窗
// @Description 获取当前时间范围内的活跃弹窗
// @Tags 弹窗管理
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Param trigger_type query string false "触发类型"
// @Param user_id query int false "用户ID"
// @Success 200 {object} response.Response "成功"
// @Router /system/popup/active [get]
func (c *PopupController) GetActivePopups(ctx *gin.Context) {
	var req dto.GetActivePopupsReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(201, "参数错误: "+err.Error(), ctx)
		return
	}

	list, err := c.popupService.GetActivePopups(&req)
	if err != nil {
		response.FailWithMessage(201, "获取活跃弹窗失败: "+err.Error(), ctx)
		return
	}

	response.OkWithDetailed(gin.H{
		"list": list,
	}, "获取成功", ctx)
}
