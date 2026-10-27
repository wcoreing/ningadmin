package distribution

import "time"

// DistributionIncome 分销收益记录
type DistributionIncome struct {
	ID               uint       `gorm:"primaryKey" json:"id"`
	UserID           uint       `gorm:"column:user_id" json:"user_id"`
	UserName         string     `gorm:"column:user_name" json:"user_name"`       // 用户名称(关联字段)
	InviteeID        *uint      `gorm:"column:invitee_id" json:"invitee_id"`     // 被邀请人ID
	InviteeName      string     `gorm:"column:invitee_name" json:"invitee_name"` // 被邀请人姓名
	CommissionAmount float64    `gorm:"column:commission_amount" json:"commission_amount"`
	CommissionType   string     `gorm:"column:commission_type" json:"commission_type"`     // 佣金类型
	CommissionStatus int        `gorm:"column:commission_status" json:"commission_status"` // 0-待结算 1-已结算 2-已取消
	Status           int        `gorm:"column:commission_status" json:"status"`            // 状态字段别名，兼容前端
	SourceOrderNo    string     `gorm:"column:source_order_no" json:"source_order_no"`     // 来源订单号
	OrderNo          string     `gorm:"column:order_no" json:"order_no"`                   // 订单号(关联字段)
	OrderID          uint       `gorm:"column:order_id" json:"order_id"`                   // 订单ID(兼容字段)
	OrderAmount      float64    `gorm:"column:order_amount" json:"order_amount"`
	CommissionRate   float64    `gorm:"column:commission_rate" json:"commission_rate"`
	Level            int        `gorm:"column:level" json:"level"`
	UnfreezeTime     *time.Time `gorm:"column:unfreeze_time" json:"unfreeze_time"`
	ProcessTime      *time.Time `gorm:"column:process_time" json:"process_time"` // 处理时间
	Remark           string     `gorm:"column:remark" json:"remark"`             // 备注
	CreateTime       time.Time  `gorm:"column:create_time;autoCreateTime" json:"create_time"`
	UpdateTime       time.Time  `gorm:"column:update_time;autoUpdateTime" json:"update_time"`
}

func (DistributionIncome) TableName() string {
	return "share_income"
}
