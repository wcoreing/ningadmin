package shop

import (
	"time"
	"wnadmin/app/provider/db"

	"gorm.io/gorm"
)

// ShopGoodsSpec 商品规格表
type ShopGoodsSpec struct {
	Id         int64                 `json:"id"`              // 规格ID
	GoodsId    int64                 `json:"goodsId"`         // 商品ID
	Name       string                `json:"name"`            // 规格名称
	Sort       int                   `json:"sort"`            // 排序
	IsDelete   int8                  `json:"isDelete"`        // 是否删除
	CreateTime time.Time             `json:"createTime"`      // 创建时间
	UpdateTime time.Time             `json:"updateTime"`      // 更新时间
	Values     []*ShopGoodsSpecValue `json:"values" gorm:"-"` // 规格值列表
}

// TableName 表名
func (s *ShopGoodsSpec) TableName() string {
	return "shop_goods_spec"
}

// Query 查询构造器
func (s *ShopGoodsSpec) Query() *gorm.DB {
	return db.WnAdmin().Model(s)
}

// Create 创建
func (s *ShopGoodsSpec) Create() error {
	return s.Query().Create(s).Error
}

// Delete 删除
func (s *ShopGoodsSpec) Delete() error {
	return s.Query().Where("id = ?", s.Id).Delete(s).Error
}

// DeleteByGoodsId 根据商品ID删除
func (s *ShopGoodsSpec) DeleteByGoodsId() error {
	return s.Query().Where("goods_id = ?", s.GoodsId).Delete(s).Error
}

// Update 更新
func (s *ShopGoodsSpec) Update() error {
	return s.Query().Where("id = ?", s.Id).Updates(s).Error
}
