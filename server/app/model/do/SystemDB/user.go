package SystemDB

import (
	"time"
	"wnadmin/app/model/dto"
	"wnadmin/app/provider/db"

	"gorm.io/gorm"
)

// UserDO 用户信息
type UserDO struct {
	ID          int64     `json:"userId" gorm:"column:user_id;primaryKey;comment:用户ID"`
	UserName    string    `json:"userName" gorm:"column:user_name;comment:用户账号"`
	Password    string    `json:"password" gorm:"column:password;comment:密码"`
	NickName    string    `json:"nickName" gorm:"column:nick_name;comment:用户昵称"`
	Email       string    `json:"email" gorm:"column:email;comment:用户邮箱"`
	Phonenumber string    `json:"phonenumber" gorm:"column:phonenumber;comment:手机号码"`
	Sex         string    `json:"sex" gorm:"column:sex;comment:用户性别（0男 1女 2未知）"`
	Avatar      string    `json:"avatar" gorm:"column:avatar;comment:头像地址"`
	Status      string    `json:"status" gorm:"column:status;comment:帐号状态（0正常 1停用）"`
	DelFlag     string    `json:"-" gorm:"column:del_flag;comment:删除标志（0代表存在 2代表删除）"`
	CreateTime  time.Time `json:"createTime" gorm:"column:create_time;autoCreateTime;comment:创建时间"`
	CreateBy    string    `json:"createBy" gorm:"column:create_by;comment:创建者"`
	UpdateTime  time.Time `json:"updateTime" gorm:"column:update_time;autoUpdateTime;comment:更新时间"`
	UpdateBy    string    `json:"updateBy" gorm:"column:update_by;comment:更新者"`
}

// TableName 表名
func (UserDO) TableName() string {
	return "sys_user"
}

// Query 获取查询构建器
func (u *UserDO) Query() *gorm.DB {
	return db.WnAdmin().Model(u).Where("del_flag = ?", "0")
}

// Create 创建记录
func (u *UserDO) Create() error {
	return u.Query().Create(u).Error
}

// Update 更新记录
func (u *UserDO) Update() error {
	return u.Query().Updates(u).Error
}

// Delete 删除记录
func (u *UserDO) Delete() error {
	return db.WnAdmin().Transaction(func(tx *gorm.DB) error {
		// 删除用户角色关联
		if err := tx.Table("sys_user_role").Where("user_id = ?", u.ID).Delete(nil).Error; err != nil {
			return err
		}

		// 软删除用户
		return tx.Model(u).Where("user_id = ?", u.ID).Update("del_flag", "2").Error
	})
}

// First 获取第一条记录
func (u *UserDO) First() error {
	return u.Query().First(u).Error
}

// Find 查询多条记录
func (u *UserDO) Find(dest interface{}) error {
	return u.Query().Find(dest).Error
}

// CheckUserNameExists 检查用户名是否存在
func (u *UserDO) CheckUserNameExists(userName string, excludeID int64) (bool, error) {
	var count int64
	db := u.Query().Where("user_name = ? AND del_flag = ?", userName, "0")
	if excludeID > 0 {
		db = db.Where("user_id != ?", excludeID)
	}
	err := db.Count(&count).Error
	return count > 0, err
}

// UpdateStatus 更新用户状态
func (u *UserDO) UpdateStatus() error {
	return u.Query().Where("user_id = ?", u.ID).Update("status", u.Status).Error
}

// UpdatePassword 更新用户密码
func (u *UserDO) UpdatePassword() error {
	return u.Query().Where("user_id = ?", u.ID).Update("password", u.Password).Error
}

// GetUserRoles 获取用户关联的角色ID列表
func (u *UserDO) GetUserRoles() ([]int64, error) {
	var roleIDs []int64
	err := db.WnAdmin().Table("sys_user_role").
		Where("user_id = ?", u.ID).
		Pluck("role_id", &roleIDs).Error
	return roleIDs, err
}

// UpdateUserRoles 更新用户关联的角色
func (u *UserDO) UpdateUserRoles(roleIDs []int64) error {
	return db.WnAdmin().Transaction(func(tx *gorm.DB) error {
		// 删除原有角色关联
		if err := tx.Table("sys_user_role").Where("user_id = ?", u.ID).Delete(nil).Error; err != nil {
			return err
		}

		// 添加新角色关联
		if len(roleIDs) > 0 {
			userRoles := make([]map[string]interface{}, len(roleIDs))
			for i, rid := range roleIDs {
				userRoles[i] = map[string]interface{}{
					"user_id": u.ID,
					"role_id": rid,
				}
			}
			return tx.Table("sys_user_role").Create(userRoles).Error
		}
		return nil
	})
}

// GetUserRolesWithName 获取用户角色（包含角色名称）
func (u *UserDO) GetUserRolesWithName() ([]int64, error) {
	var roleIDs []int64
	err := db.WnAdmin().Table("sys_user_role").
		Select("sys_role.role_id").
		Joins("LEFT JOIN sys_role ON sys_user_role.role_id = sys_role.role_id").
		Where("sys_user_role.user_id = ? AND sys_role.del_flag = ?", u.ID, "0").
		Pluck("sys_role.role_id", &roleIDs).Error
	return roleIDs, err
}

// GetUserRole 获取用户角色信息
func (u *UserDO) GetUserRole() (*dto.RoleInfo, error) {
	var role dto.RoleInfo
	err := db.WnAdmin().Table("sys_user_role").
		Select("sys_role.role_id, sys_role.role_name, sys_role.role_key, sys_role.status").
		Joins("LEFT JOIN sys_role ON sys_user_role.role_id = sys_role.role_id").
		Where("sys_user_role.user_id = ? AND sys_role.del_flag = ?", u.ID, "0").
		First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// UpdateUserRole 更新用户角色
func (u *UserDO) UpdateUserRole(roleID int64) error {
	return db.WnAdmin().Transaction(func(tx *gorm.DB) error {
		// 删除原有角色关联
		if err := tx.Table("sys_user_role").Where("user_id = ?", u.ID).Delete(nil).Error; err != nil {
			return err
		}

		// 添加新角色关联
		return tx.Table("sys_user_role").Create(map[string]interface{}{
			"user_id": u.ID,
			"role_id": roleID,
		}).Error
	})
}

// GetUserDepts 获取用户关联的部门信息
func (u *UserDO) GetUserDepts() ([]UserDeptDO, error) {
	userDeptDO := &UserDeptDO{}
	return userDeptDO.GetUserDeptWithScope(u.ID)
}

// GetUserDeptNames 获取用户关联的部门名称
func (u *UserDO) GetUserDeptNames() ([]string, error) {
	var deptNames []string
	err := db.WnAdmin().Table("sys_user_dept").
		Select("sys_dept.dept_name").
		Joins("LEFT JOIN sys_dept ON sys_user_dept.dept_id = sys_dept.dept_id").
		Where("sys_user_dept.user_id = ? AND sys_dept.del_flag = ?", u.ID, "0").
		Pluck("sys_dept.dept_name", &deptNames).Error
	return deptNames, err
}

// GetUserDeptIds 获取用户关联的部门ID列表
func (u *UserDO) GetUserDeptIds() ([]int64, error) {
	userDeptDO := &UserDeptDO{}
	return userDeptDO.GetDeptsByUserId(u.ID)
}
