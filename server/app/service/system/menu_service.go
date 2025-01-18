package system

import (
	"strings"
	"wnadmin/app/model/do/SystemDB"
	"wnadmin/app/model/dto"
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
	err := db.Order("sort").Find(&menus).Error
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
	menu := &SystemDB.MenuDO{
		ID:            req.ID,
		ParentID:      req.ParentID,
		Type:          req.Type,
		Name:          req.Name,
		Icon:          req.Icon,
		RoutePath:     req.RoutePath,
		ComponentPath: req.ComponentPath,
		Permission:    req.Permission,
		Sort:          req.Sort,
		Status:        req.Status,
		UpdateBy:      req.Operator,
	}

	// 检查权限标识是否存在
	if req.Permission != "" {
		exists, err := menu.CheckPermissionExists(req.Permission, req.ID)
		if err != nil {
			return err
		}
		if exists {
			return ErrPermissionExists
		}
	}

	return menu.Update()
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
	err := menuDO.Query().Order("sort").Find(&allMenus).Error
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
