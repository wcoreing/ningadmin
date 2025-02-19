package distribution

import "time"

// DistributionRelation 分销用户关系
type DistributionRelation struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     uint      `json:"user_id"`
	AncestorID uint      `json:"ancestor_id"`
	Level      int       `json:"level"`
	Status     int       `json:"status"`
	EndTime    time.Time `json:"end_time"`
	BindType   string    `json:"bind_type"`
	CreateTime time.Time `gorm:"autoCreateTime" json:"create_time"`
}

func (DistributionRelation) TableName() string {
	return "distribution_relation"
}
