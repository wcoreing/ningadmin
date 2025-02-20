package dto

import "time"

// DistributionIncomeQueryDTO 分销收益查询参数
type DistributionIncomeQueryDTO struct {
	UserId    string    `json:"userId"`
	Nickname  string    `json:"nickname"`
	Status    string    `json:"status"`
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
	PageNum   int       `json:"pageNum"`
	PageSize  int       `json:"pageSize"`
}

// DistributionIncomeDTO 分销收益返回数据
type DistributionIncomeDTO struct {
	Id               uint      `json:"id"`
	UserId           uint      `json:"userId"`
	UserName         string    `json:"userName"` // 用户名称
	OrderId          uint      `json:"orderId"`
	OrderAmount      float64   `json:"orderAmount"`    // 订单金额
	CommissionRate   float64   `json:"commissionRate"` // 佣金比例
	CommissionAmount float64   `json:"amount"`         // 佣金金额
	Level            int       `json:"level"`          // 分销等级
	Status           int       `json:"status"`         // 状态
	StatusText       string    `json:"statusText"`     // 状态文本
	UnfreezeTime     time.Time `json:"unfreezeTime"`   // 解冻时间
	CreateTime       time.Time `json:"createTime"`     // 创建时间
}

// DistributionIncomePageReq 分销收益分页查询请求
type DistributionIncomeListReq struct {
	UserId    string `json:"userId" form:"userId"`       // 用户ID
	Nickname  string `json:"nickname" form:"nickname"`   // 用户昵称
	Status    string `json:"status" form:"status"`       // 收益状态
	StartTime string `json:"startTime" form:"startTime"` // 开始时间
	EndTime   string `json:"endTime" form:"endTime"`     // 结束时间
	PageNum   int    `json:"pageNum" form:"pageNum"`     // 页码
	PageSize  int    `json:"pageSize" form:"pageSize"`   // 每页大小
}

// DistributionIncomeUnfreezeReq 分销收益解冻请求
type DistributionIncomeUnfreezeReq struct {
	Id int64 `json:"id" form:"id"` // 收益ID
}
