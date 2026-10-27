package goods

import (
	"errors"
	"fmt"
	"wnadmin/app/model/do/shop"
	"wnadmin/app/provider/db"

	"gorm.io/gorm"
)

// GetComboConfigService 获取套餐配置服务
type GetComboConfigService struct{}

// ComboConfig 套餐配置
type ComboConfig struct {
	ComboId     int64         `json:"comboId"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Price       float64       `json:"price"`
	Groups      []*ComboGroup `json:"groups"`
}

// ComboGroup 套餐组
type ComboGroup struct {
	Id          int64          `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	MinSelect   int            `json:"minSelect"`
	MaxSelect   int            `json:"maxSelect"`
	IsRequired  bool           `json:"isRequired"`
	Sort        int            `json:"sort"`
	Options     []*ComboOption `json:"options"`
}

// ComboOption 套餐选项
type ComboOption struct {
	Id            int64   `json:"id"`
	Name          string  `json:"name"`
	Value         string  `json:"value"`
	GoodsId       int64   `json:"goodsId"`
	SkuId         int64   `json:"skuId"`
	AddPrice      float64 `json:"addPrice"`
	Quantity      int     `json:"quantity"`
	IsDefault     bool    `json:"isDefault"`
	Sort          int     `json:"sort"`
	OriginalPrice float64 `json:"originalPrice"` // 原商品价格
	Stock         int     `json:"stock"`         // 库存
}

// Execute 执行获取套餐配置
func (s *GetComboConfigService) Execute(comboId int64) (*ComboConfig, error) {
	// 1. 获取套餐规格商品信息
	var goods shop.GoodsDO
	if err := db.WnAdmin().Where("id = ? AND spec_type = ? AND is_delete = 0", comboId, 3).First(&goods).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("套餐规格商品不存在")
		}
		return nil, fmt.Errorf("查询套餐规格商品失败: %v", err)
	}

	// 2. 获取套餐组（规格组）
	var specGroups []*shop.ShopGoodsSpecGroup
	if err := db.WnAdmin().Where("goods_id = ? AND status = 1", comboId).Order("sort ASC").Find(&specGroups).Error; err != nil {
		return nil, fmt.Errorf("查询套餐组失败: %v", err)
	}

	// 3. 获取套餐价格（从SKU表获取）
	var price float64 = 0
	var sku shop.ShopGoodsSku
	if err := db.WnAdmin().Where("goods_id = ? AND spec_type = ? AND status = 1", comboId, 3).First(&sku).Error; err == nil {
		price = sku.Price
	}

	// 4. 构建套餐配置
	config := &ComboConfig{
		ComboId:     goods.Id,
		Name:        goods.Name,
		Description: goods.Description,
		Price:       price, // 使用SKU表中的价格作为套餐价格
		Groups:      make([]*ComboGroup, 0, len(specGroups)),
	}

	// 5. 处理每个套餐组
	for _, specGroup := range specGroups {
		// 获取套餐组配置
		comboConfig, err := specGroup.GetComboConfig()
		if err != nil {
			return nil, fmt.Errorf("解析套餐组配置失败: %v", err)
		}

		group := &ComboGroup{
			Id:          specGroup.Id,
			Name:        specGroup.Name,
			Description: "",
			MinSelect:   0,
			MaxSelect:   1,
			IsRequired:  false,
			Sort:        specGroup.Sort,
			Options:     make([]*ComboOption, 0),
		}

		// 如果有套餐配置，使用配置中的值
		if comboConfig != nil {
			group.Description = comboConfig.Description
			group.MinSelect = comboConfig.MinSelect
			group.MaxSelect = comboConfig.MaxSelect
			group.IsRequired = comboConfig.IsRequired
		}

		// 6. 获取套餐选项（规格值）
		var specValues []*shop.ShopGoodsSpecValue
		if err := db.WnAdmin().Where("spec_group_id = ? AND status = 1", specGroup.Id).Order("sort ASC").Find(&specValues).Error; err != nil {
			return nil, fmt.Errorf("查询套餐选项失败: %v", err)
		}

		// 添加调试日志
		fmt.Printf("查询套餐选项 - 规格组ID: %d, 找到选项数量: %d\n", specGroup.Id, len(specValues))
		for i, value := range specValues {
			fmt.Printf("  选项%d: ID=%d, Name=%s, Value=%s, SpecType=%d\n", i+1, value.Id, value.Name, value.Value, value.SpecType)
		}

		// 7. 处理每个套餐选项
		for _, specValue := range specValues {
			option := &ComboOption{
				Id:        specValue.Id,
				Name:      specValue.Name,
				Value:     specValue.Value,
				GoodsId:   specValue.GetOptionGoodsId(),
				SkuId:     specValue.GetOptionSkuId(),
				AddPrice:  specValue.GetAddPrice(),
				Quantity:  specValue.GetQuantity(),
				IsDefault: specValue.IsDefault(),
				Sort:      specValue.Sort,
			}

			// 8. 获取选项商品的价格和库存信息
			if option.GoodsId > 0 {
				if option.SkuId > 0 {
					// 指定了SKU，获取SKU信息
					var sku shop.ShopGoodsSku
					if err := db.WnAdmin().Where("id = ?", option.SkuId).First(&sku).Error; err == nil {
						option.OriginalPrice = sku.Price
						option.Stock = sku.Stock
					}
				} else {
					// 未指定SKU，获取商品的默认SKU信息
					var sku shop.ShopGoodsSku
					if err := db.WnAdmin().Where("goods_id = ? AND status = 1", option.GoodsId).Order("sort ASC").First(&sku).Error; err == nil {
						option.OriginalPrice = sku.Price
						option.Stock = sku.Stock
					}
				}
			}

			group.Options = append(group.Options, option)
		}

		config.Groups = append(config.Groups, group)
	}

	return config, nil
}
