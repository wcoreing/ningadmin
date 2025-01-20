package form

import (
	"encoding/json"
	"time"
)

// Form 表单实体
type Form struct {
	ID              int64           `gorm:"column:id;primary_key;auto_increment" json:"id"`              // 表单ID
	TemplateID      int64           `gorm:"column:template_id;not null" json:"templateId"`               // 模板ID
	TemplateName    string          `gorm:"column:template_name;not null" json:"templateName"`           // 模板名称
	FormName        string          `gorm:"column:form_name;not null" json:"formName"`                   // 表单名称
	FormDescription string          `gorm:"column:form_description" json:"formDescription"`              // 表单描述
	FormAttributes  json.RawMessage `gorm:"column:form_attributes;type:json" json:"formAttributes"`      // 表单属性
	FormStatus      string          `gorm:"column:form_status;default:0" json:"formStatus"`              // 表单状态（0正常 1停用）
	FormDelFlag     string          `gorm:"column:form_del_flag;default:0" json:"formDelFlag"`           // 表单删除标志（0存在 2删除）
	FormCreateBy    string          `gorm:"column:form_create_by" json:"formCreateBy"`                   // 表单创建者
	FormCreateTime  *time.Time      `gorm:"column:form_create_time;type:datetime" json:"formCreateTime"` // 表单创建时间
	FormUpdateBy    string          `gorm:"column:form_update_by" json:"formUpdateBy"`                   // 表单更新者
	FormUpdateTime  *time.Time      `gorm:"column:form_update_time;type:datetime" json:"formUpdateTime"` // 表单更新时间
}

// TableName 设置表名
func (Form) TableName() string {
	return "form"
}
