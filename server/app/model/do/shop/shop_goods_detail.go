package shop

import (
	"time"
)

// ShopGoodsDetail 商品详情表
type ShopGoodsDetail struct {
	Id          int       `gorm:"primaryKey;column:id"`
	GoodsId     int       `gorm:"column:goods_id"`
	VideoUrl    string    `gorm:"column:video_url"`
	DetailDesc  string    `gorm:"column:detail_desc"`
	SpecDesc    string    `gorm:"column:spec_desc"`
	ServiceDesc string    `gorm:"column:service_desc"`
	CreateTime  time.Time `gorm:"column:create_time;autoCreateTime"`
	UpdateTime  time.Time `gorm:"column:update_time;autoUpdateTime"`
}

func (ShopGoodsDetail) TableName() string {
	return "shop_goods_detail"
}
