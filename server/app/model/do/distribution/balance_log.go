package distribution

import "time"

// DistributionBalanceLog 分销资金变动记录
type DistributionBalanceLog struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	UserID        uint      `json:"user_id"`
	Type          int       `json:"type"` // 1-收入 2-支出
	Amount        float64   `json:"amount"`
	BeforeBalance float64   `json:"before_balance"`
	AfterBalance  float64   `json:"after_balance"`
	SourceType    int       `json:"source_type"` // 1-分销收益 2-提现
	SourceID      *uint     `json:"source_id"`
	Remark        string    `json:"remark"`
	CreateTime    time.Time `gorm:"autoCreateTime" json:"create_time"`
}

func (DistributionBalanceLog) TableName() string {
	return "share_balance_log"
}
