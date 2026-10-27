package SystemDB

import (
	"wnadmin/app/provider/db"

	"gorm.io/gorm"
)

// UserDeptDO 用户部门关系
type UserDeptDO struct {
	UserId    int64  `json:"userId" gorm:"column:user_id;primaryKey;comment:用户ID"`
	DeptId    int64  `json:"deptId" gorm:"column:dept_id;primaryKey;comment:部门ID"`
	DataScope string `json:"dataScope" gorm:"column:data_scope;comment:数据范围（1：全部数据权限 2：自定数据权限 3：本部门数据权限 4：本部门及以下数据权限 5：仅本人数据权限）"`
}

// TableName 表名
func (UserDeptDO) TableName() string {
	return "sys_user_dept"
}

// Query 获取查询构建器
func (u *UserDeptDO) Query() *gorm.DB {
	return db.WnAdmin().Model(u)
}

// Create 创建记录
func (u *UserDeptDO) Create() error {
	return u.Query().Create(u).Error
}

// Delete 删除记录
func (u *UserDeptDO) Delete() error {
	return u.Query().Delete(u).Error
}

// DeleteByUserId 根据用户ID删除关系
func (u *UserDeptDO) DeleteByUserId(userId int64) error {
	return u.Query().Where("user_id = ?", userId).Delete(u).Error
}

// DeleteByDeptId 根据部门ID删除关系
func (u *UserDeptDO) DeleteByDeptId(deptId int64) error {
	return u.Query().Where("dept_id = ?", deptId).Delete(u).Error
}

// GetUsersByDeptId 根据部门ID获取用户列表
func (u *UserDeptDO) GetUsersByDeptId(deptId int64) ([]int64, error) {
	var userIds []int64
	err := u.Query().Where("dept_id = ?", deptId).Pluck("user_id", &userIds).Error
	return userIds, err
}

// GetDeptsByUserId 根据用户ID获取部门列表
func (u *UserDeptDO) GetDeptsByUserId(userId int64) ([]int64, error) {
	var deptIds []int64
	err := u.Query().Where("user_id = ?", userId).Pluck("dept_id", &deptIds).Error
	return deptIds, err
}

// BatchCreate 批量创建关系
func (u *UserDeptDO) BatchCreate(relations []UserDeptDO) error {
	if len(relations) == 0 {
		return nil
	}
	return u.Query().Create(relations).Error
}

// GetUserDeptWithScope 根据用户ID获取部门和数据权限信息
func (u *UserDeptDO) GetUserDeptWithScope(userId int64) ([]UserDeptDO, error) {
	var relations []UserDeptDO
	err := u.Query().Where("user_id = ?", userId).Find(&relations).Error
	return relations, err
}

// UpdateDataScope 更新用户在某个部门的数据权限
func (u *UserDeptDO) UpdateDataScope(userId, deptId int64, dataScope string) error {
	return u.Query().Where("user_id = ? AND dept_id = ?", userId, deptId).Update("data_scope", dataScope).Error
}
