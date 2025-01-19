package dto

import (
	"encoding/json"
	"time"
)

// FormTemplateQuery 表单模板查询参数
type FormTemplateQuery struct {
	Name     string `form:"name"`     // 模板名称
	PageNum  int    `form:"pageNum"`  // 页码
	PageSize int    `form:"pageSize"` // 每页条数
}

// FormTemplateDTO 表单模板DTO
type FormTemplateDTO struct {
	ID          int64           `json:"id"`          // 模板ID
	Name        string          `json:"name"`        // 模板名称
	Description string          `json:"description"` // 模板描述
	FormItems   json.RawMessage `json:"formItems"`   // 表单项配置
	Config      json.RawMessage `json:"config"`      // 表单配置
	CreateTime  time.Time       `json:"createTime"`  // 创建时间
	UpdateTime  *time.Time      `json:"updateTime"`  // 更新时间
	CreateBy    string          `json:"createBy"`    // 创建者
	UpdateBy    string          `json:"updateBy"`    // 更新者
}

// FormShareCreateDTO 创建表单分享DTO
type FormShareCreateDTO struct {
	TemplateID int64      `json:"templateId"` // 模板ID
	ExpireTime *time.Time `json:"expireTime"` // 过期时间
}

// FormShareDTO 表单分享DTO
type FormShareDTO struct {
	ShareCode string `json:"shareCode"` // 分享码
	ShareURL  string `json:"shareUrl"`  // 分享链接
}

// FormSubmissionDTO 表单提交DTO
type FormSubmissionDTO struct {
	TemplateID int64           `json:"templateId"` // 模板ID
	FormData   json.RawMessage `json:"formData"`   // 表单数据
}

// FormSubmissionQuery 表单提交记录查询参数
type FormSubmissionQuery struct {
	TemplateID int64 `form:"templateId"` // 模板ID
	PageNum    int   `form:"pageNum"`    // 页码
	PageSize   int   `form:"pageSize"`   // 每页条数
}

// FormSubmissionVO 表单提交记录VO
type FormSubmissionVO struct {
	ID         int64           `json:"id"`         // 提交ID
	TemplateID int64           `json:"templateId"` // 模板ID
	FormData   json.RawMessage `json:"formData"`   // 表单数据
	SubmitTime time.Time       `json:"submitTime"` // 提交时间
	SubmitBy   string          `json:"submitBy"`   // 提交人
}
