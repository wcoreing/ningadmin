package system

import (
	"strconv"
	"strings"
	"wnadmin/app/common/response"
	"wnadmin/app/model/do/SystemDB"

	"github.com/gin-gonic/gin"
)

type MenuController struct{}

// List 获取菜单列表
func (c *MenuController) List(ctx *gin.Context) {
	var query struct {
		Keyword string `form:"keyword"`
		Status  string `form:"status"`
	}

	if err := ctx.ShouldBindQuery(&query); err != nil {
		response.FailWithMessage(201, err.Error(), ctx)
		return
	}

	menuDO := &SystemDB.MenuDO{}
	db := menuDO.Query().Where("del_flag = ?", "0")

	if query.Keyword != "" {
		db = db.Where("permission_name LIKE ? OR permission_code LIKE ?", "%"+query.Keyword+"%", "%"+query.Keyword+"%")
	}
	if query.Status != "" {
		db = db.Where("status = ?", query.Status)
	}

	var menus []SystemDB.MenuDO
	err := db.Order("sort").Find(&menus).Error
	if err != nil {
		response.FailWithMessage(201, err.Error(), ctx)
		return
	}

	// 构建树形结构
	tree := buildTree(menus, 0)
	response.OkWithDetailed(tree, "获取成功", ctx)
}

// Add 新增菜单
func (c *MenuController) Add(ctx *gin.Context) {
	var menu SystemDB.MenuDO
	if err := ctx.ShouldBindJSON(&menu); err != nil {
		response.FailWithMessage(201, err.Error(), ctx)
		return
	}

	// 检查权限标识是否存在
	if menu.Permission != "" {
		exists, err := menu.CheckPermissionExists(menu.Permission, 0)
		if err != nil {
			response.FailWithMessage(201, err.Error(), ctx)
			return
		}
		if exists {
			response.FailWithMessage(201, "权限标识已存在", ctx)
			return
		}
	}

	menu.DelFlag = "0"
	if err := menu.Create(); err != nil {
		response.FailWithMessage(201, err.Error(), ctx)
		return
	}

	response.OkWithMessage("创建成功", ctx)
}

// Update 更新菜单
func (c *MenuController) Update(ctx *gin.Context) {
	var menu SystemDB.MenuDO
	if err := ctx.ShouldBindJSON(&menu); err != nil {
		response.FailWithMessage(201, err.Error(), ctx)
		return
	}

	// 检查权限标识是否存在
	if menu.Permission != "" {
		exists, err := menu.CheckPermissionExists(menu.Permission, menu.ID)
		if err != nil {
			response.FailWithMessage(201, err.Error(), ctx)
			return
		}
		if exists {
			response.FailWithMessage(201, "权限标识已存在", ctx)
			return
		}
	}

	if err := menu.Update(); err != nil {
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

	// 检查是否有子菜单
	menuDO := &SystemDB.MenuDO{}
	var count int64
	err = menuDO.Query().Where("parent_id = ? AND del_flag = ?", id, "0").Count(&count).Error
	if err != nil {
		response.FailWithMessage(201, err.Error(), ctx)
		return
	}
	if count > 0 {
		response.FailWithMessage(201, "存在子菜单，无法删除", ctx)
		return
	}

	menu := &SystemDB.MenuDO{ID: id}
	if err := menu.Delete(); err != nil {
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

	menu := &SystemDB.MenuDO{ID: id}
	if err := menu.First(); err != nil {
		response.FailWithMessage(201, err.Error(), ctx)
		return
	}

	response.OkWithDetailed(menu, "获取成功", ctx)
}

// GetTree 获取菜单树形数据
func (c *MenuController) GetTree(ctx *gin.Context) {
	var query struct {
		Keyword string `form:"keyword"`
		Status  string `form:"status"`
	}

	if err := ctx.ShouldBindQuery(&query); err != nil {
		response.FailWithMessage(201, err.Error(), ctx)
		return
	}

	menuDO := &SystemDB.MenuDO{}

	// 1. 先获取所有菜单
	var allMenus []SystemDB.MenuDO
	err := menuDO.Query().Where("del_flag = ?", "0").Order("sort").Find(&allMenus).Error
	if err != nil {
		response.FailWithMessage(201, err.Error(), ctx)
		return
	}

	// 2. 根据条件筛选匹配的菜单
	var matchedMenuIDs []int64
	for _, menu := range allMenus {
		isMatched := true

		// 关键字匹配
		if query.Keyword != "" {
			isMatched = false
			keyword := "%" + query.Keyword + "%"
			if containsKeyword(menu.Name, keyword) || containsKeyword(menu.Permission, keyword) {
				isMatched = true
			}
		}

		// 状态匹配
		if query.Status != "" && menu.Status != query.Status {
			isMatched = false
		}

		if isMatched {
			matchedMenuIDs = append(matchedMenuIDs, menu.ID)
			// 添加所有父级ID
			parentID := menu.ParentID
			for parentID != 0 {
				if !contains(matchedMenuIDs, parentID) {
					matchedMenuIDs = append(matchedMenuIDs, parentID)
				}
				// 查找父级的父级
				for _, m := range allMenus {
					if m.ID == parentID {
						parentID = m.ParentID
						break
					}
				}
			}
		}
	}

	// 3. 筛选出匹配的菜单和它们的父级菜单
	var filteredMenus []SystemDB.MenuDO
	for _, menu := range allMenus {
		if contains(matchedMenuIDs, menu.ID) {
			filteredMenus = append(filteredMenus, menu)
		}
	}

	// 4. 构建树形结构
	tree := buildTree(filteredMenus, 0)
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

	menuDO := &SystemDB.MenuDO{}
	exists, err := menuDO.CheckPermissionExists(permission, excludeID)
	if err != nil {
		response.FailWithMessage(201, err.Error(), ctx)
		return
	}

	response.OkWithDetailed(exists, "获取成功", ctx)
}

// 构建树形结构
func buildTree(menus []SystemDB.MenuDO, parentID int64) []SystemDB.MenuDO {
	var trees []SystemDB.MenuDO
	for _, menu := range menus {
		if menu.ParentID == parentID {
			menu.Children = buildTree(menus, menu.ID)
			trees = append(trees, menu)
		}
	}
	return trees
}

// 检查字符串是否包含关键字（不区分大小写）
func containsKeyword(str string, keyword string) bool {
	// 去除首尾的 % 用于实际比较
	keyword = keyword[1 : len(keyword)-1]
	return strings.Contains(strings.ToLower(str), strings.ToLower(keyword))
}

// 检查slice中是否包含某个值
func contains(slice []int64, val int64) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}
