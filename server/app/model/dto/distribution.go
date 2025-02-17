package dto

import "time"

// 分销配置DTO
type DistributionConfigDTO struct {
	ID           int64  `json:"id"`
	Enabled      bool   `json:"enabled"`      // 是否启用
	Level        int8   `json:"level"`        // 分销层级
	Mode         int8   `json:"mode"`         // 分销模式
	ApplyEnabled bool   `json:"applyEnabled"` // 是否开启申请
	BindType     int8   `json:"bindType"`     // 绑定方式
	BindMode     int8   `json:"bindMode"`     // 绑定模式
	PosterBg     string `json:"posterBg"`     // 海报背景
	CenterType   int8   `json:"centerType"`   // 分销中心类型
}

// 分销员DTO
type DistributorDTO struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"userId"`
	ParentID    int64     `json:"parentId"`
	Level       int8      `json:"level"`
	Status      int8      `json:"status"`
	TotalIncome float64   `json:"totalIncome"`
	TotalOrders int       `json:"totalOrders"`
	BindTime    time.Time `json:"bindTime"`
	ExpireTime  time.Time `json:"expireTime"`
}

// DistributorQueryDTO 分销员查询参数
type DistributorQueryDTO struct {
	PageNum  int  `form:"pageNum"`
	PageSize int  `form:"pageSize"`
	Status   int8 `form:"status"`
}

// DistributorAuditDTO 分销员审核参数
type DistributorAuditDTO struct {
	ID     int64 `json:"id" binding:"required"`
	Status int8  `json:"status" binding:"required,oneof=1 2"`
}
