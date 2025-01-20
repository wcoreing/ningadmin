package system

import (
	"strconv"
	"wnadmin/app/common/response"
	"wnadmin/app/model/do/SystemDB"
	"wnadmin/app/model/dto"
	"wnadmin/app/service/system"

	"github.com/gin-gonic/gin"
)

type RoleController struct{}

// List 获取角色列表
func (c *RoleController) List(ctx *gin.Context) {
	var query dto.RoleQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		response.FailWithMessage(201, err.Error(), ctx)
		return
	}

	// 参数校验
	if query.PageNum < 1 {
		query.PageNum = 1
	}
	if query.PageSize < 1 {
		query.PageSize = 10
	}

	// 调用service层获取数据
	roleService := system.RoleService{}
	result, err := roleService.GetRoleList(&query)
	if err != nil {
		response.FailWithMessage(201, err.Error(), ctx)
		return
	}

	response.OkWithDetailed(gin.H{
		"total": result.Total,
		"list":  result.Rows,
	}, "获取成功", ctx)
}

// GetInfo 获取角色详情
func (c *RoleController) GetInfo(ctx *gin.Context) {
	idStr := ctx.Param("roleId")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.FailWithMessage(201, "无效的角色ID", ctx)
		return
	}

	roleService := system.RoleService{}
	role, err := roleService.GetRoleInfo(id)
	if err != nil {
		response.FailWithMessage(201, err.Error(), ctx)
		return
	}

	response.OkWithDetailed(role, "获取成功", ctx)
}

// Add 新增角色
func (c *RoleController) Add(ctx *gin.Context) {
	var req dto.CreateRoleReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(201, err.Error(), ctx)
		return
	}
	// 调用service层创建角色
	roleService := system.RoleService{}
	if err := roleService.CreateRole(&req, "admin"); err != nil {
		response.FailWithMessage(201, err.Error(), ctx)
		return
	}

	response.OkWithMessage("创建成功", ctx)
}

// Update 更新角色
func (c *RoleController) Update(ctx *gin.Context) {
	var req dto.UpdateRoleReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(201, err.Error(), ctx)
		return
	}

	// 调用service层更新数据
	roleService := system.RoleService{}
	if err := roleService.UpdateRole(&req, "admin"); err != nil {
		response.FailWithMessage(201, err.Error(), ctx)
		return
	}

	response.OkWithMessage("更新成功", ctx)
}

// Delete 删除角色
func (c *RoleController) Delete(ctx *gin.Context) {
	idStr := ctx.Param("roleId")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.FailWithMessage(201, "无效的角色ID", ctx)
		return
	}

	roleService := system.RoleService{}
	if err := roleService.DeleteRole(id); err != nil {
		response.FailWithMessage(201, err.Error(), ctx)
		return
	}

	response.OkWithMessage("删除成功", ctx)
}

// UpdateStatus 更新角色状态
func (c *RoleController) UpdateStatus(ctx *gin.Context) {
	var req dto.UpdateRoleStatusReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(201, err.Error(), ctx)
		return
	}

	roleService := system.RoleService{}
	if err := roleService.UpdateRoleStatus(&req, "admin"); err != nil {
		response.FailWithMessage(201, err.Error(), ctx)
		return
	}

	response.OkWithMessage("状态更新成功", ctx)
}

// GetRolePermissions 获取角色权限
func (c *RoleController) GetRolePermissions(ctx *gin.Context) {
	idStr := ctx.Param("roleId")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.FailWithMessage(201, "无效的角色ID", ctx)
		return
	}

	roleService := system.RoleService{}
	permissions, err := roleService.GetRolePermissionsTree(id)
	if err != nil {
		response.FailWithMessage(201, err.Error(), ctx)
		return
	}

	response.OkWithDetailed(permissions, "获取成功", ctx)
}

// UpdateRolePermissions 更新角色权限
func (c *RoleController) UpdateRolePermissions(ctx *gin.Context) {
	var req struct {
		RoleID        int64   `json:"roleId"`
		PermissionIDs []int64 `json:"permissionIds"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(201, err.Error(), ctx)
		return
	}

	role := &SystemDB.RoleDO{RoleID: req.RoleID}
	if err := role.UpdateRolePermissions(req.PermissionIDs); err != nil {
		response.FailWithMessage(201, err.Error(), ctx)
		return
	}

	response.OkWithMessage("权限更新成功", ctx)
}
