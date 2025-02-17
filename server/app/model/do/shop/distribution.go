package shop

import "time"

// DistributionConfig 分销配置
type DistributionConfig struct {
	ID           int64     `gorm:"column:id;primaryKey"`
	Enabled      int8      `gorm:"column:enabled"`
	Level        int8      `gorm:"column:level"`
	Mode         int8      `gorm:"column:mode"`
	ApplyEnabled int8      `gorm:"column:apply_enabled"`
	BindType     int8      `gorm:"column:bind_type"`
	BindMode     int8      `gorm:"column:bind_mode"`
	PosterBg     string    `gorm:"column:poster_bg"`
	CenterType   int8      `gorm:"column:center_type"`
	CreateTime   time.Time `gorm:"column:create_time"`
	UpdateTime   time.Time `gorm:"column:update_time"`
}

// Distributor 分销员
type Distributor struct {
	ID          int64     `gorm:"column:id;primaryKey"`
	UserID      int64     `gorm:"column:user_id"`
	ParentID    int64     `gorm:"column:parent_id"`
	Level       int8      `gorm:"column:level"`
	Status      int8      `gorm:"column:status"`
	TotalIncome float64   `gorm:"column:total_income"`
	TotalOrders int       `gorm:"column:total_orders"`
	BindTime    time.Time `gorm:"column:bind_time"`
	ExpireTime  time.Time `gorm:"column:expire_time"`
	CreateTime  time.Time `gorm:"column:create_time"`
	UpdateTime  time.Time `gorm:"column:update_time"`
}
