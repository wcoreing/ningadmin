package shop

import (
	"encoding/json"
	"time"
	"wnadmin/app/provider/db"

	"gorm.io/gorm"
)

// 规格类型常量
const (
	SpecGroupTypeSingle = 1 // 单规格
	SpecGroupTypeMulti  = 2 // 多规格
	SpecGroupTypeCombo  = 3 // 套餐规格
)

// ComboGroupConfig 套餐组配置
type ComboGroupConfig struct {
	Description string `json:"description"`  // 组描述
	MinSelect   int    `json:"min_select"`   // 最少选择数量
	MaxSelect   int    `json:"max_select"`   // 最多选择数量
	IsRequired  bool   `json:"is_required"`  // 是否必选
	DisplayType string `json:"display_type"` // 显示类型：radio, checkbox, dropdown
}

// ShopGoodsSpecGroup 商品规格组表
type ShopGoodsSpecGroup struct {
	Id         int64                 `json:"id" gorm:"primaryKey;autoIncrement"`                  // 规格组ID
	GoodsId    int64                 `json:"goodsId" gorm:"column:goods_id;not null"`             // 商品ID
	Name       string                `json:"name" gorm:"column:name;size:100;not null"`           // 规格组名称，如：颜色、尺寸
	SpecType   int8                  `json:"specType" gorm:"column:spec_type;default:1"`          // 规格类型：1-单规格 2-多规格 3-套餐规格
	Config     *string               `json:"config" gorm:"column:config;type:json"`               // 规格组配置（JSON格式）
	Sort       int                   `json:"sort" gorm:"column:sort;default:0"`                   // 排序
	Status     int8                  `json:"status" gorm:"column:status;default:1"`               // 状态(0:禁用 1:启用)
	CreateTime time.Time             `json:"createTime" gorm:"column:create_time;autoCreateTime"` // 创建时间
	UpdateTime time.Time             `json:"updateTime" gorm:"column:update_time;autoUpdateTime"` // 更新时间
	Values     []*ShopGoodsSpecValue `json:"values" gorm:"-"`                                     // 规格值列表
}

// TableName 表名
func (ShopGoodsSpecGroup) TableName() string {
	return "shop_goods_spec_group"
}

// IsComboGroup 判断是否为套餐组
func (s *ShopGoodsSpecGroup) IsComboGroup() bool {
	return s.SpecType == SpecGroupTypeCombo
}

// IsSingleSpec 判断是否为单规格
func (s *ShopGoodsSpecGroup) IsSingleSpec() bool {
	return s.SpecType == SpecGroupTypeSingle
}

// IsMultiSpec 判断是否为多规格
func (s *ShopGoodsSpecGroup) IsMultiSpec() bool {
	return s.SpecType == SpecGroupTypeMulti
}

// GetComboConfig 获取套餐组配置
func (s *ShopGoodsSpecGroup) GetComboConfig() (*ComboGroupConfig, error) {
	if s.Config == nil || *s.Config == "" {
		return &ComboGroupConfig{}, nil
	}

	var config ComboGroupConfig
	if err := json.Unmarshal([]byte(*s.Config), &config); err != nil {
		return nil, err
	}
	return &config, nil
}

// SetComboConfig 设置套餐组配置
func (s *ShopGoodsSpecGroup) SetComboConfig(config *ComboGroupConfig) error {
	if config == nil {
		s.Config = nil
		return nil
	}

	configBytes, err := json.Marshal(config)
	if err != nil {
		return err
	}
	configStr := string(configBytes)
	s.Config = &configStr
	return nil
}

// Query 查询构建器
func (s *ShopGoodsSpecGroup) Query() *gorm.DB {
	return db.WnAdmin().Model(s)
}

// GetByGoodsIdAndSpecType 根据商品ID和规格类型获取规格组
func (s *ShopGoodsSpecGroup) GetByGoodsIdAndSpecType(goodsId int64, specType int8) ([]*ShopGoodsSpecGroup, error) {
	var specGroups []*ShopGoodsSpecGroup
	err := s.Query().Where("goods_id = ? AND spec_type = ? AND status = 1", goodsId, specType).
		Order("sort ASC").Find(&specGroups).Error
	return specGroups, err
}

// GetAllByGoodsId 根据商品ID获取所有规格类型的规格组
func (s *ShopGoodsSpecGroup) GetAllByGoodsId(goodsId int64) (map[int8][]*ShopGoodsSpecGroup, error) {
	var specGroups []*ShopGoodsSpecGroup
	err := s.Query().Where("goods_id = ? AND status = 1", goodsId).
		Order("spec_type ASC, sort ASC").Find(&specGroups).Error
	if err != nil {
		return nil, err
	}

	// 按规格类型分组
	result := make(map[int8][]*ShopGoodsSpecGroup)
	for _, group := range specGroups {
		result[group.SpecType] = append(result[group.SpecType], group)
	}
	return result, nil
}

// DeleteByGoodsIdAndSpecType 删除指定商品和规格类型的规格组
func (s *ShopGoodsSpecGroup) DeleteByGoodsIdAndSpecType(goodsId int64, specType int8) error {
	return s.Query().Where("goods_id = ? AND spec_type = ?", goodsId, specType).Delete(s).Error
}

// CreateBatch 批量创建规格组
func (s *ShopGoodsSpecGroup) CreateBatch(specGroups []*ShopGoodsSpecGroup) error {
	if len(specGroups) == 0 {
		return nil
	}
	return s.Query().Create(&specGroups).Error
}

// Create 创建
func (s *ShopGoodsSpecGroup) Create() error {
	return s.Query().Create(s).Error
}

// Delete 删除
func (s *ShopGoodsSpecGroup) Delete() error {
	return s.Query().Where("id = ?", s.Id).Delete(s).Error
}

// DeleteByGoodsId 根据商品ID删除
func (s *ShopGoodsSpecGroup) DeleteByGoodsId() error {
	return s.Query().Where("goods_id = ?", s.GoodsId).Delete(s).Error
}

// Update 更新
func (s *ShopGoodsSpecGroup) Update() error {
	return s.Query().Where("id = ?", s.Id).Updates(s).Error
}

// FindByGoodsId 根据商品ID查询规格组
func (s *ShopGoodsSpecGroup) FindByGoodsId(goodsId int64) ([]*ShopGoodsSpecGroup, error) {
	var groups []*ShopGoodsSpecGroup
	err := s.Query().Where("goods_id = ? AND status = 1", goodsId).Order("sort ASC").Find(&groups).Error
	return groups, err
}

// FindComboGroupsByGoodsId 根据商品ID查询套餐组
func (s *ShopGoodsSpecGroup) FindComboGroupsByGoodsId(goodsId int64) ([]*ShopGoodsSpecGroup, error) {
	return s.GetByGoodsIdAndSpecType(goodsId, SpecGroupTypeCombo)
}

// FindSingleSpecGroupsByGoodsId 根据商品ID查询单规格组
func (s *ShopGoodsSpecGroup) FindSingleSpecGroupsByGoodsId(goodsId int64) ([]*ShopGoodsSpecGroup, error) {
	return s.GetByGoodsIdAndSpecType(goodsId, SpecGroupTypeSingle)
}

// FindMultiSpecGroupsByGoodsId 根据商品ID查询多规格组
func (s *ShopGoodsSpecGroup) FindMultiSpecGroupsByGoodsId(goodsId int64) ([]*ShopGoodsSpecGroup, error) {
	return s.GetByGoodsIdAndSpecType(goodsId, SpecGroupTypeMulti)
}
