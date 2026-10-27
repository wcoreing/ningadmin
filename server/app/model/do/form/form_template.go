package form

import (
	"time"
)

// FormTemplate 表单模板
type FormTemplate struct {
	ID             int64      `json:"id" gorm:"primary_key"`
	Name           string     `json:"name" gorm:"not null"`
	Description    string     `json:"description"`
	SubmitURL      string     `json:"submit_url" gorm:"not null"`
	SubmitMethod   string     `json:"submit_method" gorm:"not null"`
	SuccessMessage string     `json:"success_message"`
	FailMessage    string     `json:"fail_message"`
	FormItems      []byte     `json:"form_items"`
	Config         []byte     `json:"config"`
	Status         string     `json:"status" gorm:"default:0"`
	DelFlag        string     `json:"del_flag" gorm:"default:0"`
	CreateBy       string     `json:"create_by"`
	CreateTime     time.Time  `json:"create_time"`
	UpdateBy       string     `json:"update_by"`
	UpdateTime     *time.Time `json:"update_time"`
}

// TableName 表名
func (FormTemplate) TableName() string {
	return "form_template"
}
