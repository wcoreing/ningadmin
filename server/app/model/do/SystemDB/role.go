package SystemDB

import (
	"time"
	"wnadmin/app/provider/db"

	"gorm.io/gorm"
)

// RoleDO 角色信息
type RoleDO struct {
	ID         int64     `json:"roleId" gorm:"column:role_id;primaryKey;comment:角色ID"`
	RoleName   string    `json:"roleName" gorm:"column:role_name;comment:角色名称"`
	RoleKey    string    `json:"roleKey" gorm:"column:role_key;comment:角色标识"`
	Status     string    `json:"status" gorm:"comment:角色状态（0正常 1停用）"`
	DelFlag    string    `json:"-" gorm:"column:del_flag;comment:删除标志（0代表存在 2代表删除）"`
	CreateTime time.Time `json:"createTime" gorm:"column:create_time;autoCreateTime;comment:创建时间"`
	CreateBy   string    `json:"createBy" gorm:"column:create_by;comment:创建者"`
	UpdateTime time.Time `json:"updateTime" gorm:"column:update_time;autoUpdateTime;comment:更新时间"`
	UpdateBy   string    `json:"updateBy" gorm:"column:update_by;comment:更新者"`
}

// TableName 表名
func (RoleDO) TableName() string {
	return "sys_role"
}

// Query 获取查询构建器
func (r *RoleDO) Query() *gorm.DB {
	return db.DB.WnAdmin.Model(r).Where("del_flag = ?", "0")
}

// Create 创建记录
func (r *RoleDO) Create() error {
	return r.Query().Create(r).Error
}

// Update 更新记录
func (r *RoleDO) Update() error {
	return r.Query().Updates(r).Error
}

// Delete 删除记录
func (r *RoleDO) Delete() error {
	return db.DB.WnAdmin.Transaction(func(tx *gorm.DB) error {
		// 删除角色权限关联
		if err := tx.Table("sys_role_permission").Where("role_id = ?", r.ID).Delete(nil).Error; err != nil {
			return err
		}

		// 软删除角色
		return tx.Model(r).Update("del_flag", "2").Error
	})
}

// First 获取第一条记录
func (r *RoleDO) First() error {
	return r.Query().First(r).Error
}

// Find 查询多条记录
func (r *RoleDO) Find(dest interface{}) error {
	return r.Query().Find(dest).Error
}

// CheckRoleKeyExists 检查角色标识是否存在
func (r *RoleDO) CheckRoleKeyExists(roleKey string, excludeID int64) (bool, error) {
	var count int64
	db := r.Query().Where("role_key = ? ", roleKey)
	if excludeID > 0 {
		db = db.Where("role_id != ?", excludeID)
	}
	err := db.Count(&count).Error
	return count > 0, err
}

// UpdateStatus 更新角色状态
func (r *RoleDO) UpdateStatus() error {
	return r.Query().Where("role_id = ?", r.ID).Update("status", r.Status).Error
}

// GetRolePermissions 获取角色权限ID列表
func (r *RoleDO) GetRolePermissions() ([]int64, error) {
	var permissionIDs []int64
	err := db.DB.WnAdmin.Table("sys_role_permission").
		Where("role_id = ?", r.ID).
		Pluck("permission_id", &permissionIDs).Error
	return permissionIDs, err
}

// UpdateRolePermissions 更新角色权限
func (r *RoleDO) UpdateRolePermissions(permissionIDs []int64) error {
	return db.DB.WnAdmin.Transaction(func(tx *gorm.DB) error {
		// 删除原有权限
		if err := tx.Table("sys_role_permission").Where("role_id = ?", r.ID).Delete(nil).Error; err != nil {
			return err
		}

		// 添加新权限
		if len(permissionIDs) > 0 {
			rolePermissions := make([]map[string]interface{}, len(permissionIDs))
			for i, pid := range permissionIDs {
				rolePermissions[i] = map[string]interface{}{
					"role_id":       r.ID,
					"permission_id": pid,
				}
			}
			return tx.Table("sys_role_permission").Create(rolePermissions).Error
		}
		return nil
	})
}
