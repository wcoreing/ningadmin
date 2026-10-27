package site

import (
	"time"
	"wnadmin/app/provider/db"

	"gorm.io/gorm"
)

// Profile 官网企业资料，全站只有一条
type Profile struct {
	ID             int64     `gorm:"column:id;primaryKey" json:"id"`
	Name           string    `gorm:"column:name" json:"name"`
	ShortName      string    `gorm:"column:short_name" json:"shortName"`
	Phone          string    `gorm:"column:phone" json:"phone"`
	ContactName    string    `gorm:"column:contact_name" json:"contactName"`
	Address        string    `gorm:"column:address" json:"address"`
	Email          string    `gorm:"column:email" json:"email"`
	Summary        string    `gorm:"column:summary" json:"summary"`
	SeoTitle       string    `gorm:"column:seo_title" json:"seoTitle"`
	SeoKeywords    string    `gorm:"column:seo_keywords" json:"seoKeywords"`
	SeoDescription string    `gorm:"column:seo_description" json:"seoDescription"`
	CreateTime     time.Time `gorm:"column:create_time" json:"createTime"`
	UpdateTime     time.Time `gorm:"column:update_time" json:"updateTime"`
}

// TableName 表名
func (Profile) TableName() string {
	return "site_profile"
}

// Query 获取查询构建器
func (p *Profile) Query() *gorm.DB {
	return db.WnAdmin().Model(p)
}

// UpdateFields 更新企业资料
func (p *Profile) UpdateFields(fields map[string]any) error {
	return db.WnAdmin().Model(p).Where("id = ?", p.ID).Updates(fields).Error
}
