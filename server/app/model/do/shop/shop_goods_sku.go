package shop

import (
	"time"
	"wnadmin/app/provider/db"

	"gorm.io/gorm"
)

// ShopGoodsSku 商品SKU表
type ShopGoodsSku struct {
	Id            int64     `json:"id" gorm:"primaryKey;autoIncrement"`                                         // SKU ID
	GoodsId       int64     `json:"goodsId" gorm:"column:goods_id;not null"`                                    // 商品ID
	SpecType      int8      `json:"specType" gorm:"column:spec_type;default:1"`                                 // 规格类型：1-单规格 2-多规格 3-套餐规格
	SkuCode       string    `json:"skuCode" gorm:"column:sku_code;size:100;not null"`                           // SKU编码
	Image         string    `json:"image" gorm:"column:image;size:500"`                                         // SKU图片URL
	Price         float64   `json:"price" gorm:"column:price;type:decimal(10,2);default:0.00"`                  // 售价
	OriginalPrice float64   `json:"originalPrice" gorm:"column:original_price;type:decimal(10,2);default:0.00"` // 原价
	CostPrice     float64   `json:"costPrice" gorm:"column:cost_price;type:decimal(10,2);default:0.00"`         // 成本价
	Stock         int       `json:"stock" gorm:"column:stock;default:0"`                                        // 库存
	Weight        float64   `json:"weight" gorm:"column:weight;type:decimal(8,3)"`                              // 重量(kg)
	Volume        float64   `json:"volume" gorm:"column:volume;type:decimal(8,3)"`                              // 体积(m³)
	Sort          int       `json:"sort" gorm:"column:sort;default:0"`                                          // 排序
	Status        int8      `json:"status" gorm:"column:status;default:1"`                                      // 状态(0:禁用 1:启用)
	CreateTime    time.Time `json:"createTime" gorm:"column:create_time;autoCreateTime"`                        // 创建时间
	UpdateTime    time.Time `json:"updateTime" gorm:"column:update_time;autoUpdateTime"`                        // 更新时间

	// 关联字段，不映射到数据库
	SpecRelations []*ShopGoodsSkuSpec `json:"specRelations" gorm:"-"` // SKU规格关联列表
}

// TableName 表名
func (s *ShopGoodsSku) TableName() string {
	return "shop_goods_sku"
}

// Query 查询构造器
func (s *ShopGoodsSku) Query() *gorm.DB {
	return db.WnAdmin().Model(s)
}

// Create 创建
func (s *ShopGoodsSku) Create() error {
	return s.Query().Create(s).Error
}

// Delete 删除
func (s *ShopGoodsSku) Delete() error {
	return s.Query().Where("id = ?", s.Id).Delete(s).Error
}

// DeleteByGoodsId 根据商品ID删除
func (s *ShopGoodsSku) DeleteByGoodsId() error {
	return s.Query().Where("goods_id = ?", s.GoodsId).Delete(s).Error
}

// Update 更新
func (s *ShopGoodsSku) Update() error {
	return s.Query().Where("id = ?", s.Id).Updates(s).Error
}

// FindByGoodsId 根据商品ID查询SKU
func (s *ShopGoodsSku) FindByGoodsId(goodsId int64) ([]*ShopGoodsSku, error) {
	var skus []*ShopGoodsSku
	err := s.Query().Where("goods_id = ? AND status = 1", goodsId).Order("sort ASC").Find(&skus).Error
	return skus, err
}

// GetByGoodsIdAndSpecType 根据商品ID和规格类型获取SKU
func (s *ShopGoodsSku) GetByGoodsIdAndSpecType(goodsId int64, specType int8) ([]*ShopGoodsSku, error) {
	var skus []*ShopGoodsSku
	err := s.Query().Where("goods_id = ? AND spec_type = ? AND status = 1", goodsId, specType).
		Order("sort ASC").Find(&skus).Error
	return skus, err
}

// GetAllByGoodsId 根据商品ID获取所有规格类型的SKU
func (s *ShopGoodsSku) GetAllByGoodsId(goodsId int64) (map[int8][]*ShopGoodsSku, error) {
	var skus []*ShopGoodsSku
	err := s.Query().Where("goods_id = ? AND status = 1", goodsId).
		Order("spec_type ASC, sort ASC").Find(&skus).Error
	if err != nil {
		return nil, err
	}

	// 按规格类型分组
	result := make(map[int8][]*ShopGoodsSku)
	for _, sku := range skus {
		result[sku.SpecType] = append(result[sku.SpecType], sku)
	}
	return result, nil
}

// DeleteByGoodsIdAndSpecType 删除指定商品和规格类型的SKU
func (s *ShopGoodsSku) DeleteByGoodsIdAndSpecType(goodsId int64, specType int8) error {
	return s.Query().Where("goods_id = ? AND spec_type = ?", goodsId, specType).Delete(s).Error
}

// CreateBatch 批量创建SKU
func (s *ShopGoodsSku) CreateBatch(skus []*ShopGoodsSku) error {
	if len(skus) == 0 {
		return nil
	}
	return s.Query().Create(&skus).Error
}
