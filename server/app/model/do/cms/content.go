package cms

import (
	"time"
	"wnadmin/app/provider/db"

	"gorm.io/gorm"
)

// Content 栏目下的内容
type Content struct {
	ID             int64     `gorm:"column:id;primaryKey" json:"id"`
	ChannelID      int64     `gorm:"column:channel_id" json:"channelId"`
	Title          string    `gorm:"column:title" json:"title"`
	Summary        string    `gorm:"column:summary" json:"summary"`
	Cover          string    `gorm:"column:cover" json:"cover"`
	Body           string    `gorm:"column:body" json:"body"`
	SeoTitle       string    `gorm:"column:seo_title" json:"seoTitle"`
	SeoKeywords    string    `gorm:"column:seo_keywords" json:"seoKeywords"`
	SeoDescription string    `gorm:"column:seo_description" json:"seoDescription"`
	Sort           int       `gorm:"column:sort" json:"sort"`
	Status         int       `gorm:"column:status" json:"status"`
	IsDelete       int       `gorm:"column:is_delete" json:"-"`
	PublishTime    time.Time `gorm:"column:publish_time" json:"publishTime"`
	CreateTime     time.Time `gorm:"column:create_time" json:"createTime"`
	UpdateTime     time.Time `gorm:"column:update_time" json:"updateTime"`
}

// TableName 表名
func (Content) TableName() string {
	return "cms_content"
}

// Query 获取查询构建器
func (c *Content) Query() *gorm.DB {
	return db.WnAdmin().Model(c).Where("is_delete = ?", 0)
}

// Create 新增内容。显式选择字段，避免草稿状态 0 被写成已发布。
func (c *Content) Create() error {
	return db.WnAdmin().Select(
		"ChannelID", "Title", "Summary", "Cover", "Body",
		"SeoTitle", "SeoKeywords", "SeoDescription", "Sort", "Status",
		"PublishTime", "CreateTime", "UpdateTime",
	).Create(c).Error
}

// Save 保存内容
func (c *Content) Save() error {
	return db.WnAdmin().Save(c).Error
}
