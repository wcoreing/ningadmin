package share

import (
	"time"
)

// ShareUser 分享用户模型
type ShareUser struct {
	ID               uint       `gorm:"primaryKey" json:"id"`
	UserID           uint       `json:"user_id"`
	RealName         string     `json:"real_name"`
	Phone            string     `json:"phone"`
	IDCard           string     `json:"id_card"`
	ParentID         *uint      `json:"parent_id"`
	Level            int        `json:"level"`
	Status           int        `json:"status"`
	TotalIncome      float64    `json:"total_income"`
	AvailableBalance float64    `json:"available_balance"`
	FrozenBalance    float64    `json:"frozen_balance"`
	TotalWithdraw    float64    `json:"total_withdraw"`
	TeamCount        int        `json:"team_count"`
	ApplyTime        *time.Time `json:"apply_time"`
	AuditTime        *time.Time `json:"audit_time"`
	AuditUser        string     `json:"audit_user"`
	RejectReason     string     `json:"reject_reason"`
	BindTime         *time.Time `json:"bind_time"`
	ExpireTime       *time.Time `json:"expire_time"`
	CreateTime       time.Time  `gorm:"autoCreateTime" json:"create_time"`
	UpdateTime       time.Time  `gorm:"autoUpdateTime" json:"update_time"`
	CreateBy         string     `json:"create_by"`
	UpdateBy         string     `json:"update_by"`
	Remark           string     `json:"remark"`
	Version          int        `json:"version"` // 乐观锁版本号

	// 关联字段(从user表读取)
	Username string `gorm:"column:username" json:"username"`
	Nickname string `gorm:"column:nickname" json:"nickname"`
	Avatar   string `gorm:"column:avatar" json:"avatar"`
}

func (ShareUser) TableName() string {
	return "share_user"
}
