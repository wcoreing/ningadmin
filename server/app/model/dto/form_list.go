package dto

import (
	"encoding/json"
)

// FormListQuery 表单列表查询参数
type FormListQuery struct {
	FormName   string `form:"formName"`   // 表单名称
	FormStatus string `form:"formStatus"` // 表单状态
	PageNum    int    `form:"pageNum"`    // 页码
	PageSize   int    `form:"pageSize"`   // 每页条数
}

// FormCreateDTO 创建表单请求
type FormCreateDTO struct {
	TemplateID      int64           `json:"templateId" binding:"required"` // 模板ID
	FormName        string          `json:"formName" binding:"required"`   // 表单名称
	FormDescription string          `json:"formDescription"`               // 表单描述
	FormAttributes  json.RawMessage `json:"formAttributes"`                // 表单属性
	FormStatus      string          `json:"formStatus" binding:"required"` // 表单状态
}

// FormUpdateDTO 更新表单请求
type FormUpdateDTO struct {
	ID              int64           `json:"id" binding:"required"`         // 表单ID
	TemplateID      int64           `json:"templateId" binding:"required"` // 模板ID
	FormName        string          `json:"formName" binding:"required"`   // 表单名称
	FormDescription string          `json:"formDescription"`               // 表单描述
	FormAttributes  json.RawMessage `json:"formAttributes"`                // 表单属性
	FormStatus      string          `json:"formStatus" binding:"required"` // 表单状态
}
