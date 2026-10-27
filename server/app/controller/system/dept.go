package system

import (
	"strconv"
	"wnadmin/app/common/response"
	"wnadmin/app/model/dto"
	"wnadmin/app/service/system"

	"github.com/gin-gonic/gin"
)

// DeptController 部门控制器
type DeptController struct {
	deptService *system.DeptService
}

// NewDeptController 创建部门控制器
func NewDeptController() *DeptController {
	return &DeptController{
		deptService: &system.DeptService{},
	}
}

// List 获取部门列表
// @Summary 获取部门列表
// @Description 获取部门列表，支持分页和条件查询
// @Tags 系统部门
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Param pageNum query int false "页码"
// @Param pageSize query int false "每页数量"
// @Param keyword query string false "关键字"
// @Param status query string false "状态"
// @Success 200 {object} response.Response "成功"
// @Router /system/dept/list [get]
func (c *DeptController) List(ctx *gin.Context) {
	var query dto.DeptQuery
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

	depts, total, err := c.deptService.GetDeptList(&query)
	if err != nil {
		response.FailWithMessage(201, "获取部门列表失败: "+err.Error(), ctx)
		return
	}

	response.OkWithDetailed(gin.H{
		"total": total,
		"list":  depts,
	}, "获取成功", ctx)
}

// GetTree 获取部门树
// @Summary 获取部门树
// @Description 获取部门树形结构
// @Tags 系统部门
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Success 200 {object} response.Response "成功"
// @Router /system/dept/tree [get]
func (c *DeptController) GetTree(ctx *gin.Context) {
	var query dto.DeptQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		response.FailWithMessage(201, "参数错误: "+err.Error(), ctx)
		return
	}

	// 树形查询，不分页
	query.PageNum = 0
	query.PageSize = 0

	depts, _, err := c.deptService.GetDeptList(&query)
	if err != nil {
		response.FailWithMessage(201, "获取部门树失败: "+err.Error(), ctx)
		return
	}

	response.OkWithDetailed(depts, "获取成功", ctx)
}

// Add 新增部门
// @Summary 新增部门
// @Description 新增系统部门
// @Tags 系统部门
// @Accept json
// @Produce json
// @Param data body dto.CreateDeptReq true "部门信息"
// @Success 200 {object} response.Response "成功"
// @Router /system/dept [post]
func (c *DeptController) Add(ctx *gin.Context) {
	var req dto.CreateDeptReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(201, "参数错误: "+err.Error(), ctx)
		return
	}

	// 从上下文获取创建人
	createBy := ctx.GetString("userName")
	if createBy == "" {
		createBy = "admin"
	}

	if err := c.deptService.CreateDept(&req, createBy); err != nil {
		response.FailWithMessage(201, "创建部门失败: "+err.Error(), ctx)
		return
	}

	response.OkWithMessage("创建成功", ctx)
}

// Update 更新部门
// @Summary 更新部门
// @Description 更新系统部门
// @Tags 系统部门
// @Accept json
// @Produce json
// @Param data body dto.UpdateDeptReq true "部门信息"
// @Success 200 {object} response.Response "成功"
// @Router /system/dept [put]
func (c *DeptController) Update(ctx *gin.Context) {
	var req dto.UpdateDeptReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(201, "参数错误: "+err.Error(), ctx)
		return
	}

	// 从上下文获取更新人
	updateBy := ctx.GetString("userName")
	if updateBy == "" {
		updateBy = "admin"
	}

	if err := c.deptService.UpdateDept(&req, updateBy); err != nil {
		response.FailWithMessage(201, "更新部门失败: "+err.Error(), ctx)
		return
	}

	response.OkWithMessage("更新成功", ctx)
}

// Delete 删除部门
// @Summary 删除部门
// @Description 删除系统部门
// @Tags 系统部门
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Param deptId path int true "部门ID"
// @Success 200 {object} response.Response "成功"
// @Router /system/dept/{deptId} [delete]
func (c *DeptController) Delete(ctx *gin.Context) {
	idStr := ctx.Param("deptId")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.FailWithMessage(201, "无效的部门ID", ctx)
		return
	}

	if err := c.deptService.DeleteDept(id); err != nil {
		response.FailWithMessage(201, "删除部门失败: "+err.Error(), ctx)
		return
	}

	response.OkWithMessage("删除成功", ctx)
}

// GetInfo 获取部门详情
// @Summary 获取部门详情
// @Description 获取系统部门详情
// @Tags 系统部门
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Param deptId path int true "部门ID"
// @Success 200 {object} response.Response "成功"
// @Router /system/dept/{deptId} [get]
func (c *DeptController) GetInfo(ctx *gin.Context) {
	idStr := ctx.Param("deptId")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.FailWithMessage(201, "无效的部门ID", ctx)
		return
	}

	dept, err := c.deptService.GetDeptInfo(id)
	if err != nil {
		response.FailWithMessage(201, "获取部门详情失败: "+err.Error(), ctx)
		return
	}

	response.OkWithDetailed(dept, "获取成功", ctx)
}

// UpdateStatus 更新部门状态
// @Summary 更新部门状态
// @Description 更新系统部门状态
// @Tags 系统部门
// @Accept json
// @Produce json
// @Param data body dto.UpdateDeptStatusReq true "状态信息"
// @Success 200 {object} response.Response "成功"
// @Router /system/dept/status [put]
func (c *DeptController) UpdateStatus(ctx *gin.Context) {
	var req dto.UpdateDeptStatusReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(201, "参数错误: "+err.Error(), ctx)
		return
	}

	if err := c.deptService.UpdateDeptStatus(req.DeptId, req.Status); err != nil {
		response.FailWithMessage(201, "更新部门状态失败: "+err.Error(), ctx)
		return
	}

	response.OkWithMessage("更新成功", ctx)
}

// GetDeptOptions 获取部门选项列表
// @Summary 获取部门选项列表
// @Description 获取部门选项列表（用于下拉选择）
// @Tags 系统部门
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Success 200 {object} response.Response "成功"
// @Router /system/dept/options [get]
func (c *DeptController) GetDeptOptions(ctx *gin.Context) {
	options, err := c.deptService.GetDeptOptions()
	if err != nil {
		response.FailWithMessage(201, "获取部门选项失败: "+err.Error(), ctx)
		return
	}

	response.OkWithDetailed(options, "获取成功", ctx)
}
