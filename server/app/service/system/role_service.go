package system

import (
	"time"
	"wnadmin/app/model/do/SystemDB"
	"wnadmin/app/model/dto"
)

// RoleService 角色服务
type RoleService struct{}

// GetRoleList 获取角色列表
func (s *RoleService) GetRoleList(query *dto.RoleQuery) (*dto.RoleListResp, error) {
	var roles []SystemDB.RoleDO
	roleDO := &SystemDB.RoleDO{}
	db := roleDO.Query()

	// 构建查询条件
	if query.RoleName != "" {
		db = db.Where("role_name LIKE CONCAT('%', ?, '%')", query.RoleName)
	}
	if query.RoleKey != "" {
		db = db.Where("role_key = ?", query.RoleKey)
	}
	if query.Status != "" {
		db = db.Where("status = ?", query.Status)
	}

	// 查询总数
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}

	// 分页查询
	offset := (query.PageNum - 1) * query.PageSize
	if err := db.Offset(offset).Limit(query.PageSize).Order("create_time DESC").Find(&roles).Error; err != nil {
		return nil, err
	}

	// 构建响应数据
	resp := &dto.RoleListResp{
		Total: total,
		Rows:  make([]dto.RoleInfoResp, 0, len(roles)),
	}

	for _, role := range roles {
		resp.Rows = append(resp.Rows, dto.RoleInfoResp{
			RoleID:     role.RoleID,
			RoleName:   role.RoleName,
			RoleKey:    role.RoleKey,
			Status:     role.Status,
			CreateTime: role.CreateTime.Format("2006-01-02 15:04:05"),
			CreateBy:   role.CreateBy,
			UpdateTime: role.UpdateTime.Format("2006-01-02 15:04:05"),
			UpdateBy:   role.UpdateBy,
		})
	}

	return resp, nil
}

// GetRoleInfo 获取角色详情
func (s *RoleService) GetRoleInfo(id int64) (*dto.RoleInfoResp, error) {
	role := &SystemDB.RoleDO{RoleID: id}
	if err := role.First(); err != nil {
		return nil, ErrRoleNotFound
	}

	return &dto.RoleInfoResp{
		RoleID:     role.RoleID,
		RoleName:   role.RoleName,
		RoleKey:    role.RoleKey,
		Status:     role.Status,
		CreateTime: role.CreateTime.Format("2006-01-02 15:04:05"),
		CreateBy:   role.CreateBy,
		UpdateTime: role.UpdateTime.Format("2006-01-02 15:04:05"),
		UpdateBy:   role.UpdateBy,
	}, nil
}

// CreateRole 创建角色
func (s *RoleService) CreateRole(req *dto.CreateRoleReq, createBy string) error {
	// 检查角色标识是否已存在
	roleDO := &SystemDB.RoleDO{}
	exists, err := roleDO.CheckRoleKeyExists(req.RoleKey, 0)
	if err != nil {
		return err
	}
	if exists {
		return ErrRoleKeyExists
	}
	// 创建角色
	now := time.Now()
	role := &SystemDB.RoleDO{}
	role.RoleName = req.RoleName
	role.RoleKey = req.RoleKey
	role.Status = req.Status
	role.CreateBy = createBy
	role.CreateTime = now
	role.UpdateTime = now
	role.DelFlag = "0"

	// 设置默认状态
	if role.Status == "" {
		role.Status = "0"
	}

	return role.Create()
}

// UpdateRole 更新角色
func (s *RoleService) UpdateRole(req *dto.UpdateRoleReq, updateBy string) error {
	// 检查角色是否存在
	role := &SystemDB.RoleDO{RoleID: req.RoleID}
	if err := role.First(); err != nil {
		return ErrRoleNotFound
	}

	// 检查角色标识是否已存在（排除自身）
	if req.RoleKey != role.RoleKey {
		exists, err := role.CheckRoleKeyExists(req.RoleKey, req.RoleID)
		if err != nil {
			return err
		}
		if exists {
			return ErrRoleKeyExists
		}
	}

	// 更新角色信息
	role.RoleName = req.RoleName
	role.RoleKey = req.RoleKey
	role.Status = req.Status
	role.UpdateBy = updateBy
	role.UpdateTime = time.Now()

	return role.Update()
}

// DeleteRole 删除角色
func (s *RoleService) DeleteRole(id int64) error {
	// 检查角色是否存在
	role := &SystemDB.RoleDO{RoleID: id}
	if err := role.First(); err != nil {
		return ErrRoleNotFound
	}

	// 检查角色是否已分配用户
	if hasUsers, err := role.HasUsers(); err != nil {
		return err
	} else if hasUsers {
		return ErrRoleHasUsers
	}

	return role.Delete()
}

// UpdateRoleStatus 更新角色状态
func (s *RoleService) UpdateRoleStatus(req *dto.UpdateRoleStatusReq, updateBy string) error {
	// 检查角色是否存在
	role := &SystemDB.RoleDO{RoleID: req.RoleID}
	if err := role.First(); err != nil {
		return ErrRoleNotFound
	}

	// 更新状态
	role.Status = req.Status
	role.UpdateBy = updateBy
	role.UpdateTime = time.Now()

	return role.UpdateStatus()
}

// GetRolePermissions 获取角色权限ID列表
func (s *RoleService) GetRolePermissions(roleId int64) ([]int64, error) {
	// 检查角色是否存在
	role := &SystemDB.RoleDO{RoleID: roleId}
	if err := role.First(); err != nil {
		return nil, ErrRoleNotFound
	}

	return role.GetRolePermissions()
}

// UpdateRolePermissions 更新角色权限
func (s *RoleService) UpdateRolePermissions(req *dto.RolePermissionReq, updateBy string) error {
	// 检查角色是否存在
	role := &SystemDB.RoleDO{RoleID: req.RoleID}
	if err := role.First(); err != nil {
		return ErrRoleNotFound
	}

	// 更新角色权限
	if err := role.UpdateRolePermissions(req.PermissionIds); err != nil {
		return err
	}

	// 更新角色信息
	role.UpdateBy = updateBy
	role.UpdateTime = time.Now()
	return role.Update()
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

// GetRolePermissionsTree 获取角色权限树
func (s *RoleService) GetRolePermissionsTree(roleId int64) ([]map[string]interface{}, error) {
	// 1. 获取所有菜单
	menuDO := &SystemDB.MenuDO{}
	var allMenus []SystemDB.MenuDO
	err := menuDO.Query().Where("del_flag = ?", "0").Order("sort").Find(&allMenus).Error
	if err != nil {
		return nil, err
	}

	// 2. 获取角色已分配的权限ID列表
	role := &SystemDB.RoleDO{RoleID: roleId}
	permissionIds, err := role.GetRolePermissions()
	if err != nil {
		return nil, err
	}

	// 3. 构建树形结构
	return buildPermissionTree(allMenus, 0, permissionIds), nil
}
