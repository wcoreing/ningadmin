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

// TableName 表名
func (MenuDO) TableName() string {
	return "sys_permission"
}

// Query 获取查询构建器
func (m *MenuDO) Query() *gorm.DB {
	return db.DB.WnAdmin.Model(m)
}

// Create 创建记录
func (m *MenuDO) Create() error {
	return m.Query().Create(m).Error
}

// Update 更新记录
func (m *MenuDO) Update() error {
	return m.Query().Updates(m).Error
}

// Delete 删除记录
func (m *MenuDO) Delete() error {
	return m.Query().Update("del_flag", "2").Error
}

// First 获取第一条记录
func (m *MenuDO) First() error {
	return m.Query().Where("del_flag = ?", "0").First(m).Error
}

// Find 查询多条记录
func (m *MenuDO) Find(dest interface{}) error {
	return m.Query().Where("del_flag = ?", "0").Find(dest).Error
}

// CheckPermissionExists 检查权限标识是否存在
func (m *MenuDO) CheckPermissionExists(permission string, excludeID int64) (bool, error) {
	var count int64
	db := m.Query().Where("permission_code = ? AND del_flag = ?", permission, "0")
	if excludeID > 0 {
		db = db.Where("permission_id != ?", excludeID)
	}
	err := db.Count(&count).Error
	return count > 0, err
}
