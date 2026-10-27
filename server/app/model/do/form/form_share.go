package form

import "time"

// FormShare 表单分享记录
type FormShare struct {
	ID         int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`       // 分享ID
	TemplateID int64      `gorm:"column:template_id;not null" json:"templateId"`      // 模板ID
	ShareCode  string     `gorm:"column:share_code;not null;unique" json:"shareCode"` // 分享码
	ExpireTime *time.Time `gorm:"column:expire_time" json:"expireTime"`               // 过期时间
	CreateTime time.Time  `gorm:"column:create_time;not null" json:"createTime"`      // 创建时间
	CreateBy   string     `gorm:"column:create_by" json:"createBy"`                   // 创建者
}

// TableName 表名
func (FormShare) TableName() string {
	return "form_share"
}
