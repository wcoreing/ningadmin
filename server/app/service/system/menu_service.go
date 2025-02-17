package system

import (
	"strings"
	"wnadmin/app/model/do/SystemDB"
	"wnadmin/app/model/dto"
	"wnadmin/app/provider/db"
)

type MenuService struct{}

// GetMenuList 获取菜单列表
func (s *MenuService) GetMenuList(req *dto.MenuListReq) ([]SystemDB.MenuDO, error) {
	menuDO := &SystemDB.MenuDO{}
	db := menuDO.Query()

	if req.Keyword != "" {
		db = db.Where("permission_name LIKE ? OR permission_code LIKE ?", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}
	if req.Status != "" {
		db = db.Where("status = ?", req.Status)
	}

	var menus []SystemDB.MenuDO
	err := db.Order("sort desc").Find(&menus).Error
	return menus, err
}

// CreateMenu 创建菜单
func (s *MenuService) CreateMenu(req *dto.MenuDTO) error {
	menu := &SystemDB.MenuDO{
		ParentID:      req.ParentID,
		Type:          req.Type,
		Name:          req.Name,
		Icon:          req.Icon,
		RoutePath:     req.RoutePath,
		ComponentPath: req.ComponentPath,
		Permission:    req.Permission,
		Sort:          req.Sort,
		Status:        req.Status,
		CreateBy:      req.Operator,
		DelFlag:       "0",
	}

	// 检查权限标识是否存在
	if req.Permission != "" {
		exists, err := menu.CheckPermissionExists(req.Permission, 0)
		if err != nil {
			return err
		}
		if exists {
			return ErrPermissionExists
		}
	}

	return menu.Create()
}

// UpdateMenu 更新菜单
func (s *MenuService) UpdateMenu(req *dto.MenuDTO) error {
	// 如果只传入了 ID 和 Sort，则只更新排序
	if len(req.Type) == 0 && len(req.Name) == 0 && len(req.Permission) == 0 {
		return db.DB.WnAdmin.Model(&SystemDB.MenuDO{}).
			Where("id = ?", req.ID).
			Update("sort", req.Sort).Error
	}

	// 构建更新字段map
	updateFields := map[string]interface{}{
		"parent_id":       req.ParentID,
		"permission_type": req.Type,
		"permission_name": req.Name,
		"icon":            req.Icon,
		"route_path":      req.RoutePath,
		"component_path":  req.ComponentPath,
		"permission_code": req.Permission,
		"sort":            req.Sort,
		"status":          req.Status,
		"update_by":       req.Operator,
	}

	// 检查权限标识是否存在
	if req.Permission != "" {
		exists, err := (&SystemDB.MenuDO{}).CheckPermissionExists(req.Permission, req.ID)
		if err != nil {
			return err
		}
		if exists {
			return ErrPermissionExists
		}
	}

	return db.DB.WnAdmin.Model(&SystemDB.MenuDO{}).
		Where("permission_id = ?", req.ID).
		Updates(updateFields).Error
}

// DeleteMenu 删除菜单
func (s *MenuService) DeleteMenu(id int64) error {
	menu := &SystemDB.MenuDO{ID: id}

	// 检查是否有子菜单
	hasChildren, err := menu.HasChildren()
	if err != nil {
		return err
	}
	if hasChildren {
		return ErrHasChildren
	}

	// 检查是否已分配角色
	hasRoles, err := menu.HasRoles()
	if err != nil {
		return err
	}
	if hasRoles {
		return ErrMenuHasRoles
	}

	return menu.Delete()
}

// GetMenuInfo 获取菜单信息
func (s *MenuService) GetMenuInfo(id int64) (*SystemDB.MenuDO, error) {
	menu := &SystemDB.MenuDO{ID: id}
	err := menu.First()
	return menu, err
}

// GetMenuTree 获取菜单树
func (s *MenuService) GetMenuTree(req *dto.MenuListReq) ([]SystemDB.MenuDO, error) {
	menuDO := &SystemDB.MenuDO{}

	// 获取所有菜单
	var allMenus []SystemDB.MenuDO
	err := menuDO.Query().Order("sort desc").Find(&allMenus).Error
	if err != nil {
		return nil, err
	}

	// 根据条件筛选匹配的菜单
	var matchedMenuIDs []int64
	for _, menu := range allMenus {
		isMatched := true

		// 关键字匹配
		if req.Keyword != "" {
			isMatched = false
			if containsKeyword(menu.Name, req.Keyword) || containsKeyword(menu.Permission, req.Keyword) {
				isMatched = true
			}
		}

		// 状态匹配
		if req.Status != "" && menu.Status != req.Status {
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

	// 筛选出匹配的菜单
	var filteredMenus []SystemDB.MenuDO
	for _, menu := range allMenus {
		if contains(matchedMenuIDs, menu.ID) {
			filteredMenus = append(filteredMenus, menu)
		}
	}

	return filteredMenus, nil
}

// BuildTree 构建树形结构
func (s *MenuService) BuildTree(menus []SystemDB.MenuDO, parentID int64) []SystemDB.MenuDO {
	var trees []SystemDB.MenuDO
	for _, menu := range menus {
		if menu.ParentID == parentID {
			menu.Children = s.BuildTree(menus, menu.ID)
			trees = append(trees, menu)
		}
	}
	return trees
}

// 检查字符串是否包含关键字（不区分大小写）
func containsKeyword(str, keyword string) bool {
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

// GetRoleMenuTree 获取角色的菜单树
func (s *MenuService) GetRoleMenuTree(roleID int64, permissionType string) ([]dto.MenuTree, error) {
	// 获取角色的权限菜单
	var menus []SystemDB.MenuDO
	err := db.DB.WnAdmin.Table("sys_permission").
		Joins("LEFT JOIN sys_role_permission ON sys_permission.permission_id = sys_role_permission.permission_id").
		Where(" sys_permission.permission_type = ? AND sys_role_permission.role_id = ? AND sys_permission.del_flag = ? AND sys_permission.status = ?",
			permissionType, roleID, "0", "0").
		Order("sys_permission.sort desc").
		Find(&menus).Error
	if err != nil {
		return nil, err
	}

	// 构建菜单树
	return buildRoleMenuTree(menus, 0), nil
}

// buildRoleMenuTree 构建角色菜单树
func buildRoleMenuTree(menus []SystemDB.MenuDO, parentID int64) []dto.MenuTree {
	var trees []dto.MenuTree
	for _, menu := range menus {
		if menu.ParentID == parentID {
			tree := dto.MenuTree{
				ID:        menu.ID,
				ParentID:  menu.ParentID,
				Name:      menu.Name,
				Path:      menu.RoutePath,
				Component: menu.ComponentPath,
				Icon:      menu.Icon,
				Sort:      menu.Sort,
				Children:  buildRoleMenuTree(menus, menu.ID),
			}
			trees = append(trees, tree)
		}
	}
	return trees
}

// UpdateSort 更新菜单排序
func (s *MenuService) UpdateSort(id int64, sort int64) error {
	return db.DB.WnAdmin.Model(&SystemDB.MenuDO{}).
		Where("permission_id = ?", id).
		Update("sort", sort).Error
}
