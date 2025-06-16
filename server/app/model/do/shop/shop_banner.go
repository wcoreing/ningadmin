package shop

import (
	"time"
	"wnadmin/app/provider/db"

	"gorm.io/gorm"
)

type ShopBanner struct {
	Id         int64      `json:"id" gorm:"primaryKey"`
	Title      string     `json:"title" gorm:"not null"`
	Image      string     `json:"image" gorm:"not null"`
	Url        string     `json:"url"`
	Position   int        `json:"position" gorm:"not null;default:1"`
	Type       int        `json:"type" gorm:"not null;default:1"`
	TargetId   int64      `json:"targetId"`
	Sort       int        `json:"sort" gorm:"not null;default:0"`
	Status     int        `json:"status" gorm:"not null;default:1"`
	IsDelete   int        `json:"isDelete" gorm:"not null;default:0"`
	CreateTime time.Time  `json:"createTime" gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdateTime time.Time  `json:"updateTime" gorm:"not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
	Goods      *ShopGoods `gorm:"foreignKey:Id;references:TargetId"`
}

// 定义常量
const (
	// 位置
	BannerPositionHome     = 1 // 首页
	BannerPositionCategory = 2 // 分类页
	BannerPositionActivity = 3 // 活动页

	// 类型
	BannerTypeGoods    = 1 // 商品
	BannerTypeActivity = 2 // 活动
	BannerTypeUrl      = 3 // 外链

	// 状态
	BannerStatusDisable = 0 // 禁用
	BannerStatusEnable  = 1 // 启用
)

// Query 返回查询构建器
func (m *ShopBanner) Query() *gorm.DB {
	return db.WnAdmin().Model(m).Where("is_delete = ?", 0)
}

// Create 创建记录
func (m *ShopBanner) Create() error {
	return db.WnAdmin().Create(m).Error
}

// Update 更新记录
func (m *ShopBanner) Update() error {
	return db.WnAdmin().Save(m).Error
}

// First 获取单条记录
func (m *ShopBanner) First() error {
	return m.Query().First(m).Error
}
