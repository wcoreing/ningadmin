package distribution

import "time"

// DistributionWithdraw 分销提现模型
type DistributionWithdraw struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      uint      `json:"user_id"`
	Amount      float64   `json:"amount"`
	Method      string    `json:"method"`
	Account     string    `json:"account"`
	RealName    string    `json:"real_name"`
	Status      int       `json:"status"`
	AuditTime   time.Time `json:"audit_time"`
	AuditRemark string    `json:"audit_remark"`
	PayTime     time.Time `json:"pay_time"`
	PayRemark   string    `json:"pay_remark"`
	CreateTime  time.Time `gorm:"autoCreateTime" json:"create_time"`
	UpdateTime  time.Time `gorm:"autoUpdateTime" json:"update_time"`
}

func (DistributionWithdraw) TableName() string {
	return "share_withdraw"
}
