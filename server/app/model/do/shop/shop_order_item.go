package shop

import (
	"time"
)

// ShopOrderGoods 订单商品表
type ShopOrderGoods struct {
	Id          int64     `json:"id" gorm:"primaryKey"`
	OrderId     int64     `json:"orderId" gorm:"not null"`
	GoodsId     int64     `json:"goodsId" gorm:"not null"`
	GoodsName   string    `json:"goodsName" gorm:"not null"`
	GoodsImage  string    `json:"goodsImage" gorm:"not null"`
	SkuCode     string    `json:"skuCode" gorm:"not null"`
	Specs       string    `json:"specs" gorm:"not null"`
	Price       float64   `json:"price" gorm:"not null"`
	Quantity    int       `json:"quantity" gorm:"not null"`
	TotalAmount float64   `json:"totalAmount" gorm:"not null"`
	IsComment   int       `json:"isComment" gorm:"not null;default:0"`
	IsDelete    int       `json:"isDelete" gorm:"not null;default:0"`
	CreateTime  time.Time `json:"createTime" gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdateTime  time.Time `json:"updateTime" gorm:"not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
}
