package goods

import (
	"fmt"
	"wnadmin/app/model/do/shop"
	"wnadmin/app/provider/db"
)

// GetDefaultComboSelectionService 获取默认套餐选择服务
type GetDefaultComboSelectionService struct{}

// Execute 执行获取默认套餐选择
func (s *GetDefaultComboSelectionService) Execute(comboId int64) ([]*ComboSelection, error) {
	// 获取所有默认选项（使用JSON查询）
	var specValues []*shop.ShopGoodsSpecValue
	if err := db.WnAdmin().Where("goods_id = ? AND spec_type = ? AND status = 1 AND JSON_EXTRACT(config, '$.is_default') = true", comboId, shop.SpecValueTypeCombo).Find(&specValues).Error; err != nil {
		return nil, fmt.Errorf("查询默认选项失败: %v", err)
	}

	selections := make([]*ComboSelection, 0, len(specValues))
	for _, specValue := range specValues {
		selection := &ComboSelection{
			GroupId:       specValue.SpecGroupId,
			SpecValueId:   specValue.Id,
			OptionGoodsId: specValue.GetOptionGoodsId(),
			Quantity:      specValue.GetQuantity(),
		}
		selections = append(selections, selection)
	}

	return selections, nil
}
