package SystemDB

import (
	"encoding/json"
	"time"
)

// FormSubmission 表单提交记录
type FormSubmission struct {
	ID         int64           `gorm:"column:id;primaryKey;autoIncrement" json:"id"`        // 提交ID
	TemplateID int64           `gorm:"column:template_id;not null" json:"templateId"`       // 模板ID
	FormData   json.RawMessage `gorm:"column:form_data;type:json;not null" json:"formData"` // 表单数据
	SubmitTime time.Time       `gorm:"column:submit_time;not null" json:"submitTime"`       // 提交时间
	SubmitBy   string          `gorm:"column:submit_by" json:"submitBy"`                    // 提交人
}

// TableName 表名
func (FormSubmission) TableName() string {
	return "form_submission"
}
