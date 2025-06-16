package shop

import (
	"time"
	"wnadmin/app/provider/db"

	"gorm.io/gorm"
)

// Feedback 反馈模型
type Feedback struct {
	ID         int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserID     int64      `gorm:"column:user_id;not null" json:"userId"`
	Content    string     `gorm:"column:content;type:text;not null" json:"content"`
	Contact    string     `gorm:"column:contact;size:100" json:"contact"`
	Status     int        `gorm:"column:status;default:0;not null" json:"status"`
	Reply      string     `gorm:"column:reply;type:text" json:"reply"`
	ReplyTime  *time.Time `gorm:"column:reply_time" json:"replyTime"`
	CreateTime time.Time  `gorm:"column:create_time;default:CURRENT_TIMESTAMP" json:"createTime"`
	UpdateTime time.Time  `gorm:"column:update_time;default:CURRENT_TIMESTAMP" json:"updateTime"`
}

// TableName 指定表名
func (f *Feedback) TableName() string {
	return "feedback"
}

// Query 获取查询构建器
func (f *Feedback) Query() *gorm.DB {
	return db.WnAdmin().Model(f)
}

// Create 创建反馈
func (f *Feedback) Create() error {
	return db.WnAdmin().Create(f).Error
}

// Update 更新反馈
func (f *Feedback) Update() error {
	return db.WnAdmin().Save(f).Error
}

// Delete 删除反馈
func (f *Feedback) Delete() error {
	return db.WnAdmin().Delete(f).Error
}

// GetByID 根据ID获取反馈
func (f *Feedback) GetByID(id int64) error {
	return db.WnAdmin().Where("id = ?", id).First(f).Error
}

// GetList 获取反馈列表
func (f *Feedback) GetList(offset, limit int, conditions map[string]interface{}) ([]Feedback, int64, error) {
	var feedbacks []Feedback
	var total int64

	query := db.WnAdmin().Model(f)

	// 添加查询条件
	for key, value := range conditions {
		if value != nil && value != "" {
			switch key {
			case "user_id":
				query = query.Where("user_id = ?", value)
			case "status":
				query = query.Where("status = ?", value)
			case "keyword":
				query = query.Where("content LIKE ?", "%"+value.(string)+"%")
			}
		}
	}

	// 获取总数
	query.Count(&total)

	// 获取列表
	err := query.Offset(offset).Limit(limit).Order("create_time DESC").Find(&feedbacks).Error

	return feedbacks, total, err
}
