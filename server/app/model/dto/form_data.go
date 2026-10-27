package dto

import (
	"encoding/json"
	"time"
)

// FormDataQuery 表单数据查询参数
type FormDataQuery struct {
	FormName  string    `form:"formName"`  // 表单名称
	StartTime time.Time `form:"startTime"` // 开始时间
	EndTime   time.Time `form:"endTime"`   // 结束时间
	PageNum   int       `form:"pageNum"`   // 页码
	PageSize  int       `form:"pageSize"`  // 每页数量
}

// FormSubmitRequest 提交表单请求
type FormSubmitRequest struct {
	FormID    int64           `json:"formId" binding:"required"` // 表单ID
	Data      json.RawMessage `json:"data" binding:"required"`   // 表单数据
	Submitter json.RawMessage `json:"submitter"`                 // 提交者信息
}

// FormDataResponse 表单数据返回结构
type FormDataResponse struct {
	ID         int64           `json:"id"`         // 数据ID
	FormID     int64           `json:"formId"`     // 表单ID
	FormName   string          `json:"formName"`   // 表单名称
	Data       json.RawMessage `json:"data"`       // 表单数据
	Submitter  json.RawMessage `json:"submitter"`  // 提交者信息
	IP         string          `json:"ip"`         // IP地址
	Status     string          `json:"status"`     // 状态
	CreateBy   string          `json:"createBy"`   // 创建人
	CreateTime time.Time       `json:"createTime"` // 创建时间
}
