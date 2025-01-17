package SystemDB

import (
	"time"
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
	return db.DB.WnAdmin.Model(u)
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
	return u.Query().Update("del_flag", "2").Error
}

// First 获取第一条记录
func (u *UserDO) First() error {
	return u.Query().Where("del_flag = ?", "0").First(u).Error
}

// Find 查询多条记录
func (u *UserDO) Find(dest interface{}) error {
	return u.Query().Where("del_flag = ?", "0").Find(dest).Error
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
