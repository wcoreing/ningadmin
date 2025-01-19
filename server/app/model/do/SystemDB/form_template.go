package SystemDB

import (
	"encoding/json"
	"time"
)

// FormTemplate 表单模板
type FormTemplate struct {
	ID          int64           `gorm:"column:id;primaryKey;autoIncrement" json:"id"`          // 模板ID
	Name        string          `gorm:"column:name;not null" json:"name"`                      // 模板名称
	Description string          `gorm:"column:description" json:"description"`                 // 模板描述
	FormItems   json.RawMessage `gorm:"column:form_items;type:json;not null" json:"formItems"` // 表单项配置
	Config      json.RawMessage `gorm:"column:config;type:json;not null" json:"config"`        // 表单配置
	CreateTime  time.Time       `gorm:"column:create_time;not null" json:"createTime"`         // 创建时间
	UpdateTime  *time.Time      `gorm:"column:update_time" json:"updateTime"`                  // 更新时间
	CreateBy    string          `gorm:"column:create_by" json:"createBy"`                      // 创建者
	UpdateBy    string          `gorm:"column:update_by" json:"updateBy"`                      // 更新者
}

// TableName 表名
func (FormTemplate) TableName() string {
	return "form_template"
}
