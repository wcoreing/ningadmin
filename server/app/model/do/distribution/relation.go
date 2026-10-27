package distribution

import "time"

// DistributionRelation 分销用户关系（基于通用关系表）
type DistributionRelation struct {
	ID            uint       `gorm:"primaryKey" json:"id"`
	UserID        uint       `gorm:"column:user_id" json:"user_id"`
	RelatedUserID uint       `gorm:"column:related_user_id" json:"related_user_id"` // 对应原来的 ancestor_id
	RelationType  string     `gorm:"column:relation_type" json:"relation_type"`     // 固定为 "invite" 或 "share"
	Direction     string     `gorm:"column:direction" json:"direction"`             // 固定为 "one_way"
	Status        int        `gorm:"column:status" json:"status"`                   // 1-有效 0-无效 2-待确认
	Source        string     `gorm:"column:source" json:"source"`                   // 关系来源
	Metadata      string     `gorm:"column:metadata;type:json" json:"metadata"`     // 扩展数据，可存储 level 等信息
	CreateTime    time.Time  `gorm:"column:create_time" json:"create_time"`
	UpdateTime    time.Time  `gorm:"column:update_time" json:"update_time"`
	ExpireTime    *time.Time `gorm:"column:expire_time" json:"expire_time"`
}

func (DistributionRelation) TableName() string {
	return "user_relation"
}

// GetAncestorID 获取祖先ID（兼容旧字段名）
func (d *DistributionRelation) GetAncestorID() uint {
	return d.RelatedUserID
}

// SetAncestorID 设置祖先ID（兼容旧字段名）
func (d *DistributionRelation) SetAncestorID(ancestorID uint) {
	d.RelatedUserID = ancestorID
}

// GetLevel 从Metadata中获取层级信息
func (d *DistributionRelation) GetLevel() int {
	// 这里可以解析JSON metadata获取level信息
	// 简化实现，可以根据实际需要扩展
	return 1
}

// SetLevel 设置层级信息到Metadata
func (d *DistributionRelation) SetLevel(level int) {
	// 这里可以将level信息存储到JSON metadata中
	// 简化实现，可以根据实际需要扩展
}
