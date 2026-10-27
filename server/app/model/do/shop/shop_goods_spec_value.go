package shop

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
	"wnadmin/app/provider/db"

	"gorm.io/gorm"
)

// 规格值类型常量
const (
	SpecValueTypeSingle = 1 // 单规格值
	SpecValueTypeMulti  = 2 // 多规格值
	SpecValueTypeCombo  = 3 // 套餐选项
)

// ComboOptionConfig 套餐选项配置
type ComboOptionConfig struct {
	OptionGoodsId *int64  `json:"option_goods_id"` // 关联商品ID
	OptionSkuId   *int64  `json:"option_sku_id"`   // 关联SKU ID
	AddPrice      float64 `json:"add_price"`       // 加价金额
	Quantity      int     `json:"quantity"`        // 数量
	IsDefault     bool    `json:"is_default"`      // 是否默认选择
	Description   string  `json:"description"`     // 选项详细描述（可与name不同）
	Image         string  `json:"image"`           // 选项图片
}

// ShopGoodsSpecValue 商品规格值表
type ShopGoodsSpecValue struct {
	Id          int64     `json:"id" gorm:"primaryKey;autoIncrement"`                  // 规格值ID
	GoodsId     int64     `json:"goodsId" gorm:"column:goods_id;not null"`             // 商品ID
	SpecGroupId int64     `json:"specGroupId" gorm:"column:spec_group_id;not null"`    // 规格组ID
	Name        string    `json:"name" gorm:"column:name;size:100;not null"`           // 规格值显示名称（用户可见）
	Value       string    `json:"value" gorm:"column:value;size:100;not null"`         // 系统标识符（API/查询用）
	SpecType    int8      `json:"specType" gorm:"column:spec_type;default:1"`          // 规格类型：1-单规格值 2-多规格值 3-套餐选项
	Config      *string   `json:"config" gorm:"column:config;type:json"`               // 规格值配置（JSON格式）
	Sort        int       `json:"sort" gorm:"column:sort;default:0"`                   // 排序
	Status      int8      `json:"status" gorm:"column:status;default:1"`               // 状态(0:禁用 1:启用)
	CreateTime  time.Time `json:"createTime" gorm:"column:create_time;autoCreateTime"` // 创建时间
	UpdateTime  time.Time `json:"updateTime" gorm:"column:update_time;autoUpdateTime"` // 更新时间

	// 兼容旧字段（仅用于数据转换，不映射到数据库）
	SpecId    int64     `json:"specId" gorm:"-"`    // 兼容旧的规格ID字段
	CreatedAt time.Time `json:"createdAt" gorm:"-"` // 兼容旧的创建时间字段
	UpdatedAt time.Time `json:"updatedAt" gorm:"-"` // 兼容旧的更新时间字段
}

// TableName 表名
func (s *ShopGoodsSpecValue) TableName() string {
	return "shop_goods_spec_value"
}

// BeforeCreate GORM钩子：创建前自动生成value（如果为空）
func (s *ShopGoodsSpecValue) BeforeCreate(tx *gorm.DB) error {
	if s.Value == "" {
		s.Value = s.generateValue()
	}
	return nil
}

// generateValue 自动生成value值
func (s *ShopGoodsSpecValue) generateValue() string {
	switch s.SpecType {
	case SpecValueTypeSingle, SpecValueTypeMulti:
		// 单规格值和多规格值：使用name的简化版本
		return s.normalizeValueFromName(s.Name)
	case SpecValueTypeCombo:
		// 套餐选项：尝试从config中获取option_goods_id
		if s.Config != nil {
			var config ComboOptionConfig
			if err := json.Unmarshal([]byte(*s.Config), &config); err == nil {
				if config.OptionGoodsId != nil {
					return fmt.Sprintf("option_%d", *config.OptionGoodsId)
				}
			}
		}
		// 如果无法从config获取，使用name的简化版本
		return s.normalizeValueFromName(s.Name)
	default:
		return s.normalizeValueFromName(s.Name)
	}
}

// normalizeValueFromName 从name生成标准化的value
func (s *ShopGoodsSpecValue) normalizeValueFromName(name string) string {
	// 简单的标准化逻辑：
	// 1. 转小写
	// 2. 替换空格和特殊字符为下划线
	// 3. 移除连续的下划线
	value := strings.ToLower(name)
	value = strings.ReplaceAll(value, " ", "_")
	value = strings.ReplaceAll(value, "-", "_")
	value = strings.ReplaceAll(value, ".", "_")
	value = strings.ReplaceAll(value, ",", "_")
	value = strings.ReplaceAll(value, "，", "_")

	// 移除连续的下划线
	for strings.Contains(value, "__") {
		value = strings.ReplaceAll(value, "__", "_")
	}

	// 移除首尾下划线
	value = strings.Trim(value, "_")

	// 如果结果为空，使用时间戳
	if value == "" {
		value = fmt.Sprintf("value_%d", time.Now().UnixNano())
	}

	return value
}

// IsComboOption 判断是否为套餐选项
func (s *ShopGoodsSpecValue) IsComboOption() bool {
	return s.SpecType == SpecValueTypeCombo
}

// IsSingleSpec 判断是否为单规格值
func (s *ShopGoodsSpecValue) IsSingleSpec() bool {
	return s.SpecType == SpecValueTypeSingle
}

// IsMultiSpec 判断是否为多规格值
func (s *ShopGoodsSpecValue) IsMultiSpec() bool {
	return s.SpecType == SpecValueTypeMulti
}

// GetComboConfig 获取套餐选项配置
func (s *ShopGoodsSpecValue) GetComboConfig() (*ComboOptionConfig, error) {
	if !s.IsComboOption() || s.Config == nil {
		return nil, nil
	}

	var config ComboOptionConfig
	if err := json.Unmarshal([]byte(*s.Config), &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// SetComboConfig 设置套餐选项配置
func (s *ShopGoodsSpecValue) SetComboConfig(config *ComboOptionConfig) error {
	if config == nil {
		s.Config = nil
		s.SpecType = SpecValueTypeSingle
		return nil
	}

	configBytes, err := json.Marshal(config)
	if err != nil {
		return err
	}

	configStr := string(configBytes)
	s.Config = &configStr
	s.SpecType = SpecValueTypeCombo

	// 如果value为空且有option_goods_id，自动生成value
	if s.Value == "" && config.OptionGoodsId != nil {
		s.Value = fmt.Sprintf("option_%d", *config.OptionGoodsId)
	}

	return nil
}

// GetDisplayName 获取显示名称（优先使用name，套餐选项可回退到description）
func (s *ShopGoodsSpecValue) GetDisplayName() string {
	if s.Name != "" {
		return s.Name
	}

	// 对于套餐选项，如果name为空，尝试使用config中的description
	if s.IsComboOption() {
		if config, err := s.GetComboConfig(); err == nil && config != nil && config.Description != "" {
			return config.Description
		}
	}

	return s.Value // 最后回退到value
}

// GetDescription 获取详细描述
func (s *ShopGoodsSpecValue) GetDescription() string {
	if s.IsComboOption() {
		if config, err := s.GetComboConfig(); err == nil && config != nil {
			return config.Description
		}
	}
	return s.Name // 普通规格值返回name
}

// GetOptionGoodsId 获取选项商品ID
func (s *ShopGoodsSpecValue) GetOptionGoodsId() int64 {
	config, err := s.GetComboConfig()
	if err != nil || config == nil || config.OptionGoodsId == nil {
		return 0
	}
	return *config.OptionGoodsId
}

// GetOptionSkuId 获取选项SKU ID
func (s *ShopGoodsSpecValue) GetOptionSkuId() int64 {
	config, err := s.GetComboConfig()
	if err != nil || config == nil || config.OptionSkuId == nil {
		return 0
	}
	return *config.OptionSkuId
}

// GetAddPrice 获取加价金额
func (s *ShopGoodsSpecValue) GetAddPrice() float64 {
	config, err := s.GetComboConfig()
	if err != nil || config == nil {
		return 0
	}
	return config.AddPrice
}

// GetQuantity 获取数量
func (s *ShopGoodsSpecValue) GetQuantity() int {
	config, err := s.GetComboConfig()
	if err != nil || config == nil {
		return 1
	}
	return config.Quantity
}

// IsDefault 是否为默认选择
func (s *ShopGoodsSpecValue) IsDefault() bool {
	config, err := s.GetComboConfig()
	if err != nil || config == nil {
		return false
	}
	return config.IsDefault
}

// ValidateConsistency 验证数据一致性
func (s *ShopGoodsSpecValue) ValidateConsistency() error {
	switch s.SpecType {
	case SpecValueTypeSingle, SpecValueTypeMulti:
		// 单规格值和多规格值不应该有复杂配置
		if s.Config != nil {
			var config map[string]interface{}
			if err := json.Unmarshal([]byte(*s.Config), &config); err == nil && len(config) > 0 {
				return fmt.Errorf("单规格值和多规格值不应该有配置信息")
			}
		}
	case SpecValueTypeCombo:
		// 套餐选项必须有配置
		if s.Config == nil {
			return fmt.Errorf("套餐选项必须有配置信息")
		}

		config, err := s.GetComboConfig()
		if err != nil {
			return fmt.Errorf("套餐选项配置格式错误: %v", err)
		}

		if config.OptionGoodsId == nil {
			return fmt.Errorf("套餐选项必须指定关联商品ID")
		}

		// 检查value和option_goods_id的一致性（警告级别）
		expectedValue := fmt.Sprintf("option_%d", *config.OptionGoodsId)
		if s.Value != expectedValue {
			// 这里只是记录警告，不强制失败
			// 可以根据需要调整为严格模式
		}
	}

	return nil
}

// Query 查询构造器
func (s *ShopGoodsSpecValue) Query() *gorm.DB {
	return db.WnAdmin().Model(s)
}

// GetByGoodsIdAndSpecType 根据商品ID和规格类型获取规格值
func (s *ShopGoodsSpecValue) GetByGoodsIdAndSpecType(goodsId int64, specType int8) ([]*ShopGoodsSpecValue, error) {
	var specValues []*ShopGoodsSpecValue
	err := s.Query().Where("goods_id = ? AND spec_type = ? AND status = 1", goodsId, specType).
		Order("sort ASC").Find(&specValues).Error
	return specValues, err
}

// GetAllByGoodsId 根据商品ID获取所有规格类型的规格值
func (s *ShopGoodsSpecValue) GetAllByGoodsId(goodsId int64) (map[int8][]*ShopGoodsSpecValue, error) {
	var specValues []*ShopGoodsSpecValue
	err := s.Query().Where("goods_id = ? AND status = 1", goodsId).
		Order("spec_type ASC, sort ASC").Find(&specValues).Error
	if err != nil {
		return nil, err
	}

	// 按规格类型分组
	result := make(map[int8][]*ShopGoodsSpecValue)
	for _, value := range specValues {
		result[value.SpecType] = append(result[value.SpecType], value)
	}
	return result, nil
}

// DeleteByGoodsIdAndSpecType 删除指定商品和规格类型的规格值
func (s *ShopGoodsSpecValue) DeleteByGoodsIdAndSpecType(goodsId int64, specType int8) error {
	return s.Query().Where("goods_id = ? AND spec_type = ?", goodsId, specType).Delete(s).Error
}

// CreateBatch 批量创建规格值
func (s *ShopGoodsSpecValue) CreateBatch(specValues []*ShopGoodsSpecValue) error {
	if len(specValues) == 0 {
		return nil
	}
	return s.Query().Create(&specValues).Error
}

// Create 创建
func (s *ShopGoodsSpecValue) Create() error {
	// 创建前验证一致性
	if err := s.ValidateConsistency(); err != nil {
		return err
	}
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

// DeleteBySpecGroupId 根据规格组ID删除
func (s *ShopGoodsSpecValue) DeleteBySpecGroupId() error {
	return s.Query().Where("spec_group_id = ?", s.SpecGroupId).Delete(s).Error
}

// Update 更新
func (s *ShopGoodsSpecValue) Update() error {
	// 更新前验证一致性
	if err := s.ValidateConsistency(); err != nil {
		return err
	}
	return s.Query().Where("id = ?", s.Id).Updates(s).Error
}

// FindBySpecGroupId 根据规格组ID查询规格值
func (s *ShopGoodsSpecValue) FindBySpecGroupId(specGroupId int64) ([]*ShopGoodsSpecValue, error) {
	var values []*ShopGoodsSpecValue
	err := s.Query().Where("spec_group_id = ? AND status = 1", specGroupId).Order("sort ASC").Find(&values).Error
	return values, err
}

// FindByGoodsId 根据商品ID查询规格值
func (s *ShopGoodsSpecValue) FindByGoodsId(goodsId int64) ([]*ShopGoodsSpecValue, error) {
	var values []*ShopGoodsSpecValue
	err := s.Query().Where("goods_id = ? AND status = 1", goodsId).Order("spec_group_id ASC, sort ASC").Find(&values).Error
	return values, err
}

// FindByGoodsIdAndType 根据商品ID和规格类型查询规格值
func (s *ShopGoodsSpecValue) FindByGoodsIdAndType(goodsId int64, specType int8) ([]*ShopGoodsSpecValue, error) {
	var values []*ShopGoodsSpecValue
	err := s.Query().Where("goods_id = ? AND spec_type = ? AND status = 1", goodsId, specType).Order("spec_group_id ASC, sort ASC").Find(&values).Error
	return values, err
}

// FindComboOptionsByGoodsId 根据商品ID查询套餐选项
func (s *ShopGoodsSpecValue) FindComboOptionsByGoodsId(goodsId int64) ([]*ShopGoodsSpecValue, error) {
	return s.FindByGoodsIdAndType(goodsId, SpecValueTypeCombo)
}

// FindSingleSpecValuesByGoodsId 根据商品ID查询单规格值
func (s *ShopGoodsSpecValue) FindSingleSpecValuesByGoodsId(goodsId int64) ([]*ShopGoodsSpecValue, error) {
	return s.FindByGoodsIdAndType(goodsId, SpecValueTypeSingle)
}

// FindMultiSpecValuesByGoodsId 根据商品ID查询多规格值
func (s *ShopGoodsSpecValue) FindMultiSpecValuesByGoodsId(goodsId int64) ([]*ShopGoodsSpecValue, error) {
	return s.FindByGoodsIdAndType(goodsId, SpecValueTypeMulti)
}

// FindDefaultByGoodsId 根据商品ID查询默认套餐选项
func (s *ShopGoodsSpecValue) FindDefaultByGoodsId(goodsId int64) ([]*ShopGoodsSpecValue, error) {
	var values []*ShopGoodsSpecValue
	err := s.Query().Where("goods_id = ? AND spec_type = ? AND status = 1 AND JSON_EXTRACT(config, '$.is_default') = true", goodsId, SpecValueTypeCombo).Order("spec_group_id ASC, sort ASC").Find(&values).Error
	return values, err
}
