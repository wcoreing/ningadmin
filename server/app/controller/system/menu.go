package system

import (
	"strconv"
	"wnadmin/app/common/response"
	"wnadmin/app/model/do/SystemDB"
	"wnadmin/app/model/dto"
	"wnadmin/app/service/system"

	"github.com/gin-gonic/gin"
)

type MenuController struct {
	menuService *system.MenuService
}

// NewMenuController 创建菜单控制器
func NewMenuController() *MenuController {
	return &MenuController{
		menuService: &system.MenuService{},
	}
}

// List 获取菜单列表
func (c *MenuController) List(ctx *gin.Context) {
	var query dto.MenuListReq
	if err := ctx.ShouldBindQuery(&query); err != nil {
		response.FailWithMessage(201, err.Error(), ctx)
		return
	}

	menus, err := c.menuService.GetMenuList(&query)
	if err != nil {
		response.FailWithMessage(201, err.Error(), ctx)
		return
	}

	// 构建树形结构
	tree := c.menuService.BuildTree(menus, 0)
	response.OkWithDetailed(tree, "获取成功", ctx)
}

// Add 新增菜单
func (c *MenuController) Add(ctx *gin.Context) {
	var menuDTO dto.MenuDTO
	if err := ctx.ShouldBindJSON(&menuDTO); err != nil {
		response.FailWithMessage(201, err.Error(), ctx)
		return
	}

	menuDTO.Operator = ctx.GetString("username")
	if err := c.menuService.CreateMenu(&menuDTO); err != nil {
		response.FailWithMessage(201, err.Error(), ctx)
		return
	}

	response.OkWithMessage("创建成功", ctx)
}

// Update 更新菜单
func (c *MenuController) Update(ctx *gin.Context) {
	var menuDTO dto.MenuDTO
	if err := ctx.ShouldBindJSON(&menuDTO); err != nil {
		response.FailWithMessage(201, err.Error(), ctx)
		return
	}

	menuDTO.Operator = ctx.GetString("username")
	if err := c.menuService.UpdateMenu(&menuDTO); err != nil {
		response.FailWithMessage(201, err.Error(), ctx)
		return
	}

	response.OkWithMessage("更新成功", ctx)
}

// Delete 删除菜单
func (c *MenuController) Delete(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.FailWithMessage(201, "无效的ID", ctx)
		return
	}

	if err := c.menuService.DeleteMenu(id); err != nil {
		response.FailWithMessage(201, err.Error(), ctx)
		return
	}

	response.OkWithMessage("删除成功", ctx)
}

// GetInfo 获取菜单详情
func (c *MenuController) GetInfo(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.FailWithMessage(201, "无效的ID", ctx)
		return
	}

	menu, err := c.menuService.GetMenuInfo(id)
	if err != nil {
		response.FailWithMessage(201, err.Error(), ctx)
		return
	}

	response.OkWithDetailed(menu, "获取成功", ctx)
}

// GetTree 获取菜单树形数据
func (c *MenuController) GetTree(ctx *gin.Context) {
	var query dto.MenuListReq
	if err := ctx.ShouldBindQuery(&query); err != nil {
		response.FailWithMessage(201, err.Error(), ctx)
		return
	}

	menus, err := c.menuService.GetMenuTree(&query)
	if err != nil {
		response.FailWithMessage(201, err.Error(), ctx)
		return
	}

	tree := c.menuService.BuildTree(menus, 0)
	response.OkWithDetailed(tree, "获取成功", ctx)
}

// CheckPermission 检查权限标识是否存在
func (c *MenuController) CheckPermission(ctx *gin.Context) {
	permission := ctx.Query("permission")
	if permission == "" {
		response.FailWithMessage(201, "权限标识不能为空", ctx)
		return
	}

	excludeIDStr := ctx.Query("excludeId")
	var excludeID int64
	if excludeIDStr != "" {
		var err error
		excludeID, err = strconv.ParseInt(excludeIDStr, 10, 64)
		if err != nil {
			response.FailWithMessage(201, "无效的排除ID", ctx)
			return
		}
	}

	menu := &SystemDB.MenuDO{}
	exists, err := menu.CheckPermissionExists(permission, excludeID)
	if err != nil {
		response.FailWithMessage(201, err.Error(), ctx)
		return
	}

	response.OkWithDetailed(exists, "获取成功", ctx)
}

// UpdateSort 更新菜单排序
func (c *MenuController) UpdateSort(ctx *gin.Context) {
	var req struct {
		ID   int64 `json:"id" validate:"required"`
		Sort int64 `json:"sort" validate:"gte=0"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(201, err.Error(), ctx)
		return
	}

	if err := c.menuService.UpdateSort(req.ID, req.Sort); err != nil {
		response.FailWithMessage(201, "更新排序失败", ctx)
		return
	}

	response.Ok(ctx)
}

// UpdateStatus 更新菜单状态
func (c *MenuController) UpdateStatus(ctx *gin.Context) {
	var req struct {
		ID     int64  `json:"id" validate:"required"`
		Status string `json:"status" validate:"required,oneof=0 1"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(201, err.Error(), ctx)
		return
	}

	if err := c.menuService.UpdateStatus(req.ID, req.Status); err != nil {
		response.FailWithMessage(201, "更新状态失败", ctx)
		return
	}

	response.OkWithMessage("更新成功", ctx)
}
