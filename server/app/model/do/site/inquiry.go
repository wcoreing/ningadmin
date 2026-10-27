package site

import (
	"time"
	"wnadmin/app/provider/db"

	"gorm.io/gorm"
)

// Inquiry 官网询盘
type Inquiry struct {
	ID         int64     `gorm:"column:id;primaryKey" json:"id"`
	GoodsID    int64     `gorm:"column:goods_id" json:"goodsId"`
	GoodsName  string    `gorm:"column:goods_name" json:"goodsName"`
	Name       string    `gorm:"column:name" json:"name"`
	Company    string    `gorm:"column:company" json:"company"`
	Phone      string    `gorm:"column:phone" json:"phone"`
	Email      string    `gorm:"column:email" json:"email"`
	Country    string    `gorm:"column:country" json:"country"`
	Quantity   string    `gorm:"column:quantity" json:"quantity"`
	Message    string    `gorm:"column:message" json:"message"`
	Status     int       `gorm:"column:status" json:"status"`
	IsDelete   int       `gorm:"column:is_delete" json:"-"`
	CreateTime time.Time `gorm:"column:create_time;autoCreateTime" json:"createTime"`
}

// TableName 表名
func (Inquiry) TableName() string {
	return "site_inquiry"
}

// Query 获取查询构建器
func (m *Inquiry) Query() *gorm.DB {
	return db.WnAdmin().Model(m).Where("is_delete = ?", 0)
}

// Create 写入询盘
func (m *Inquiry) Create() error {
	return db.WnAdmin().Create(m).Error
}

func (m *Inquiry) scope(status string) *gorm.DB {
	query := m.Query()
	if status == "0" || status == "1" {
		query = query.Where("status = ?", status)
	}
	return query
}

// Page 分页查询询盘
func (m *Inquiry) Page(status string, page, size int) ([]Inquiry, int64, error) {
	var total int64
	if err := m.scope(status).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []Inquiry
	err := m.scope(status).Order("id desc").Offset((page - 1) * size).Limit(size).Find(&rows).Error
	return rows, total, err
}

// Mark 更新处理状态
func (m *Inquiry) Mark(id int64, status int) error {
	return m.Query().Where("id = ?", id).Update("status", status).Error
}
