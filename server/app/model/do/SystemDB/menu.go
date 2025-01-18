package SystemDB

import (
	"time"
	"wnadmin/app/provider/db"

	"gorm.io/gorm"
)

// MenuDO 菜单权限
type MenuDO struct {
	ID            int64     `json:"id" gorm:"column:permission_id;primaryKey;comment:权限ID"`
	ParentID      int64     `json:"parentId" gorm:"column:parent_id;comment:上级权限ID"`
	Type          string    `json:"type" gorm:"column:permission_type;comment:权限类型（菜单、按钮等）"`
	Name          string    `json:"name" gorm:"column:permission_name;comment:权限名称"`
	Icon          string    `json:"icon" gorm:"comment:图标"`
	RoutePath     string    `json:"routePath" gorm:"column:route_path;comment:路由路径"`
	ComponentPath string    `json:"componentPath" gorm:"column:component_path;comment:组件路径"`
	Permission    string    `json:"permission" gorm:"column:permission_code;comment:权限标识"`
	Sort          int       `json:"sort" gorm:"comment:排序"`
	Status        string    `json:"status" gorm:"comment:状态（0显示 1隐藏）"`
	DelFlag       string    `json:"-" gorm:"column:del_flag;comment:删除标志（0代表存在 2代表删除）"`
	CreateTime    time.Time `json:"createTime" gorm:"column:create_time;autoCreateTime;comment:创建时间"`
	CreateBy      string    `json:"createBy" gorm:"column:create_by;comment:创建者"`
	UpdateTime    time.Time `json:"updateTime" gorm:"column:update_time;autoUpdateTime;comment:更新时间"`
	UpdateBy      string    `json:"updateBy" gorm:"column:update_by;comment:更新者"`
	Children      []MenuDO  `json:"children" gorm:"-"`
	IsExpand      bool      `json:"isExpand" gorm:"-"`
}

// RoleMenuDO 角色菜单关联
type RoleMenuDO struct {
	RoleID       int64 `gorm:"column:role_id;primaryKey;comment:角色ID"`
	PermissionID int64 `gorm:"column:permission_id;primaryKey;comment:权限ID"`
}

// TableName 表名
func (RoleMenuDO) TableName() string {
	return "sys_role_permission"
}

// TableName 表名
func (MenuDO) TableName() string {
	return "sys_permission"
}

// Query 获取查询构建器
func (m *MenuDO) Query() *gorm.DB {
	return db.DB.WnAdmin.Model(m).Where("del_flag = ?", "0")
}

// Create 创建记录
func (m *MenuDO) Create() error {
	return m.Query().Create(m).Error
}

// Update 更新记录
func (m *MenuDO) Update() error {
	return m.Query().Where("permission_id = ?", m.ID).Updates(m).Error
}

// Delete 删除记录
func (m *MenuDO) Delete() error {
	return db.DB.WnAdmin.Transaction(func(tx *gorm.DB) error {
		// 删除角色菜单关联
		if err := tx.Table("sys_role_permission").Where("permission_id = ?", m.ID).Delete(nil).Error; err != nil {
			return err
		}

		// 软删除菜单
		return tx.Model(m).Where("permission_id = ?", m.ID).Update("del_flag", "2").Error
	})
}

// First 获取第一条记录
func (m *MenuDO) First() error {
	return m.Query().Where("permission_id = ?", m.ID).First(m).Error
}

// Find 查询多条记录
func (m *MenuDO) Find(dest interface{}) error {
	return m.Query().Find(dest).Error
}

// CheckPermissionExists 检查权限标识是否存在
func (m *MenuDO) CheckPermissionExists(permission string, excludeID int64) (bool, error) {
	var count int64
	db := m.Query().Where("permission_code = ? ", permission)
	if excludeID > 0 {
		db = db.Where("permission_id != ?", excludeID)
	}
	err := db.Count(&count).Error
	return count > 0, err
}

// UpdateStatus 更新菜单状态
func (m *MenuDO) UpdateStatus() error {
	return m.Query().Where("permission_id = ?", m.ID).Update("status", m.Status).Error
}

// HasRoles 检查菜单是否已分配角色
func (m *MenuDO) HasRoles() (bool, error) {
	var count int64
	err := db.DB.WnAdmin.Model(&RoleMenuDO{}).Where("permission_id = ?", m.ID).Count(&count).Error
	return count > 0, err
}

// GetMenuRoles 获取菜单关联的角色ID列表
func (m *MenuDO) GetMenuRoles() ([]int64, error) {
	var roleIDs []int64
	err := db.DB.WnAdmin.Table("sys_role_permission").
		Where("permission_id = ?", m.ID).
		Pluck("role_id", &roleIDs).Error
	return roleIDs, err
}

// UpdateMenuRoles 更新菜单关联的角色
func (m *MenuDO) UpdateMenuRoles(roleIDs []int64) error {
	return db.DB.WnAdmin.Transaction(func(tx *gorm.DB) error {
		// 删除原有角色关联
		if err := tx.Table("sys_role_permission").Where("permission_id = ?", m.ID).Delete(nil).Error; err != nil {
			return err
		}

		// 添加新角色关联
		if len(roleIDs) > 0 {
			menuRoles := make([]map[string]interface{}, len(roleIDs))
			for i, rid := range roleIDs {
				menuRoles[i] = map[string]interface{}{
					"permission_id": m.ID,
					"role_id":       rid,
				}
			}
			return tx.Table("sys_role_permission").Create(menuRoles).Error
		}
		return nil
	})
}

// HasChildren 检查是否有子菜单
func (m *MenuDO) HasChildren() (bool, error) {
	var count int64
	err := m.Query().Where("parent_id = ?", m.ID).Count(&count).Error
	return count > 0, err
}

// GetChildMenus 获取子菜单列表
func (m *MenuDO) GetChildMenus() ([]MenuDO, error) {
	var children []MenuDO
	err := m.Query().Where("parent_id = ?", m.ID).Order("sort").Find(&children).Error
	return children, err
}
