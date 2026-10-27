package share

import "time"

// ShareIncome 分享收益记录
type ShareIncome struct {
	ID               uint       `gorm:"primaryKey" json:"id"`
	UserID           uint       `json:"user_id"`           // 用户ID
	InviteeID        *uint      `json:"invitee_id"`        // 被邀请人ID
	InviteeName      string     `json:"invitee_name"`      // 被邀请人姓名
	CommissionAmount float64    `json:"commission_amount"` // 佣金金额
	CommissionType   string     `json:"commission_type"`   // 佣金类型
	CommissionStatus int        `json:"commission_status"` // 佣金状态(0-待结算,1-已结算,2-已取消)
	SourceOrderNo    string     `json:"source_order_no"`   // 来源订单号
	OrderAmount      float64    `json:"order_amount"`      // 订单金额
	CommissionRate   float64    `json:"commission_rate"`   // 佣金比例
	Level            int        `json:"level"`             // 分销层级
	UnfreezeTime     *time.Time `json:"unfreeze_time"`     // 解冻时间
	ProcessTime      *time.Time `json:"process_time"`      // 处理时间
	Remark           string     `json:"remark"`            // 备注
	CreateTime       time.Time  `gorm:"autoCreateTime" json:"create_time"`
	UpdateTime       time.Time  `gorm:"autoUpdateTime" json:"update_time"`

	// 关联字段
	UserName        string `gorm:"column:user_name" json:"user_name"`           // 推荐官名称(关联字段)
	InviteeNickname string `gorm:"column:invitee_name" json:"invitee_nickname"` // 被邀请人昵称(关联字段)
	OrderNo         string `gorm:"column:order_no" json:"order_no"`             // 订单号(关联字段)
}

func (ShareIncome) TableName() string {
	return "share_income"
}
