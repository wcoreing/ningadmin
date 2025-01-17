package system

import (
	"strconv"
	"wnadmin/app/common/response"
	"wnadmin/app/model/do/SystemDB"

	"github.com/gin-gonic/gin"
)

type RoleController struct{}

// List 获取角色列表
func (c *RoleController) List(ctx *gin.Context) {
	var query struct {
		PageNum  int    `form:"pageNum"`
		PageSize int    `form:"pageSize"`
		RoleName string `form:"roleName"`
		Status   string `form:"status"`
	}

	if err := ctx.ShouldBindQuery(&query); err != nil {
		response.FailWithMessage(201, err.Error(), ctx)
		return
	}

	if query.PageNum < 1 {
		query.PageNum = 1
	}
	if query.PageSize < 1 {
		query.PageSize = 10
	}

	roleDO := &SystemDB.RoleDO{}
	db := roleDO.Query().Where("del_flag = ?", "0")

	if query.RoleName != "" {
		db = db.Where("role_name LIKE ?", "%"+query.RoleName+"%")
	}
	if query.Status != "" {
		db = db.Where("status = ?", query.Status)
	}

	var total int64
	err := db.Count(&total).Error
	if err != nil {
		response.FailWithMessage(201, err.Error(), ctx)
		return
	}

	var roles []SystemDB.RoleDO
	err = db.Offset((query.PageNum - 1) * query.PageSize).
		Limit(query.PageSize).
		Order("role_id desc").
		Find(&roles).Error
	if err != nil {
		response.FailWithMessage(201, err.Error(), ctx)
		return
	}

	response.OkWithDetailed(gin.H{
		"total": total,
		"list":  roles,
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

	role := &SystemDB.RoleDO{ID: id}
	if err := role.First(); err != nil {
		response.FailWithMessage(201, err.Error(), ctx)
		return
	}

	response.OkWithDetailed(role, "获取成功", ctx)
}

// Add 新增角色
func (c *RoleController) Add(ctx *gin.Context) {
	var role SystemDB.RoleDO
	if err := ctx.ShouldBindJSON(&role); err != nil {
		response.FailWithMessage(201, err.Error(), ctx)
		return
	}

	// 检查角色标识是否存在
	exists, err := role.CheckRoleKeyExists(role.RoleKey, 0)
	if err != nil {
		response.FailWithMessage(201, err.Error(), ctx)
		return
	}
	if exists {
		response.FailWithMessage(201, "角色标识已存在", ctx)
		return
	}

	role.DelFlag = "0"
	if err := role.Create(); err != nil {
		response.FailWithMessage(201, err.Error(), ctx)
		return
	}

	response.OkWithMessage("创建成功", ctx)
}

// Update 更新角色
func (c *RoleController) Update(ctx *gin.Context) {
	var role SystemDB.RoleDO
	if err := ctx.ShouldBindJSON(&role); err != nil {
		response.FailWithMessage(201, err.Error(), ctx)
		return
	}

	// 检查角色标识是否存在
	exists, err := role.CheckRoleKeyExists(role.RoleKey, role.ID)
	if err != nil {
		response.FailWithMessage(201, err.Error(), ctx)
		return
	}
	if exists {
		response.FailWithMessage(201, "角色标识已存在", ctx)
		return
	}

	if err := role.Update(); err != nil {
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

	role := &SystemDB.RoleDO{ID: id}
	if err := role.Delete(); err != nil {
		response.FailWithMessage(201, err.Error(), ctx)
		return
	}

	response.OkWithMessage("删除成功", ctx)
}

// UpdateStatus 更新角色状态
func (c *RoleController) UpdateStatus(ctx *gin.Context) {
	var req struct {
		RoleID int64  `json:"roleId"`
		Status string `json:"status"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(201, err.Error(), ctx)
		return
	}

	role := &SystemDB.RoleDO{
		ID:     req.RoleID,
		Status: req.Status,
	}

	if err := role.UpdateStatus(); err != nil {
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

	// 1. 获取所有菜单权限
	menuDO := &SystemDB.MenuDO{}
	var allMenus []SystemDB.MenuDO
	err = menuDO.Query().Where("del_flag = ?", "0").Order("sort").Find(&allMenus).Error
	if err != nil {
		response.FailWithMessage(201, err.Error(), ctx)
		return
	}

	// 2. 获取角色已分配的权限ID列表
	role := &SystemDB.RoleDO{ID: id}
	permissionIDs, err := role.GetRolePermissions()
	if err != nil {
		response.FailWithMessage(201, err.Error(), ctx)
		return
	}

	// 3. 构建树形结构并标记选中状态
	var result []map[string]interface{}
	for _, menu := range allMenus {
		if menu.ParentID == 0 {
			node := convertMenuToMap(menu)
			// 检查是否在已分配的权限列表中
			for _, pid := range permissionIDs {
				if menu.ID == pid {
					node["isChecked"] = true
					break
				}
			}
			children := buildPermissionTree(allMenus, menu.ID, permissionIDs)
			if len(children) > 0 {
				node["children"] = children
			}
			result = append(result, node)
		}
	}

	response.OkWithDetailed(result, "获取成功", ctx)
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

	role := &SystemDB.RoleDO{ID: req.RoleID}
	if err := role.UpdateRolePermissions(req.PermissionIDs); err != nil {
		response.FailWithMessage(201, err.Error(), ctx)
		return
	}

	response.OkWithMessage("权限更新成功", ctx)
}

// convertMenuToMap 将菜单对象转换为map
func convertMenuToMap(menu SystemDB.MenuDO) map[string]interface{} {
	return map[string]interface{}{
		"id":            menu.ID,
		"parentId":      menu.ParentID,
		"type":          menu.Type,
		"name":          menu.Name,
		"icon":          menu.Icon,
		"routePath":     menu.RoutePath,
		"componentPath": menu.ComponentPath,
		"permission":    menu.Permission,
		"sort":          menu.Sort,
		"status":        menu.Status,
		"createTime":    menu.CreateTime,
		"createBy":      menu.CreateBy,
		"updateTime":    menu.UpdateTime,
		"updateBy":      menu.UpdateBy,
		"isChecked":     false,
		"isExpand":      false,
	}
}

// buildPermissionTree 构建权限树
func buildPermissionTree(menus []SystemDB.MenuDO, parentID int64, checkedIds []int64) []map[string]interface{} {
	var nodes []map[string]interface{}
	for _, menu := range menus {
		if menu.ParentID == parentID {
			node := convertMenuToMap(menu)
			// 检查是否在已分配的权限列表中
			for _, pid := range checkedIds {
				if menu.ID == pid {
					node["isChecked"] = true
					break
				}
			}
			children := buildPermissionTree(menus, menu.ID, checkedIds)
			if len(children) > 0 {
				node["children"] = children
			}
			nodes = append(nodes, node)
		}
	}
	return nodes
}
