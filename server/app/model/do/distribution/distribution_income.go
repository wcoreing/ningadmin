package distribution

import "time"

// DistributionIncome 分销收益记录
type DistributionIncome struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	UserID           uint      `json:"user_id"`
	UserName         string    `gorm:"column:user_name" json:"user_name"` // 用户名称(关联字段)
	OrderID          uint      `json:"order_id"`
	OrderNo          string    `gorm:"column:order_no" json:"order_no"` // 订单号(关联字段)
	OrderAmount      float64   `json:"order_amount"`
	CommissionRate   float64   `json:"commission_rate"`
	CommissionAmount float64   `json:"commission_amount"`
	Level            int       `json:"level"`
	Status           int       `json:"status"` // 0-冻结中 1-已解冻 2-已失效 3-已入账
	UnfreezeTime     time.Time `json:"unfreeze_time"`
	ProcessTime      time.Time `json:"process_time"`   // 处理时间
	ProcessRemark    string    `json:"process_remark"` // 处理备注
	CreateTime       time.Time `gorm:"autoCreateTime" json:"create_time"`
	UpdateTime       time.Time `gorm:"autoUpdateTime" json:"update_time"`
}

func (DistributionIncome) TableName() string {
	return "distribution_income"
}
