package goods

import (
	"fmt"
	"time"
	"wnadmin/app/common/utils"
	"wnadmin/app/model/do/shop"
	"wnadmin/app/model/dto"
	"wnadmin/app/provider/db"

	"gorm.io/gorm"
)

// MultiSpecService 多规格数据管理服务
type MultiSpecService struct{}

// SaveAllSpecData 保存所有规格类型的数据
func (s *MultiSpecService) SaveAllSpecData(tx *gorm.DB, goodsId int64, req *dto.GoodsUpdateReq) error {
	fmt.Printf("MultiSpecService - 开始保存多规格数据，商品ID: %d, specType: %d\n", goodsId, req.SpecType)
	fmt.Printf("MultiSpecService - 单规格数据: %+v\n", req.SingleSpecData)
	fmt.Printf("MultiSpecService - 多规格数据: %+v\n", req.MultiSpecData)
	fmt.Printf("MultiSpecService - 套餐规格数据: %+v\n", req.ComboSpecData)

	// 详细打印套餐规格数据
	if req.ComboSpecData != nil {
		for i, group := range req.ComboSpecData.ComboGroups {
			fmt.Printf("MultiSpecService - 套餐组[%d]: %s\n", i, group.Name)
			for j, option := range group.Options {
				fmt.Printf("MultiSpecService - 选项[%d]: 名称=%s, 商品ID=%d, 图片=%s\n", j, option.Name, option.GoodsId, option.GoodsImage)
			}
		}
	}

	// 根据当前的specType决定保存哪种类型的数据，避免相互覆盖
	switch req.SpecType {
	case 1: // 单规格
		if req.SingleSpecData != nil {
			fmt.Printf("MultiSpecService - 保存单规格数据\n")
			if err := s.saveSingleSpecData(tx, goodsId, req.SingleSpecData); err != nil {
				return fmt.Errorf("保存单规格数据失败: %v", err)
			}
		}
	case 2: // 多规格
		if req.MultiSpecData != nil {
			fmt.Printf("MultiSpecService - 保存多规格数据\n")
			if err := s.saveMultiSpecData(tx, goodsId, req.MultiSpecData); err != nil {
				return fmt.Errorf("保存多规格数据失败: %v", err)
			}
		}
	case 3: // 套餐规格
		if req.ComboSpecData != nil {
			fmt.Printf("MultiSpecService - 保存套餐规格数据\n")
			if err := s.saveComboSpecData(tx, goodsId, req.ComboSpecData); err != nil {
				return fmt.Errorf("保存套餐规格数据失败: %v", err)
			}
		}
	}

	// 兼容旧版本数据格式（只有在新格式数据都为空时才执行）
	if req.SingleSpecData == nil && req.MultiSpecData == nil && req.ComboSpecData == nil {
		if err := s.saveCompatibleData(tx, goodsId, req); err != nil {
			return fmt.Errorf("保存兼容数据失败: %v", err)
		}
	}

	fmt.Printf("MultiSpecService - 多规格数据保存完成\n")
	return nil
}

// LoadAllSpecData 加载所有规格类型的数据
func (s *MultiSpecService) LoadAllSpecData(goodsId int64) (*dto.GoodsInfoResp, error) {
	resp := &dto.GoodsInfoResp{}

	// 1. 加载单规格数据
	singleData, err := s.loadSingleSpecData(goodsId)
	if err != nil {
		return nil, fmt.Errorf("加载单规格数据失败: %v", err)
	}
	resp.SingleSpecData = singleData

	// 2. 加载多规格数据
	multiData, err := s.loadMultiSpecData(goodsId)
	if err != nil {
		return nil, fmt.Errorf("加载多规格数据失败: %v", err)
	}
	resp.MultiSpecData = multiData

	// 3. 加载套餐规格数据
	comboData, err := s.loadComboSpecData(goodsId)
	if err != nil {
		return nil, fmt.Errorf("加载套餐规格数据失败: %v", err)
	}
	resp.ComboSpecData = comboData

	return resp, nil
}

// saveSingleSpecData 保存单规格数据
func (s *MultiSpecService) saveSingleSpecData(tx *gorm.DB, goodsId int64, data *dto.SingleSpecData) error {
	fmt.Printf("MultiSpecService - 开始删除单规格SKU\n")

	// 优化删除操作
	if err := tx.Where("goods_id = ? AND spec_type = ?", goodsId, shop.SpecGroupTypeSingle).Delete(&shop.ShopGoodsSku{}).Error; err != nil {
		fmt.Printf("删除单规格SKU失败: %v\n", err)
		return err
	}

	fmt.Printf("MultiSpecService - 单规格SKU删除完成\n")

	// 创建单规格SKU
	sku := &shop.ShopGoodsSku{
		GoodsId:       goodsId,
		SpecType:      shop.SpecGroupTypeSingle,
		SkuCode:       fmt.Sprintf("SINGLE_%d_%d", goodsId, time.Now().Unix()),
		Price:         data.Price,
		OriginalPrice: data.OriginalPrice,
		CostPrice:     data.CostPrice,
		Stock:         data.Stock,
		Weight:        data.Weight,
		Volume:        data.Volume,
		Sort:          0,
		Status:        1,
		CreateTime:    time.Now(),
		UpdateTime:    time.Now(),
	}

	fmt.Printf("MultiSpecService - 创建单规格SKU: %+v\n", sku)
	if err := tx.Create(sku).Error; err != nil {
		fmt.Printf("MultiSpecService - 创建单规格SKU失败: %v\n", err)
		return err
	}
	fmt.Printf("MultiSpecService - 单规格SKU创建成功，ID: %d\n", sku.Id)
	return nil
}

// saveMultiSpecData 保存多规格数据
func (s *MultiSpecService) saveMultiSpecData(tx *gorm.DB, goodsId int64, data *dto.MultiSpecData) error {
	fmt.Printf("MultiSpecService - 开始删除多规格数据\n")

	// 优化删除操作，先删除关联数据，再删除主数据
	// 1. 删除SKU规格关联
	if err := tx.Exec("DELETE FROM shop_goods_sku_spec WHERE sku_id IN (SELECT id FROM shop_goods_sku WHERE goods_id = ? AND spec_type = ?)", goodsId, shop.SpecGroupTypeMulti).Error; err != nil {
		fmt.Printf("删除SKU规格关联失败: %v\n", err)
		return err
	}

	// 2. 删除SKU
	if err := tx.Where("goods_id = ? AND spec_type = ?", goodsId, shop.SpecGroupTypeMulti).Delete(&shop.ShopGoodsSku{}).Error; err != nil {
		fmt.Printf("删除SKU失败: %v\n", err)
		return err
	}

	// 3. 删除规格值
	if err := tx.Where("goods_id = ? AND spec_type = ?", goodsId, shop.SpecGroupTypeMulti).Delete(&shop.ShopGoodsSpecValue{}).Error; err != nil {
		fmt.Printf("删除规格值失败: %v\n", err)
		return err
	}

	// 4. 删除规格组
	if err := tx.Where("goods_id = ? AND spec_type = ?", goodsId, shop.SpecGroupTypeMulti).Delete(&shop.ShopGoodsSpecGroup{}).Error; err != nil {
		fmt.Printf("删除规格组失败: %v\n", err)
		return err
	}

	fmt.Printf("MultiSpecService - 多规格数据删除完成\n")

	// 如果没有新数据，直接返回
	if len(data.SpecGroups) == 0 {
		fmt.Printf("MultiSpecService - 没有多规格数据需要保存\n")
		return nil
	}

	fmt.Printf("MultiSpecService - 开始创建多规格数据\n")
	// 保存规格组和规格值
	specMaps, err := s.createSpecGroupsAndValues(tx, goodsId, shop.SpecGroupTypeMulti, data.SpecGroups)
	if err != nil {
		fmt.Printf("创建规格组和规格值失败: %v\n", err)
		return err
	}

	// 保存SKU
	if err := s.createSkus(tx, goodsId, shop.SpecGroupTypeMulti, data.SkuList, specMaps); err != nil {
		fmt.Printf("创建SKU失败: %v\n", err)
		return err
	}

	fmt.Printf("MultiSpecService - 多规格数据创建完成\n")
	return nil
}

// saveComboSpecData 保存套餐规格数据
func (s *MultiSpecService) saveComboSpecData(tx *gorm.DB, goodsId int64, data *dto.ComboSpecData) error {
	// 优化删除操作
	// 1. 删除SKU规格关联
	if err := tx.Where("sku_id IN (SELECT id FROM shop_goods_sku WHERE goods_id = ? AND spec_type = ?)", goodsId, shop.SpecGroupTypeCombo).Delete(&shop.ShopGoodsSkuSpec{}).Error; err != nil {
		return err
	}

	// 2. 删除SKU
	if err := tx.Where("goods_id = ? AND spec_type = ?", goodsId, shop.SpecGroupTypeCombo).Delete(&shop.ShopGoodsSku{}).Error; err != nil {
		return err
	}

	// 3. 删除规格值
	if err := tx.Where("goods_id = ? AND spec_type = ?", goodsId, shop.SpecGroupTypeCombo).Delete(&shop.ShopGoodsSpecValue{}).Error; err != nil {
		return err
	}

	// 4. 删除规格组
	if err := tx.Where("goods_id = ? AND spec_type = ?", goodsId, shop.SpecGroupTypeCombo).Delete(&shop.ShopGoodsSpecGroup{}).Error; err != nil {
		return err
	}

	// 保存套餐规格组和规格值（包含图片信息）
	_, err := s.createComboSpecGroupsAndValues(tx, goodsId, data.ComboGroups)
	if err != nil {
		return err
	}

	// 创建套餐商品的单一SKU（类似单规格商品）
	return s.createComboSingleSku(tx, goodsId, data)
}

// saveCompatibleData 保存兼容旧版本的数据
func (s *MultiSpecService) saveCompatibleData(tx *gorm.DB, goodsId int64, req *dto.GoodsUpdateReq) error {
	// 如果新格式数据都为空，则使用旧格式数据
	if req.SingleSpecData == nil && req.MultiSpecData == nil && req.ComboSpecData == nil {
		// 根据specType决定保存哪种类型的数据
		switch req.SpecType {
		case 1: // 单规格
			singleData := &dto.SingleSpecData{
				Price:         req.Price,
				OriginalPrice: req.OriginalPrice,
				CostPrice:     req.CostPrice,
				Stock:         req.Stock,
				Weight:        req.Weight,
				Volume:        req.Volume,
			}
			return s.saveSingleSpecData(tx, goodsId, singleData)

		case 2: // 多规格
			multiData := &dto.MultiSpecData{
				SpecGroups: req.SpecGroups,
				SkuList:    req.SkuList,
			}
			return s.saveMultiSpecData(tx, goodsId, multiData)

		case 3: // 套餐规格
			comboData := &dto.ComboSpecData{
				ComboGroups: req.ComboGroups,
			}
			return s.saveComboSpecData(tx, goodsId, comboData)
		}
	}

	return nil
}

// createComboSpecGroupsAndValues 创建套餐规格组和规格值（包含图片信息）
func (s *MultiSpecService) createComboSpecGroupsAndValues(tx *gorm.DB, goodsId int64, comboGroups []dto.ComboSpecGroup) (*SpecMaps, error) {
	specGroupMap := make(map[string]int64)
	specValueMap := make(map[string]int64)

	for _, comboGroup := range comboGroups {
		if comboGroup.Name == "" || len(comboGroup.Options) == 0 {
			continue
		}

		// 创建规格组
		specGroup := &shop.ShopGoodsSpecGroup{
			GoodsId:  goodsId,
			Name:     comboGroup.Name,
			SpecType: shop.SpecGroupTypeCombo,
			Sort:     comboGroup.Sort,
			Status:   1,
		}

		// 设置套餐组配置
		config := &shop.ComboGroupConfig{
			Description: comboGroup.Description,
			MinSelect:   comboGroup.MinSelect,
			MaxSelect:   comboGroup.MaxSelect,
			IsRequired:  comboGroup.IsRequired,
			DisplayType: "radio", // 默认单选
		}
		if err := specGroup.SetComboConfig(config); err != nil {
			return nil, err
		}

		if err := tx.Create(specGroup).Error; err != nil {
			return nil, err
		}

		specGroupMap[comboGroup.Name] = specGroup.Id

		// 创建规格值（套餐选项）
		for _, option := range comboGroup.Options {
			specValue := &shop.ShopGoodsSpecValue{
				GoodsId:     goodsId,
				SpecGroupId: specGroup.Id,
				Name:        option.Name,
				Value:       option.Value,
				SpecType:    shop.SpecGroupTypeCombo,
				Sort:        option.Sort,
				Status:      1,
			}

			// 设置套餐选项配置（包含图片信息）
			optionConfig := &shop.ComboOptionConfig{
				OptionGoodsId: &option.GoodsId,
				OptionSkuId:   option.SkuId,
				AddPrice:      option.AddPrice,
				Quantity:      option.Quantity,
				IsDefault:     option.IsDefault,
				Image:         option.GoodsImage, // 保存图片信息
			}
			if err := specValue.SetComboConfig(optionConfig); err != nil {
				return nil, err
			}

			if err := tx.Create(specValue).Error; err != nil {
				return nil, err
			}

			specValueKey := fmt.Sprintf("%s:%s", comboGroup.Name, option.Value)
			specValueMap[specValueKey] = specValue.Id
		}
	}

	return &SpecMaps{
		GroupMap: specGroupMap,
		ValueMap: specValueMap,
	}, nil
}

// createComboSingleSku 为套餐商品创建单一SKU（类似单规格商品）
func (s *MultiSpecService) createComboSingleSku(tx *gorm.DB, goodsId int64, data *dto.ComboSpecData) error {
	// 获取商品基础信息来设置SKU价格
	var goods shop.GoodsDO
	if err := tx.Where("id = ?", goodsId).First(&goods).Error; err != nil {
		return fmt.Errorf("获取商品信息失败: %v", err)
	}

	// 确定套餐价格
	var comboPrice float64 = 0
	if data != nil && data.Price > 0 {
		comboPrice = data.Price
		fmt.Printf("MultiSpecService - 使用套餐规格数据中的价格: %.2f\n", comboPrice)
	}

	// 创建套餐商品的单一SKU
	sku := &shop.ShopGoodsSku{
		GoodsId:       goodsId,
		SpecType:      shop.SpecGroupTypeCombo,
		SkuCode:       fmt.Sprintf("COMBO_%d_%d", goodsId, time.Now().UnixNano()),
		Image:         goods.CoverImg, // 使用商品封面图
		Price:         comboPrice,     // 使用套餐规格数据中的价格
		OriginalPrice: 0,              // 原价也为0
		CostPrice:     0,              // 套餐商品成本价可以为0
		Stock:         999,            // 套餐默认库存
		Weight:        0,              // 套餐商品重量为0
		Volume:        0,              // 套餐商品体积为0
		Sort:          0,
		Status:        1,
	}

	if err := tx.Create(sku).Error; err != nil {
		return fmt.Errorf("创建套餐SKU失败: %v", err)
	}

	return nil
}

// createSpecGroupsAndValues 创建规格组和规格值
func (s *MultiSpecService) createSpecGroupsAndValues(tx *gorm.DB, goodsId int64, specType int8, specGroups []dto.SpecGroup) (*SpecMaps, error) {
	specGroupMap := make(map[string]int64)
	specValueMap := make(map[string]int64)

	for _, group := range specGroups {
		if group.Name == "" || len(group.Values) == 0 {
			continue
		}

		// 创建规格组
		specGroup := &shop.ShopGoodsSpecGroup{
			GoodsId:  goodsId,
			Name:     group.Name,
			SpecType: specType,
			Sort:     group.Sort,
			Status:   1,
		}

		// 如果是套餐规格，设置配置
		if specType == shop.SpecGroupTypeCombo {
			config := &shop.ComboGroupConfig{
				Description: group.Description,
				MinSelect:   group.MinSelect,
				MaxSelect:   group.MaxSelect,
				IsRequired:  group.IsRequired,
				DisplayType: "radio", // 默认单选
			}
			if err := specGroup.SetComboConfig(config); err != nil {
				return nil, err
			}
		}

		if err := tx.Create(specGroup).Error; err != nil {
			return nil, err
		}

		specGroupMap[group.Name] = specGroup.Id

		// 创建规格值
		for _, value := range group.Values {
			specValue := &shop.ShopGoodsSpecValue{
				GoodsId:     goodsId,
				SpecGroupId: specGroup.Id,
				Name:        value.Name,
				Value:       value.Value,
				SpecType:    specType,
				Sort:        value.Sort,
				Status:      1,
			}

			// 如果是套餐规格，设置配置
			if specType == shop.SpecGroupTypeCombo {
				config := &shop.ComboOptionConfig{
					OptionGoodsId: value.OptionGoodsId,
					OptionSkuId:   value.OptionSkuId,
					AddPrice:      value.AddPrice,
					Quantity:      value.Quantity,
					IsDefault:     value.IsDefault,
					Image:         "", // 先设置为空，后面从套餐选项数据中获取
				}
				if err := specValue.SetComboConfig(config); err != nil {
					return nil, err
				}
			}

			if err := tx.Create(specValue).Error; err != nil {
				return nil, err
			}

			specValueKey := fmt.Sprintf("%s:%s", group.Name, value.Value)
			specValueMap[specValueKey] = specValue.Id
		}
	}

	return &SpecMaps{
		GroupMap: specGroupMap,
		ValueMap: specValueMap,
	}, nil
}

// createSkus 创建SKU
func (s *MultiSpecService) createSkus(tx *gorm.DB, goodsId int64, specType int8, skuList []dto.Sku, specMaps *SpecMaps) error {
	for _, sku := range skuList {
		// 去除图片URL中的域名
		cleanImageUrl := utils.RemoveDomainFromURL(sku.Image)

		skuData := &shop.ShopGoodsSku{
			GoodsId:       goodsId,
			SpecType:      specType,
			SkuCode:       fmt.Sprintf("%s_%d_%d", getSpecTypePrefix(specType), goodsId, time.Now().UnixNano()),
			Image:         cleanImageUrl,
			Price:         sku.Price,
			OriginalPrice: sku.OriginalPrice,
			CostPrice:     sku.CostPrice,
			Stock:         sku.Stock,
			Weight:        sku.Weight,
			Volume:        sku.Volume,
			Sort:          0,
			Status:        1,
		}

		if err := tx.Create(skuData).Error; err != nil {
			return err
		}

		// 创建SKU规格关联
		if err := s.createSkuSpecRelations(tx, skuData.Id, sku.Specs, specMaps); err != nil {
			return err
		}
	}

	return nil
}

// createSkuSpecRelations 创建SKU规格关联
func (s *MultiSpecService) createSkuSpecRelations(tx *gorm.DB, skuId int64, specs map[string]dto.SpecValue, specMaps *SpecMaps) error {
	for groupName, specValue := range specs {
		groupId, exists := specMaps.GroupMap[groupName]
		if !exists {
			continue
		}

		specValueKey := fmt.Sprintf("%s:%s", groupName, specValue.Value)
		valueId, exists := specMaps.ValueMap[specValueKey]
		if !exists {
			continue
		}

		relation := &shop.ShopGoodsSkuSpec{
			SkuId:       skuId,
			SpecGroupId: groupId,
			SpecValueId: valueId,
		}

		if err := tx.Create(relation).Error; err != nil {
			return err
		}
	}

	return nil
}

// loadSingleSpecData 加载单规格数据
func (s *MultiSpecService) loadSingleSpecData(goodsId int64) (*dto.SingleSpecData, error) {
	skuModel := &shop.ShopGoodsSku{}
	skus, err := skuModel.GetByGoodsIdAndSpecType(goodsId, shop.SpecGroupTypeSingle)
	if err != nil {
		return nil, err
	}

	if len(skus) == 0 {
		return nil, nil
	}

	sku := skus[0] // 单规格只有一个SKU
	return &dto.SingleSpecData{
		Price:         sku.Price,
		OriginalPrice: sku.OriginalPrice,
		CostPrice:     sku.CostPrice,
		Stock:         sku.Stock,
		Weight:        sku.Weight,
		Volume:        sku.Volume,
	}, nil
}

// loadMultiSpecData 加载多规格数据
func (s *MultiSpecService) loadMultiSpecData(goodsId int64) (*dto.MultiSpecData, error) {
	// 加载规格组和规格值
	specGroups, err := s.loadSpecGroups(goodsId, shop.SpecGroupTypeMulti)
	if err != nil {
		return nil, err
	}

	// 加载SKU
	skus, err := s.loadSkus(goodsId, shop.SpecGroupTypeMulti)
	if err != nil {
		return nil, err
	}

	if len(specGroups) == 0 && len(skus) == 0 {
		return nil, nil
	}

	return &dto.MultiSpecData{
		SpecGroups: specGroups,
		SkuList:    skus,
	}, nil
}

// loadComboSpecData 加载套餐规格数据
func (s *MultiSpecService) loadComboSpecData(goodsId int64) (*dto.ComboSpecData, error) {
	// 加载规格组和规格值
	specGroups, err := s.loadSpecGroups(goodsId, shop.SpecGroupTypeCombo)
	if err != nil {
		return nil, err
	}

	if len(specGroups) == 0 {
		return nil, nil
	}

	// 转换为套餐组格式
	var comboGroups []dto.ComboSpecGroup
	for _, group := range specGroups {
		comboGroup := dto.ComboSpecGroup{
			ID:          group.ID,
			Name:        group.Name,
			Description: group.Description,
			MinSelect:   group.MinSelect,
			MaxSelect:   group.MaxSelect,
			IsRequired:  group.IsRequired,
			Sort:        group.Sort,
		}

		// 转换规格值为套餐选项
		for _, value := range group.Values {
			option := dto.ComboSpecOption{
				ID:        value.ID,
				Name:      value.Name,
				Value:     value.Value,
				AddPrice:  value.AddPrice,
				Quantity:  value.Quantity,
				IsDefault: value.IsDefault,
				Sort:      value.Sort,
			}

			if value.OptionGoodsId != nil {
				option.GoodsId = *value.OptionGoodsId
			}
			if value.OptionSkuId != nil {
				option.SkuId = value.OptionSkuId
			}

			// 从规格值配置中获取图片和关联商品信息
			specValueModel := &shop.ShopGoodsSpecValue{}
			if err := db.WnAdmin().First(specValueModel, value.ID).Error; err == nil {
				if config, err := specValueModel.GetComboConfig(); err == nil && config != nil {
					option.GoodsImage = config.Image // 从配置中获取图片

					// 如果配置中没有图片，尝试从关联商品信息中获取
					if option.GoodsImage == "" && config.OptionGoodsId != nil {
						var goods shop.GoodsDO
						if err := db.WnAdmin().First(&goods, *config.OptionGoodsId).Error; err == nil {
							option.GoodsImage = goods.CoverImg
							option.GoodsName = goods.Name
						}
					}
				}
			}

			// 如果还是没有商品名称，尝试从关联商品中获取
			if option.GoodsName == "" && option.GoodsId > 0 {
				var goods shop.GoodsDO
				if err := db.WnAdmin().First(&goods, option.GoodsId).Error; err == nil {
					option.GoodsName = goods.Name
				}
			}

			comboGroup.Options = append(comboGroup.Options, option)
		}
		comboGroups = append(comboGroups, comboGroup)
	}

	// 获取套餐SKU的价格
	var comboPrice float64 = 0
	skuModel := &shop.ShopGoodsSku{}
	skus, err := skuModel.GetByGoodsIdAndSpecType(goodsId, shop.SpecGroupTypeCombo)
	if err == nil && len(skus) > 0 {
		comboPrice = skus[0].Price // 套餐商品只有一个SKU
		fmt.Printf("MultiSpecService - 从SKU中读取套餐价格: %.2f\n", comboPrice)
	}

	return &dto.ComboSpecData{
		Price:       comboPrice,
		ComboGroups: comboGroups,
	}, nil
}

// loadSpecGroups 加载规格组
func (s *MultiSpecService) loadSpecGroups(goodsId int64, specType int8) ([]dto.SpecGroup, error) {
	specGroupModel := &shop.ShopGoodsSpecGroup{}
	groups, err := specGroupModel.GetByGoodsIdAndSpecType(goodsId, specType)
	if err != nil {
		return nil, err
	}

	var result []dto.SpecGroup
	for _, group := range groups {
		// 加载规格值
		specValueModel := &shop.ShopGoodsSpecValue{}
		values, err := specValueModel.GetByGoodsIdAndSpecType(goodsId, specType)
		if err != nil {
			return nil, err
		}

		// 过滤属于当前组的规格值
		var groupValues []dto.SpecValue
		for _, value := range values {
			if value.SpecGroupId == group.Id {
				specValue := dto.SpecValue{
					ID:    value.Id,
					Name:  value.Name,
					Value: value.Value,
					Sort:  value.Sort,
				}

				// 如果是套餐规格，加载配置
				if specType == shop.SpecGroupTypeCombo {
					if config, err := value.GetComboConfig(); err == nil && config != nil {
						specValue.OptionGoodsId = config.OptionGoodsId
						specValue.OptionSkuId = config.OptionSkuId
						specValue.AddPrice = config.AddPrice
						specValue.Quantity = config.Quantity
						specValue.IsDefault = config.IsDefault
						// 注意：这里不设置图片，因为SpecValue结构体没有图片字段
						// 图片信息会在转换为ComboSpecOption时从配置中获取
					}
				}

				groupValues = append(groupValues, specValue)
			}
		}

		specGroup := dto.SpecGroup{
			ID:     group.Id,
			Name:   group.Name,
			Sort:   group.Sort,
			Values: groupValues,
		}

		// 如果是套餐规格，加载配置
		if specType == shop.SpecGroupTypeCombo {
			if config, err := group.GetComboConfig(); err == nil && config != nil {
				specGroup.Description = config.Description
				specGroup.MinSelect = config.MinSelect
				specGroup.MaxSelect = config.MaxSelect
				specGroup.IsRequired = config.IsRequired
			}
		}

		result = append(result, specGroup)
	}

	return result, nil
}

// loadSkus 加载SKU
func (s *MultiSpecService) loadSkus(goodsId int64, specType int8) ([]dto.Sku, error) {
	skuModel := &shop.ShopGoodsSku{}
	skus, err := skuModel.GetByGoodsIdAndSpecType(goodsId, specType)
	if err != nil {
		return nil, err
	}

	var result []dto.Sku
	for _, sku := range skus {
		// 加载SKU规格关联
		specs, err := s.loadSkuSpecs(sku.Id)
		if err != nil {
			return nil, err
		}

		dtoSku := dto.Sku{
			ID:            sku.Id,
			Specs:         specs,
			Price:         sku.Price,
			OriginalPrice: sku.OriginalPrice,
			CostPrice:     sku.CostPrice,
			Stock:         sku.Stock,
			Weight:        sku.Weight,
			Volume:        sku.Volume,
			Image:         utils.AddDomainToURL(sku.Image), // 添加域名前缀
		}
		result = append(result, dtoSku)
	}

	return result, nil
}

// loadSkuSpecs 加载SKU规格关联
func (s *MultiSpecService) loadSkuSpecs(skuId int64) (map[string]dto.SpecValue, error) {
	relationModel := &shop.ShopGoodsSkuSpec{}
	relations, err := relationModel.FindBySkuId(skuId)
	if err != nil {
		return nil, err
	}

	specs := make(map[string]dto.SpecValue)
	for _, relation := range relations {
		// 获取规格组名称
		specGroupModel := &shop.ShopGoodsSpecGroup{Id: relation.SpecGroupId}
		if err := specGroupModel.Query().Where("id = ?", relation.SpecGroupId).First(specGroupModel).Error; err != nil {
			continue
		}

		// 获取规格值
		specValueModel := &shop.ShopGoodsSpecValue{Id: relation.SpecValueId}
		if err := specValueModel.Query().Where("id = ?", relation.SpecValueId).First(specValueModel).Error; err != nil {
			continue
		}

		specs[specGroupModel.Name] = dto.SpecValue{
			ID:    specValueModel.Id,
			Name:  specValueModel.Name,
			Value: specValueModel.Value,
			Sort:  specValueModel.Sort,
		}
	}

	return specs, nil
}

// getSpecTypePrefix 获取规格类型前缀
func getSpecTypePrefix(specType int8) string {
	switch specType {
	case shop.SpecGroupTypeSingle:
		return "SINGLE"
	case shop.SpecGroupTypeMulti:
		return "MULTI"
	case shop.SpecGroupTypeCombo:
		return "COMBO"
	default:
		return "UNKNOWN"
	}
}

// createComboSkus 为套餐选项创建SKU（已弃用，改为单一SKU模式）
func (s *MultiSpecService) createComboSkus(tx *gorm.DB, goodsId int64, comboGroups []dto.ComboSpecGroup, specMaps *SpecMaps) error {
	// 此方法已弃用，套餐商品现在只创建一个SKU
	// 保留此方法是为了兼容性，实际不会被调用
	return nil
}
