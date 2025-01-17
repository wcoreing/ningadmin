package system

import (
	"strconv"
	"wnadmin/app/common/response"
	"wnadmin/app/model/dto"
	"wnadmin/app/service/system"

	"github.com/gin-gonic/gin"
)

// ConfigController 配置控制器
type ConfigController struct {
	configService *system.ConfigService
}

// NewConfigController 创建配置控制器
func NewConfigController() *ConfigController {
	return &ConfigController{
		configService: &system.ConfigService{},
	}
}

// List 获取配置列表
// @Summary 获取配置列表
// @Description 获取配置列表，支持分页和条件查询
// @Tags 系统配置
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Param pageNum query int false "页码"
// @Param pageSize query int false "每页数量"
// @Param keyword query string false "关键字"
// @Param moduleKey query string false "模块标识"
// @Success 200 {object} response.Response "成功"
// @Router /system/config/list [get]
func (c *ConfigController) List(ctx *gin.Context) {
	var query dto.ConfigQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		response.FailWithMessage(201, "参数错误: "+err.Error(), ctx)
		return
	}

	// 设置默认值
	if query.PageNum <= 0 {
		query.PageNum = 1
	}
	if query.PageSize <= 0 {
		query.PageSize = 10
	}

	list, total, err := c.configService.GetConfigList(&query)
	if err != nil {
		response.FailWithMessage(201, "获取配置列表失败: "+err.Error(), ctx)
		return
	}

	response.OkWithDetailed(gin.H{
		"list":  list,
		"total": total,
	}, "获取成功", ctx)
}

// Add 新增配置
// @Summary 新增配置
// @Description 新增系统配置
// @Tags 系统配置
// @Accept json
// @Produce json
// @Param data body dto.CreateConfigReq true "配置信息"
// @Success 200 {object} response.Response "成功"
// @Router /system/config [post]
func (c *ConfigController) Add(ctx *gin.Context) {
	var req dto.CreateConfigReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(201, "参数错误: "+err.Error(), ctx)
		return
	}

	// 从上下文获取创建人
	createBy := ctx.GetString("userName")
	if err := c.configService.CreateConfig(&req, createBy); err != nil {
		response.FailWithMessage(201, "创建配置失败: "+err.Error(), ctx)
		return
	}

	response.OkWithMessage("创建成功", ctx)
}

// Update 更新配置
// @Summary 更新配置
// @Description 更新系统配置
// @Tags 系统配置
// @Accept json
// @Produce json
// @Param data body dto.UpdateConfigReq true "配置信息"
// @Success 200 {object} response.Response "成功"
// @Router /system/config [put]
func (c *ConfigController) Update(ctx *gin.Context) {
	var req dto.UpdateConfigReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(201, "参数错误: "+err.Error(), ctx)
		return
	}

	// 从上下文获取更新人
	updateBy := ctx.GetString("userName")
	if err := c.configService.UpdateConfig(&req, updateBy); err != nil {
		response.FailWithMessage(201, "更新配置失败: "+err.Error(), ctx)
		return
	}

	response.OkWithMessage("更新成功", ctx)
}

// Delete 删除配置
// @Summary 删除配置
// @Description 删除系统配置
// @Tags 系统配置
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Param id path int true "配置ID"
// @Success 200 {object} response.Response "成功"
// @Router /system/config/{id} [delete]
func (c *ConfigController) Delete(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.FailWithMessage(201, "无效的ID", ctx)
		return
	}

	if err := c.configService.DeleteConfig(id); err != nil {
		response.FailWithMessage(201, "删除配置失败: "+err.Error(), ctx)
		return
	}

	response.OkWithMessage("删除成功", ctx)
}

// GetInfo 获取配置详情
// @Summary 获取配置详情
// @Description 获取系统配置详情
// @Tags 系统配置
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Param id path int true "配置ID"
// @Success 200 {object} response.Response "成功"
// @Router /system/config/{id} [get]
func (c *ConfigController) GetInfo(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.FailWithMessage(201, "无效的ID", ctx)
		return
	}

	config, err := c.configService.GetConfigInfo(id)
	if err != nil {
		response.FailWithMessage(201, "获取配置详情失败: "+err.Error(), ctx)
		return
	}

	response.OkWithDetailed(config, "获取成功", ctx)
}
