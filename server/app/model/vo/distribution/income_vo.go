package distribution

import "time"

// IncomeVO 分销收益列表项
type IncomeVO struct {
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
	CreateTime       string    `json:"createTime"`     // 创建时间，改为 string 类型
	OrderNo          string    `json:"orderNo"`        // 订单号
}
