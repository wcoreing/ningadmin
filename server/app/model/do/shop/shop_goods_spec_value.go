package shop

import (
	"time"
	"wnadmin/app/provider/db"

	"gorm.io/gorm"
)

// ShopGoodsSpecValue 商品规格值表
type ShopGoodsSpecValue struct {
	Id        int64     `json:"id"`        // ID
	GoodsId   int64     `json:"goodsId"`   // 商品ID
	SpecId    int64     `json:"specId"`    // 规格ID
	Value     string    `json:"value"`     // 规格值
	Sort      int       `json:"sort"`      // 排序
	CreatedAt time.Time `json:"createdAt"` // 创建时间
	UpdatedAt time.Time `json:"updatedAt"` // 更新时间
}

// TableName 表名
func (s *ShopGoodsSpecValue) TableName() string {
	return "shop_goods_spec_value"
}

// Query 查询构造器
func (s *ShopGoodsSpecValue) Query() *gorm.DB {
	return db.DB.WnAdmin.Model(s)
}

// Create 创建
func (s *ShopGoodsSpecValue) Create() error {
	return s.Query().Create(s).Error
}

// Delete 删除
func (s *ShopGoodsSpecValue) Delete() error {
	return s.Query().Where("id = ?", s.Id).Delete(s).Error
}

// DeleteByGoodsId 根据商品ID删除
func (s *ShopGoodsSpecValue) DeleteByGoodsId() error {
	return s.Query().Where("goods_id = ?", s.GoodsId).Delete(s).Error
}

// Update 更新
func (s *ShopGoodsSpecValue) Update() error {
	return s.Query().Where("id = ?", s.Id).Updates(s).Error
}
