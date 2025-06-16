package do

import (
	"time"
	"wnadmin/app/provider/db"

	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID            int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Openid        string     `gorm:"column:openid;size:255" json:"openid"`
	Username      string     `gorm:"column:username;size:50;not null" json:"username"`
	Password      string     `gorm:"column:password;size:255;not null" json:"password"`
	Nickname      string     `gorm:"column:nickname;size:50" json:"nickname"`
	Gender        *int       `gorm:"column:gender;type:tinyint(1)" json:"gender"`
	Avatar        string     `gorm:"column:avatar;size:255" json:"avatar"`
	Phone         string     `gorm:"column:phone;size:20" json:"phone"`
	Email         string     `gorm:"column:email;size:100" json:"email"`
	Birthday      string     `gorm:"column:birthday;size:50" json:"birthday"`
	Balance       float64    `gorm:"column:balance;type:decimal(20,2);default:0.00;not null" json:"balance"`
	Points        int        `gorm:"column:points;default:0;not null" json:"points"`
	Level         int        `gorm:"column:level;type:tinyint(4);default:1;not null" json:"level"`
	TotalAmount   float64    `gorm:"column:total_amount;type:decimal(20,2);default:0.00;not null" json:"totalAmount"`
	LastLoginTime *time.Time `gorm:"column:last_login_time" json:"lastLoginTime"`
	LastLoginIP   string     `gorm:"column:last_login_ip;size:50" json:"lastLoginIp"`
	Status        int        `gorm:"column:status;type:tinyint(4);default:1;not null" json:"status"`
	IsDelete      int        `gorm:"column:is_delete;type:tinyint(4);default:0;not null" json:"isDelete"`
	CreateTime    time.Time  `gorm:"column:create_time;default:CURRENT_TIMESTAMP" json:"createTime"`
	UpdateTime    time.Time  `gorm:"column:update_time;default:CURRENT_TIMESTAMP" json:"updateTime"`
}

// TableName 指定表名
func (u *User) TableName() string {
	return "user"
}

// Query 获取查询构建器
func (u *User) Query() *gorm.DB {
	return db.WnAdmin().Model(u).Where("is_delete = ?", 0)
}

// Create 创建用户
func (u *User) Create() error {
	return db.WnAdmin().Create(u).Error
}

// Update 更新用户
func (u *User) Update() error {
	return db.WnAdmin().Save(u).Error
}

// Delete 软删除用户
func (u *User) Delete() error {
	return db.WnAdmin().Model(u).Where("id = ?", u.ID).Update("is_delete", 1).Error
}

// GetByID 根据ID获取用户
func (u *User) GetByID(id int64) error {
	return u.Query().Where("id = ?", id).First(u).Error
}

// GetByUsername 根据用户名获取用户
func (u *User) GetByUsername(username string) error {
	return u.Query().Where("username = ?", username).First(u).Error
}

// GetByPhone 根据手机号获取用户
func (u *User) GetByPhone(phone string) error {
	return u.Query().Where("phone = ?", phone).First(u).Error
}

// GetByEmail 根据邮箱获取用户
func (u *User) GetByEmail(email string) error {
	return u.Query().Where("email = ?", email).First(u).Error
}

// CheckUsernameExists 检查用户名是否存在
func (u *User) CheckUsernameExists(username string, excludeID int64) (bool, error) {
	var count int64
	query := u.Query().Where("username = ?", username)
	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
	}
	err := query.Count(&count).Error
	return count > 0, err
}

// CheckPhoneExists 检查手机号是否存在
func (u *User) CheckPhoneExists(phone string, excludeID int64) (bool, error) {
	var count int64
	query := u.Query().Where("phone = ?", phone)
	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
	}
	err := query.Count(&count).Error
	return count > 0, err
}

// CheckEmailExists 检查邮箱是否存在
func (u *User) CheckEmailExists(email string, excludeID int64) (bool, error) {
	var count int64
	query := u.Query().Where("email = ?", email)
	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
	}
	err := query.Count(&count).Error
	return count > 0, err
}

// GetList 获取用户列表
func (u *User) GetList(offset, limit int, conditions map[string]interface{}) ([]User, int64, error) {
	var users []User
	var total int64

	query := u.Query()

	// 添加查询条件
	for key, value := range conditions {
		if value != nil && value != "" {
			switch key {
			case "keyword":
				keyword := "%" + value.(string) + "%"
				query = query.Where("username LIKE ? OR nickname LIKE ? OR phone LIKE ? OR email LIKE ?",
					keyword, keyword, keyword, keyword)
			case "status":
				query = query.Where("status = ?", value)
			case "level":
				query = query.Where("level = ?", value)
			case "gender":
				query = query.Where("gender = ?", value)
			}
		}
	}

	// 获取总数
	query.Count(&total)

	// 获取列表
	err := query.Offset(offset).Limit(limit).Order("create_time DESC").Find(&users).Error

	return users, total, err
}

// UpdateStatus 更新用户状态
func (u *User) UpdateStatus(id int64, status int) error {
	return db.WnAdmin().Model(u).Where("id = ?", id).Update("status", status).Error
}

// UpdateLastLogin 更新最后登录信息
func (u *User) UpdateLastLogin(id int64, loginTime time.Time, loginIP string) error {
	return db.WnAdmin().Model(u).Where("id = ?", id).Updates(map[string]interface{}{
		"last_login_time": loginTime,
		"last_login_ip":   loginIP,
	}).Error
}

// UpdateBalance 更新用户余额
func (u *User) UpdateBalance(id int64, balance float64) error {
	return db.WnAdmin().Model(u).Where("id = ?", id).Update("balance", balance).Error
}

// UpdatePoints 更新用户积分
func (u *User) UpdatePoints(id int64, points int) error {
	return db.WnAdmin().Model(u).Where("id = ?", id).Update("points", points).Error
}

// UpdateLevel 更新用户等级
func (u *User) UpdateLevel(id int64, level int) error {
	return db.WnAdmin().Model(u).Where("id = ?", id).Update("level", level).Error
}

// UpdatePassword 更新用户密码
func (u *User) UpdatePassword(id int64, password string) error {
	return db.WnAdmin().Model(u).Where("id = ?", id).Update("password", password).Error
}
