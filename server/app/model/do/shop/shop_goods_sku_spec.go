package shop

import (
	"time"
	"wnadmin/app/provider/db"

	"gorm.io/gorm"
)

// ShopGoodsSkuSpec SKU规格关联表
type ShopGoodsSkuSpec struct {
	Id          int64     `json:"id" gorm:"primaryKey;autoIncrement"`                  // 关联ID
	SkuId       int64     `json:"skuId" gorm:"column:sku_id;not null"`                 // SKU ID
	SpecGroupId int64     `json:"specGroupId" gorm:"column:spec_group_id;not null"`    // 规格组ID
	SpecValueId int64     `json:"specValueId" gorm:"column:spec_value_id;not null"`    // 规格值ID
	CreateTime  time.Time `json:"createTime" gorm:"column:create_time;autoCreateTime"` // 创建时间
}

// TableName 表名
func (s *ShopGoodsSkuSpec) TableName() string {
	return "shop_goods_sku_spec"
}

// Query 查询构造器
func (s *ShopGoodsSkuSpec) Query() *gorm.DB {
	return db.WnAdmin().Model(s)
}

// Create 创建
func (s *ShopGoodsSkuSpec) Create() error {
	return s.Query().Create(s).Error
}

// Delete 删除
func (s *ShopGoodsSkuSpec) Delete() error {
	return s.Query().Where("id = ?", s.Id).Delete(s).Error
}

// DeleteBySkuId 根据SKU ID删除
func (s *ShopGoodsSkuSpec) DeleteBySkuId() error {
	return s.Query().Where("sku_id = ?", s.SkuId).Delete(s).Error
}

// Update 更新
func (s *ShopGoodsSkuSpec) Update() error {
	return s.Query().Where("id = ?", s.Id).Updates(s).Error
}

// FindBySkuId 根据SKU ID查询规格关联
func (s *ShopGoodsSkuSpec) FindBySkuId(skuId int64) ([]*ShopGoodsSkuSpec, error) {
	var specs []*ShopGoodsSkuSpec
	err := s.Query().Where("sku_id = ?", skuId).Find(&specs).Error
	return specs, err
}

// BatchCreate 批量创建
func (s *ShopGoodsSkuSpec) BatchCreate(specs []*ShopGoodsSkuSpec) error {
	return s.Query().Create(&specs).Error
}
