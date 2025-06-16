package shop

import (
	"time"
)

// ShopGoodsComboDO 套餐主表
type ShopGoodsComboDO struct {
	Id            int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name          string    `gorm:"column:name;not null" json:"name"`
	Description   string    `gorm:"column:description" json:"description"`
	Image         string    `gorm:"column:image" json:"image"`
	CategoryId    int64     `gorm:"column:category_id;not null" json:"categoryId"`
	BasePrice     float64   `gorm:"column:base_price;not null;default:0.00" json:"basePrice"`
	OriginalPrice float64   `gorm:"column:original_price;default:0.00" json:"originalPrice"`
	Sort          int       `gorm:"column:sort;default:0" json:"sort"`
	Status        int8      `gorm:"column:status;default:1" json:"status"`
	IsDelete      int8      `gorm:"column:is_delete;default:0" json:"isDelete"`
	CreateTime    time.Time `gorm:"column:create_time;autoCreateTime" json:"createTime"`
	UpdateTime    time.Time `gorm:"column:update_time;autoUpdateTime" json:"updateTime"`
}

// TableName 表名
func (ShopGoodsComboDO) TableName() string {
	return "shop_goods_combo"
}
