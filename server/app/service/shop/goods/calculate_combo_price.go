package goods

import (
	"fmt"
	"wnadmin/app/model/do/shop"
	"wnadmin/app/provider/db"
)

// CalculateComboService 计算套餐价格服务
type CalculateComboService struct{}

// ComboSelection 套餐选择
type ComboSelection struct {
	GroupId       int64 `json:"groupId"`
	SpecValueId   int64 `json:"specValueId"`
	OptionGoodsId int64 `json:"optionGoodsId"`
	Quantity      int   `json:"quantity"`
}

// Execute 执行套餐价格计算
func (s *CalculateComboService) Execute(comboId int64, selections []*ComboSelection) (float64, error) {
	// 1. 获取套餐规格商品基础价格（从SKU表获取）
	var goods shop.GoodsDO
	if err := db.WnAdmin().Where("id = ? AND spec_type = ? AND is_delete = 0", comboId, 3).First(&goods).Error; err != nil {
		return 0, fmt.Errorf("套餐规格商品不存在: %v", err)
	}

	// 2. 获取套餐商品的SKU价格作为基础价格
	var sku shop.ShopGoodsSku
	if err := db.WnAdmin().Where("goods_id = ? AND spec_type = ? AND status = 1", comboId, 3).First(&sku).Error; err != nil {
		return 0, fmt.Errorf("套餐SKU不存在: %v", err)
	}

	totalPrice := sku.Price // 使用SKU的price字段作为套餐基础价格

	// 2. 计算选项加价
	for _, selection := range selections {
		var specValue shop.ShopGoodsSpecValue
		if err := db.WnAdmin().Where("id = ?", selection.SpecValueId).First(&specValue).Error; err != nil {
			return 0, fmt.Errorf("套餐选项不存在: %v", err)
		}

		// 加价 = 选项加价 * 数量
		addPrice := specValue.GetAddPrice() * float64(selection.Quantity)
		totalPrice += addPrice
	}

	return totalPrice, nil
}
