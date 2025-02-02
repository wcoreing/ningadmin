package shop

import (
	"time"
	"wnadmin/app/provider/db"

	"gorm.io/gorm"
)

// ShopGoodsSku 商品SKU表
type ShopGoodsSku struct {
	Id            int64     `json:"id"`            // SKU ID
	GoodsId       int64     `json:"goodsId"`       // 商品ID
	SkuCode       string    `json:"skuCode"`       // sku编码
	SpecsJson     string    `json:"specsJson"`     // 规格值字符串
	Image         string    `json:"image"`         // SKU图片
	Price         float64   `json:"price"`         // 售价
	OriginalPrice float64   `json:"originalPrice"` // 原价
	CostPrice     float64   `json:"costPrice"`     // 成本价
	Stock         int       `json:"stock"`         // 库存
	Code          string    `json:"code"`          // 商品编码
	Weight        float64   `json:"weight"`        // 重量(kg)
	Volume        float64   `json:"volume"`        // 体积(m³)
	Sort          int       `json:"sort"`          // 排序
	Status        int8      `json:"status"`        // 状态(0:禁用 1:启用)
	CreateTime    time.Time `json:"createTime"`    // 创建时间
	UpdateTime    time.Time `json:"updateTime"`    // 更新时间
}

// TableName 表名
func (s *ShopGoodsSku) TableName() string {
	return "shop_goods_sku"
}

// Query 查询构造器
func (s *ShopGoodsSku) Query() *gorm.DB {
	return db.DB.WnAdmin.Model(s)
}
