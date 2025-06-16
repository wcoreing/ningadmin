package shop

import (
	"fmt"
	"wnadmin/app/model/do/shop"
	"wnadmin/app/provider/db"

	"gorm.io/gorm"
)

// ComboServiceOptimized 优化的套餐服务
type ComboServiceOptimized struct{}

// ComboConfig 套餐配置
type ComboConfig struct {
	GoodsId   int64        `json:"goodsId"`   // 套餐商品ID
	BasePrice float64      `json:"basePrice"` // 基础价格
	Groups    []ComboGroup `json:"groups"`    // 套餐组列表
}

// ComboGroup 套餐组
type ComboGroup struct {
	Id          int64         `json:"id"`          // 组ID
	Name        string        `json:"name"`        // 组名称
	Description string        `json:"description"` // 组描述
	MinSelect   int           `json:"minSelect"`   // 最少选择数量
	MaxSelect   int           `json:"maxSelect"`   // 最多选择数量
	IsRequired  bool          `json:"isRequired"`  // 是否必选
	DisplayType string        `json:"displayType"` // 显示类型
	Sort        int           `json:"sort"`        // 排序
	Options     []ComboOption `json:"options"`     // 选项列表
}

// ComboOption 套餐选项
type ComboOption struct {
	Id            int64   `json:"id"`            // 选项ID
	Name          string  `json:"name"`          // 选项名称
	Value         string  `json:"value"`         // 选项值
	OptionGoodsId *int64  `json:"optionGoodsId"` // 关联商品ID
	OptionSkuId   *int64  `json:"optionSkuId"`   // 关联SKU ID
	AddPrice      float64 `json:"addPrice"`      // 加价金额
	Quantity      int     `json:"quantity"`      // 数量
	IsDefault     bool    `json:"isDefault"`     // 是否默认选择
	Description   string  `json:"description"`   // 选项描述
	Image         string  `json:"image"`         // 选项图片
	Sort          int     `json:"sort"`          // 排序
}

// ComboSelection 套餐选择
type ComboSelection struct {
	GroupId   int64   `json:"groupId"`   // 组ID
	OptionIds []int64 `json:"optionIds"` // 选择的选项ID列表
}

// ComboPriceResult 套餐价格计算结果
type ComboPriceResult struct {
	BasePrice  float64     `json:"basePrice"`  // 基础价格
	AddPrice   float64     `json:"addPrice"`   // 加价金额
	TotalPrice float64     `json:"totalPrice"` // 总价格
	Items      []ComboItem `json:"items"`      // 套餐项目明细
}

// ComboItem 套餐项目
type ComboItem struct {
	GroupName     string  `json:"groupName"`     // 组名称
	OptionName    string  `json:"optionName"`    // 选项名称
	OptionGoodsId *int64  `json:"optionGoodsId"` // 关联商品ID
	OptionSkuId   *int64  `json:"optionSkuId"`   // 关联SKU ID
	AddPrice      float64 `json:"addPrice"`      // 加价金额
	Quantity      int     `json:"quantity"`      // 数量
}

// GetComboConfig 获取套餐配置
func (s *ComboServiceOptimized) GetComboConfig(goodsId int64) (*ComboConfig, error) {
	// 验证商品是否为套餐规格类型
	var goods shop.GoodsDO
	if err := db.WnAdmin().Where("id = ? AND spec_type = ? AND is_delete = 0", goodsId, 3).First(&goods).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("套餐规格商品不存在")
		}
		return nil, fmt.Errorf("查询套餐规格商品失败: %v", err)
	}

	// 查询套餐组
	var specGroups []*shop.ShopGoodsSpecGroup
	if err := db.WnAdmin().Where("goods_id = ? AND spec_type = ? AND status = 1", goodsId, shop.SpecGroupTypeCombo).
		Order("sort ASC").Find(&specGroups).Error; err != nil {
		return nil, fmt.Errorf("查询套餐组失败: %v", err)
	}

	config := &ComboConfig{
		GoodsId:   goodsId,
		BasePrice: 0, // 套餐基础价格为0，由选项计算
		Groups:    make([]ComboGroup, 0, len(specGroups)),
	}

	// 构建套餐组配置
	for _, specGroup := range specGroups {
		group := ComboGroup{
			Id:   specGroup.Id,
			Name: specGroup.Name,
			Sort: specGroup.Sort,
		}

		// 解析组配置
		if groupConfig, err := specGroup.GetComboConfig(); err == nil && groupConfig != nil {
			group.Description = groupConfig.Description
			group.MinSelect = groupConfig.MinSelect
			group.MaxSelect = groupConfig.MaxSelect
			group.IsRequired = groupConfig.IsRequired
			group.DisplayType = groupConfig.DisplayType
		}

		// 查询组内选项
		var specValues []*shop.ShopGoodsSpecValue
		if err := db.WnAdmin().Where("spec_group_id = ? AND spec_type = ? AND status = 1", specGroup.Id, shop.SpecValueTypeCombo).
			Order("sort ASC").Find(&specValues).Error; err != nil {
			return nil, fmt.Errorf("查询套餐选项失败: %v", err)
		}

		group.Options = make([]ComboOption, 0, len(specValues))
		for _, specValue := range specValues {
			option := ComboOption{
				Id:    specValue.Id,
				Name:  specValue.Name,
				Value: specValue.Value,
				Sort:  specValue.Sort,
			}

			// 解析选项配置
			if optionConfig, err := specValue.GetComboConfig(); err == nil && optionConfig != nil {
				option.OptionGoodsId = optionConfig.OptionGoodsId
				option.OptionSkuId = optionConfig.OptionSkuId
				option.AddPrice = optionConfig.AddPrice
				option.Quantity = optionConfig.Quantity
				option.IsDefault = optionConfig.IsDefault
				option.Description = optionConfig.Description
				option.Image = optionConfig.Image
			}

			group.Options = append(group.Options, option)
		}

		config.Groups = append(config.Groups, group)
	}

	return config, nil
}

// CalculateComboPrice 计算套餐价格
func (s *ComboServiceOptimized) CalculateComboPrice(goodsId int64, selections []ComboSelection) (*ComboPriceResult, error) {
	// 获取套餐配置
	config, err := s.GetComboConfig(goodsId)
	if err != nil {
		return nil, err
	}

	result := &ComboPriceResult{
		BasePrice:  config.BasePrice,
		AddPrice:   0,
		TotalPrice: config.BasePrice,
		Items:      make([]ComboItem, 0),
	}

	// 创建选择映射
	selectionMap := make(map[int64][]int64)
	for _, selection := range selections {
		selectionMap[selection.GroupId] = selection.OptionIds
	}

	// 计算每个组的价格
	for _, group := range config.Groups {
		selectedOptionIds, exists := selectionMap[group.Id]
		if !exists || len(selectedOptionIds) == 0 {
			if group.IsRequired {
				return nil, fmt.Errorf("必选组 '%s' 未选择选项", group.Name)
			}
			continue
		}

		// 验证选择数量
		if len(selectedOptionIds) < group.MinSelect {
			return nil, fmt.Errorf("组 '%s' 至少需要选择 %d 个选项", group.Name, group.MinSelect)
		}
		if len(selectedOptionIds) > group.MaxSelect {
			return nil, fmt.Errorf("组 '%s' 最多只能选择 %d 个选项", group.Name, group.MaxSelect)
		}

		// 计算选中选项的价格
		for _, optionId := range selectedOptionIds {
			var selectedOption *ComboOption
			for _, option := range group.Options {
				if option.Id == optionId {
					selectedOption = &option
					break
				}
			}

			if selectedOption == nil {
				return nil, fmt.Errorf("组 '%s' 中不存在选项ID %d", group.Name, optionId)
			}

			// 累加价格
			result.AddPrice += selectedOption.AddPrice

			// 添加到明细
			result.Items = append(result.Items, ComboItem{
				GroupName:     group.Name,
				OptionName:    selectedOption.Name,
				OptionGoodsId: selectedOption.OptionGoodsId,
				OptionSkuId:   selectedOption.OptionSkuId,
				AddPrice:      selectedOption.AddPrice,
				Quantity:      selectedOption.Quantity,
			})
		}
	}

	result.TotalPrice = result.BasePrice + result.AddPrice
	return result, nil
}

// ValidateComboSelection 验证套餐选择
func (s *ComboServiceOptimized) ValidateComboSelection(goodsId int64, selections []ComboSelection) error {
	_, err := s.CalculateComboPrice(goodsId, selections)
	return err
}

// GetDefaultComboSelection 获取默认套餐选择
func (s *ComboServiceOptimized) GetDefaultComboSelection(goodsId int64) ([]ComboSelection, error) {
	config, err := s.GetComboConfig(goodsId)
	if err != nil {
		return nil, err
	}

	selections := make([]ComboSelection, 0)

	for _, group := range config.Groups {
		var defaultOptionIds []int64

		// 查找默认选项
		for _, option := range group.Options {
			if option.IsDefault {
				defaultOptionIds = append(defaultOptionIds, option.Id)
			}
		}

		// 如果没有默认选项但是必选组，选择第一个选项
		if len(defaultOptionIds) == 0 && group.IsRequired && len(group.Options) > 0 {
			defaultOptionIds = append(defaultOptionIds, group.Options[0].Id)
		}

		// 确保选择数量在范围内
		if len(defaultOptionIds) > group.MaxSelect {
			defaultOptionIds = defaultOptionIds[:group.MaxSelect]
		}

		if len(defaultOptionIds) > 0 {
			selections = append(selections, ComboSelection{
				GroupId:   group.Id,
				OptionIds: defaultOptionIds,
			})
		}
	}

	return selections, nil
}

// CreateCombo 创建套餐
func (s *ComboServiceOptimized) CreateCombo(goodsId int64, groups []ComboGroup) error {
	return db.WnAdmin().Transaction(func(tx *gorm.DB) error {
		// 删除现有的套餐组和选项
		if err := tx.Where("goods_id = ? AND spec_type = ?", goodsId, shop.SpecGroupTypeCombo).Delete(&shop.ShopGoodsSpecGroup{}).Error; err != nil {
			return fmt.Errorf("删除现有套餐组失败: %v", err)
		}

		if err := tx.Where("goods_id = ? AND spec_type = ?", goodsId, shop.SpecValueTypeCombo).Delete(&shop.ShopGoodsSpecValue{}).Error; err != nil {
			return fmt.Errorf("删除现有套餐选项失败: %v", err)
		}

		// 创建新的套餐组和选项
		for _, group := range groups {
			// 创建套餐组
			specGroup := &shop.ShopGoodsSpecGroup{
				GoodsId:  goodsId,
				Name:     group.Name,
				SpecType: shop.SpecGroupTypeCombo,
				Sort:     group.Sort,
				Status:   1,
			}

			// 设置组配置
			groupConfig := &shop.ComboGroupConfig{
				Description: group.Description,
				MinSelect:   group.MinSelect,
				MaxSelect:   group.MaxSelect,
				IsRequired:  group.IsRequired,
				DisplayType: group.DisplayType,
			}

			if err := specGroup.SetComboConfig(groupConfig); err != nil {
				return fmt.Errorf("设置套餐组配置失败: %v", err)
			}

			if err := tx.Create(specGroup).Error; err != nil {
				return fmt.Errorf("创建套餐组失败: %v", err)
			}

			// 创建套餐选项
			for _, option := range group.Options {
				specValue := &shop.ShopGoodsSpecValue{
					GoodsId:     goodsId,
					SpecGroupId: specGroup.Id,
					Name:        option.Name,
					Value:       option.Value,
					SpecType:    shop.SpecValueTypeCombo,
					Sort:        option.Sort,
					Status:      1,
				}

				// 设置选项配置
				optionConfig := &shop.ComboOptionConfig{
					OptionGoodsId: option.OptionGoodsId,
					OptionSkuId:   option.OptionSkuId,
					AddPrice:      option.AddPrice,
					Quantity:      option.Quantity,
					IsDefault:     option.IsDefault,
					Description:   option.Description,
					Image:         option.Image,
				}

				if err := specValue.SetComboConfig(optionConfig); err != nil {
					return fmt.Errorf("设置套餐选项配置失败: %v", err)
				}

				if err := tx.Create(specValue).Error; err != nil {
					return fmt.Errorf("创建套餐选项失败: %v", err)
				}
			}
		}

		return nil
	})
}
