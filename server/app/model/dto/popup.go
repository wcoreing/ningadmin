package dto

import (
	"time"
	"wnadmin/app/model/dto/base"
)

// PopupQuery 弹窗查询参数
type PopupQuery struct {
	base.PageQuery
	Title       string `form:"title" json:"title"`               // 标题
	TriggerType string `form:"trigger_type" json:"trigger_type"` // 触发类型
	Status      *int   `form:"status" json:"status"`             // 状态
}

// CreatePopupReq 创建弹窗请求
type CreatePopupReq struct {
	Title        string                   `json:"title" binding:"max=100"`                                    // 弹窗标题
	Content      string                   `json:"content" binding:"required"`                                 // 弹窗内容
	TitleStyle   map[string]interface{}   `json:"title_style"`                                                // 标题样式
	Buttons      []map[string]interface{} `json:"buttons"`                                                    // 按钮配置
	TriggerType  string                   `json:"trigger_type" binding:"required,oneof=manual auto schedule"` // 触发类型
	TriggerRule  map[string]interface{}   `json:"trigger_rule"`                                               // 触发规则
	Priority     int                      `json:"priority"`                                                   // 优先级
	AutoClose    int                      `json:"auto_close"`                                                 // 自动关闭时间(毫秒)
	Delay        int                      `json:"delay"`                                                      // 延迟显示时间(毫秒)
	MaskClosable int                      `json:"mask_closable"`                                              // 点击遮罩是否关闭
	StartTime    *time.Time               `json:"start_time"`                                                 // 开始时间
	EndTime      *time.Time               `json:"end_time"`                                                   // 结束时间
	Status       int                      `json:"status" binding:"oneof=0 1"`                                 // 状态
	Remark       string                   `json:"remark"`                                                     // 备注
}

// UpdatePopupReq 更新弹窗请求
type UpdatePopupReq struct {
	ID           int64                    `json:"id" binding:"required"`                                      // 弹窗ID
	Title        string                   `json:"title" binding:"max=100"`                                    // 弹窗标题
	Content      string                   `json:"content" binding:"required"`                                 // 弹窗内容
	TitleStyle   map[string]interface{}   `json:"title_style"`                                                // 标题样式
	Buttons      []map[string]interface{} `json:"buttons"`                                                    // 按钮配置
	TriggerType  string                   `json:"trigger_type" binding:"required,oneof=manual auto schedule"` // 触发类型
	TriggerRule  map[string]interface{}   `json:"trigger_rule"`                                               // 触发规则
	Priority     int                      `json:"priority"`                                                   // 优先级
	AutoClose    int                      `json:"auto_close"`                                                 // 自动关闭时间(毫秒)
	Delay        int                      `json:"delay"`                                                      // 延迟显示时间(毫秒)
	MaskClosable int                      `json:"mask_closable"`                                              // 点击遮罩是否关闭
	StartTime    *time.Time               `json:"start_time"`                                                 // 开始时间
	EndTime      *time.Time               `json:"end_time"`                                                   // 结束时间
	Status       int                      `json:"status" binding:"oneof=0 1"`                                 // 状态
	Remark       string                   `json:"remark"`                                                     // 备注
}

// PopupVO 弹窗视图对象
type PopupVO struct {
	ID           int64                    `json:"id"`
	Title        string                   `json:"title"`
	Content      string                   `json:"content"`
	TitleStyle   map[string]interface{}   `json:"title_style"`
	Buttons      []map[string]interface{} `json:"buttons"`
	TriggerType  string                   `json:"trigger_type"`
	TriggerRule  map[string]interface{}   `json:"trigger_rule"`
	Priority     int                      `json:"priority"`
	AutoClose    int                      `json:"auto_close"`
	Delay        int                      `json:"delay"`
	MaskClosable int                      `json:"mask_closable"`
	StartTime    *time.Time               `json:"start_time"`
	EndTime      *time.Time               `json:"end_time"`
	Status       int                      `json:"status"`
	ShowCount    int64                    `json:"show_count"`
	ClickCount   int64                    `json:"click_count"`
	CreateTime   time.Time                `json:"create_time"`
	UpdateTime   time.Time                `json:"update_time"`
}

// PopupStatusReq 弹窗状态修改请求
type PopupStatusReq struct {
	ID     int64 `json:"id" binding:"required"`      // 弹窗ID
	Status int   `json:"status" binding:"oneof=0 1"` // 状态
}

// PopupStatsReq 弹窗统计更新请求
type PopupStatsReq struct {
	ID   int64  `json:"id" binding:"required"`                    // 弹窗ID
	Type string `json:"type" binding:"required,oneof=show click"` // 统计类型：show=显示次数 click=点击次数
}

// GetActivePopupsReq 获取活跃弹窗请求
type GetActivePopupsReq struct {
	TriggerType string `form:"trigger_type" json:"trigger_type"` // 触发类型
	UserID      int64  `form:"user_id" json:"user_id"`           // 用户ID（用于个性化规则）
}
