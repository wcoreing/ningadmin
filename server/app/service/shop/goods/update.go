package goods

import (
	"fmt"
	"math/rand"
	"time"
	"wnadmin/app/common/utils"
	"wnadmin/app/model/do/shop"
	"wnadmin/app/model/dto"
	"wnadmin/app/provider/db"

	"gorm.io/gorm"
)

// UpdateService 商品更新服务
type UpdateService struct{}

// Execute 执行商品更新
func (s *UpdateService) Execute(req *dto.GoodsUpdateReq) error {
	// 设置事务超时时间为30秒
	return db.WnAdmin().Transaction(func(tx *gorm.DB) error {
		// 1. 查询商品
		goods := &shop.GoodsDO{Id: req.GoodsId}
		if err := tx.First(goods).Error; err != nil {
			return fmt.Errorf("商品不存在")
		}

		// 2. 更新商品基本信息
		if err := s.updateGoodsBasicInfo(tx, goods, req); err != nil {
			return err
		}

		// 3. 更新商品详情
		if err := s.updateGoodsDetail(tx, goods.Id, req.DetailDesc); err != nil {
			return err
		}

		// 4. 处理规格信息
		if err := s.handleSpecsUpdate(tx, goods, req); err != nil {
			return err
		}

		// 5. 处理图片
		if err := s.handleImagesUpdate(tx, goods.Id, req); err != nil {
			return err
		}

		// 6. 更新标签关联
		if err := s.updateTagRelations(tx, goods.Id, req.TagNames); err != nil {
			return err
		}

		return nil // 让事务自动提交
	})
}

// updateGoodsBasicInfo 更新商品基本信息
func (s *UpdateService) updateGoodsBasicInfo(tx *gorm.DB, goods *shop.GoodsDO, req *dto.GoodsUpdateReq) error {
	// 更新商品基本信息
	goods.Name = req.GoodsName
	goods.CategoryId = req.CategoryId
	goods.Brand = req.Brand
	goods.Supplier = req.Supplier
	goods.Type = req.Type // 商品类型：1-实物商品 2-虚拟商品 3-套餐商品
	goods.Unit = req.Unit // 商品单位
	goods.Sort = req.Sort
	goods.Description = req.Description
	goods.SpecType = int8(req.SpecType)
	goods.Status = int8(req.Status)

	if err := tx.Save(goods).Error; err != nil {
		return fmt.Errorf("更新商品失败: %v", err)
	}

	return nil
}

// updateGoodsDetail 更新商品详情
func (s *UpdateService) updateGoodsDetail(tx *gorm.DB, goodsId int64, detailDesc string) error {
	var detail shop.ShopGoodsDetail
	if err := tx.Where("goods_id = ?", goodsId).First(&detail).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			detail = shop.ShopGoodsDetail{
				GoodsId: int(goodsId),
			}
		} else {
			return fmt.Errorf("查询商品详情失败: %v", err)
		}
	}
	detail.DetailDesc = detailDesc
	if err := tx.Save(&detail).Error; err != nil {
		return fmt.Errorf("保存商品详情失败: %v", err)
	}

	return nil
}

// handleSpecsUpdate 处理规格信息更新 - 支持多规格数据同时存储
func (s *UpdateService) handleSpecsUpdate(tx *gorm.DB, goods *shop.GoodsDO, req *dto.GoodsUpdateReq) error {
	fmt.Printf("UpdateService - 开始处理规格更新\n")

	// 使用统一的多规格服务
	multiSpecService := &MultiSpecService{}
	if err := multiSpecService.SaveAllSpecData(tx, goods.Id, req); err != nil {
		return fmt.Errorf("保存多规格数据失败: %v", err)
	}

	fmt.Printf("UpdateService - 规格更新完成\n")
	return nil
}

// handleMultiSpecsUpdate 处理多规格商品更新
func (s *UpdateService) handleMultiSpecsUpdate(tx *gorm.DB, goods *shop.GoodsDO, req *dto.GoodsUpdateReq) error {
	// 删除原有规格组和规格值
	if err := s.deleteOldSpecs(tx, goods.Id); err != nil {
		return err
	}

	// 删除原有SKU及其关联
	if err := s.deleteOldSkus(tx, goods.Id); err != nil {
		return err
	}

	// 保存新的规格组和规格值
	specGroupMap := make(map[string]int64) // 规格组名称到ID的映射
	specValueMap := make(map[string]int64) // 规格值到ID的映射

	// 过滤有效的规格组（名称不为空且有规格值）
	validSpecGroups := make([]dto.SpecGroup, 0)
	for _, group := range req.SpecGroups {
		if group.Name != "" && len(group.Values) > 0 {
			validSpecGroups = append(validSpecGroups, group)
		}
	}

	// 保存规格组
	for _, group := range validSpecGroups {
		specGroup := &shop.ShopGoodsSpecGroup{
			GoodsId:    goods.Id,
			Name:       group.Name,
			Sort:       group.Sort,
			Status:     1,
			CreateTime: time.Now(),
			UpdateTime: time.Now(),
		}

		if err := tx.Create(specGroup).Error; err != nil {
			return fmt.Errorf("保存规格组失败: %v", err)
		}

		specGroupMap[group.Name] = specGroup.Id

		// 保存规格值
		for _, value := range group.Values {
			specValue := &shop.ShopGoodsSpecValue{
				GoodsId:     goods.Id,
				SpecGroupId: specGroup.Id,
				Name:        value.Name,
				Value:       value.Value,
				Sort:        value.Sort,
				Status:      1,
				CreateTime:  time.Now(),
				UpdateTime:  time.Now(),
			}

			if err := tx.Create(specValue).Error; err != nil {
				return fmt.Errorf("保存规格值失败: %v", err)
			}

			// 使用组合键作为映射key
			specValueKey := fmt.Sprintf("%s:%s", group.Name, value.Value)
			specValueMap[specValueKey] = specValue.Id
		}
	}

	// 保存SKU列表
	for _, sku := range req.SkuList {
		skuCode := fmt.Sprintf("SKU%d%d%d", goods.Id, time.Now().Unix(), rand.Int63n(1000000))
		// 去除图片URL中的域名
		cleanImageUrl := utils.RemoveDomainFromURL(sku.Image)

		skuData := &shop.ShopGoodsSku{
			GoodsId:       goods.Id,
			SpecType:      2, // 多规格类型
			SkuCode:       skuCode,
			Price:         sku.Price,
			OriginalPrice: sku.OriginalPrice,
			CostPrice:     sku.CostPrice,
			Stock:         sku.Stock,
			Weight:        sku.Weight,
			Volume:        sku.Volume,
			Status:        1,
			Sort:          0,
			Image:         cleanImageUrl,
			CreateTime:    time.Now(),
			UpdateTime:    time.Now(),
		}

		if err := tx.Create(skuData).Error; err != nil {
			return fmt.Errorf("保存SKU失败: %v", err)
		}

		// 保存SKU规格关联
		for specGroupName, specValue := range sku.Specs {
			specGroupId := specGroupMap[specGroupName]
			specValueKey := fmt.Sprintf("%s:%s", specGroupName, specValue.Value)
			specValueId := specValueMap[specValueKey]

			if specGroupId > 0 && specValueId > 0 {
				skuSpec := &shop.ShopGoodsSkuSpec{
					SkuId:       skuData.Id,
					SpecGroupId: specGroupId,
					SpecValueId: specValueId,
					CreateTime:  time.Now(),
				}

				if err := tx.Create(skuSpec).Error; err != nil {
					return fmt.Errorf("保存SKU规格关联失败: %v", err)
				}
			}
		}
	}

	return nil
}

// handleSingleSpecUpdate 处理单规格商品更新
func (s *UpdateService) handleSingleSpecUpdate(tx *gorm.DB, goods *shop.GoodsDO, req *dto.GoodsUpdateReq) error {
	fmt.Printf("处理单规格商品\n")

	// 删除原有规格组和规格值
	if err := s.deleteOldSpecs(tx, goods.Id); err != nil {
		return err
	}

	// 删除原有SKU及其关联
	if err := s.deleteOldSkus(tx, goods.Id); err != nil {
		return err
	}

	// 从SKU列表中获取价格库存信息，如果没有SKU列表则使用前端传递的字段
	var price, originalPrice, costPrice, weight, volume float64
	var stock int
	var image string

	if len(req.SkuList) > 0 {
		// 如果有SKU列表，使用第一个SKU的信息
		firstSku := req.SkuList[0]
		price = firstSku.Price
		originalPrice = firstSku.OriginalPrice
		costPrice = firstSku.CostPrice
		stock = firstSku.Stock
		weight = firstSku.Weight
		volume = firstSku.Volume
		image = utils.RemoveDomainFromURL(firstSku.Image) // 去除域名
	} else {
		// 如果没有SKU列表，使用前端传递的单规格字段
		price = req.Price
		originalPrice = req.OriginalPrice
		costPrice = req.CostPrice
		stock = req.Stock
		weight = req.Weight
		volume = req.Volume
		image = utils.RemoveDomainFromURL(req.CoverImg) // 使用封面图作为SKU图片，去除域名
	}

	sku := &shop.ShopGoodsSku{
		GoodsId:       goods.Id,
		SpecType:      1, // 单规格类型
		SkuCode:       fmt.Sprintf("SKU%d%d%d", goods.Id, time.Now().Unix(), rand.Int63n(1000000)),
		Price:         price,
		OriginalPrice: originalPrice,
		CostPrice:     costPrice,
		Stock:         stock,
		Weight:        weight,
		Volume:        volume,
		Status:        1,
		Sort:          0,
		Image:         image,
		CreateTime:    time.Now(),
		UpdateTime:    time.Now(),
	}

	if err := tx.Create(sku).Error; err != nil {
		return fmt.Errorf("保存SKU失败: %v", err)
	}

	return nil
}

// deleteOldSpecs 删除原有规格组和规格值
func (s *UpdateService) deleteOldSpecs(tx *gorm.DB, goodsId int64) error {
	if err := tx.Where("goods_id = ?", goodsId).Delete(&shop.ShopGoodsSpecGroup{}).Error; err != nil {
		return fmt.Errorf("删除原有规格组失败: %v", err)
	}
	if err := tx.Where("goods_id = ?", goodsId).Delete(&shop.ShopGoodsSpecValue{}).Error; err != nil {
		return fmt.Errorf("删除原有规格值失败: %v", err)
	}
	return nil
}

// deleteOldSkus 删除原有SKU及其关联
func (s *UpdateService) deleteOldSkus(tx *gorm.DB, goodsId int64) error {
	var oldSkus []shop.ShopGoodsSku
	if err := tx.Where("goods_id = ?", goodsId).Find(&oldSkus).Error; err != nil {
		return fmt.Errorf("查询原有SKU失败: %v", err)
	}
	for _, sku := range oldSkus {
		if err := tx.Where("sku_id = ?", sku.Id).Delete(&shop.ShopGoodsSkuSpec{}).Error; err != nil {
			return fmt.Errorf("删除SKU规格关联失败: %v", err)
		}
	}
	if err := tx.Where("goods_id = ?", goodsId).Delete(&shop.ShopGoodsSku{}).Error; err != nil {
		return fmt.Errorf("删除原有SKU失败: %v", err)
	}
	return nil
}

// handleImagesUpdate 处理图片更新
func (s *UpdateService) handleImagesUpdate(tx *gorm.DB, goodsId int64, req *dto.GoodsUpdateReq) error {
	// 处理主图
	// 先删除所有主图
	if err := tx.Where("goods_id = ? AND position = ?", goodsId, shop.ImagePositionMain).Delete(&shop.ShopGoodsImage{}).Error; err != nil {
		return fmt.Errorf("删除旧主图失败: %v", err)
	}

	// 只有当 CoverImg 不为空时才添加新的主图
	if req.CoverImg != "" {
		// 去除域名，只保存相对路径
		cleanUrl := utils.RemoveDomainFromURL(req.CoverImg)
		mainImage := &shop.ShopGoodsImage{
			GoodsId:    goodsId,
			Position:   shop.ImagePositionMain,
			Url:        cleanUrl,
			Sort:       0,
			IsDelete:   0,
			CreateTime: time.Now(),
			UpdateTime: time.Now(),
		}
		if err := tx.Create(mainImage).Error; err != nil {
			return fmt.Errorf("添加主图失败: %v", err)
		}
	}

	// 处理商品图片列表
	// 先删除所有非主图的图片
	if err := tx.Where("goods_id = ? AND position != ?", goodsId, shop.ImagePositionMain).Delete(&shop.ShopGoodsImage{}).Error; err != nil {
		return fmt.Errorf("删除旧图片失败: %v", err)
	}

	// 添加新的图片列表，确保不添加主图
	if len(req.Images) > 0 {
		for _, img := range req.Images {
			// 跳过主图位置的图片
			if img.Position == shop.ImagePositionMain {
				continue
			}
			// 去除域名，只保存相对路径
			cleanUrl := utils.RemoveDomainFromURL(img.Url)
			image := &shop.ShopGoodsImage{
				GoodsId:    goodsId,
				Position:   img.Position,
				Url:        cleanUrl,
				Sort:       img.Sort,
				IsDelete:   0,
				CreateTime: time.Now(),
				UpdateTime: time.Now(),
			}
			if err := tx.Create(image).Error; err != nil {
				return fmt.Errorf("添加图片失败: %v", err)
			}
		}
	}

	return nil
}

// updateTagRelations 更新标签关联
func (s *UpdateService) updateTagRelations(tx *gorm.DB, goodsId int64, tagNames []string) error {
	// 删除旧的标签关联
	if err := tx.Where("goods_id = ?", goodsId).Delete(&shop.ShopGoodsTagRel{}).Error; err != nil {
		return fmt.Errorf("删除旧标签失败: %v", err)
	}

	// 添加新的标签关联
	if len(tagNames) > 0 {
		for _, tagName := range tagNames {
			if tagName == "" {
				continue // 跳过空标签名
			}

			// 检查标签是否已存在，如果不存在则创建
			if err := s.ensureTagExists(tx, tagName); err != nil {
				return fmt.Errorf("确保标签存在失败: %v", err)
			}

			rel := &shop.ShopGoodsTagRel{
				GoodsId:    goodsId,
				TagName:    tagName,
				CreateTime: time.Now(),
			}
			if err := tx.Create(rel).Error; err != nil {
				return fmt.Errorf("添加标签关联失败: %v", err)
			}
		}
	}

	return nil
}

// ensureTagExists 确保标签存在，如果不存在则创建
func (s *UpdateService) ensureTagExists(tx *gorm.DB, tagName string) error {
	var tag shop.ShopGoodsTag
	err := tx.Where("name = ? AND is_delete = 0", tagName).First(&tag).Error
	if err == nil {
		return nil // 标签已存在
	}

	if err != gorm.ErrRecordNotFound {
		return err // 查询出错
	}

	// 标签不存在，创建新标签
	newTag := &shop.ShopGoodsTag{
		Name:       tagName,
		Sort:       0,
		IsDelete:   0,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}

	if err := tx.Create(newTag).Error; err != nil {
		return fmt.Errorf("创建新标签失败: %v", err)
	}

	return nil
}

// handleComboSpecsUpdate 处理套餐规格商品更新
func (s *UpdateService) handleComboSpecsUpdate(tx *gorm.DB, goods *shop.GoodsDO, req *dto.GoodsUpdateReq) error {
	fmt.Printf("处理套餐规格商品\n")

	// 删除原有规格组和规格值
	if err := s.deleteOldSpecs(tx, goods.Id); err != nil {
		return err
	}

	// 删除原有SKU及其关联
	if err := s.deleteOldSkus(tx, goods.Id); err != nil {
		return err
	}

	// 处理套餐组（comboGroups）
	if len(req.ComboGroups) > 0 {
		for _, comboGroup := range req.ComboGroups {
			// 创建套餐组（规格组）
			specGroup := &shop.ShopGoodsSpecGroup{
				GoodsId:  goods.Id,
				Name:     comboGroup.Name,
				SpecType: shop.SpecGroupTypeCombo, // 套餐规格类型
				Sort:     comboGroup.Sort,
				Status:   1,
			}

			// 设置套餐组配置
			groupConfig := &shop.ComboGroupConfig{
				Description: comboGroup.Description,
				MinSelect:   comboGroup.MinSelect,
				MaxSelect:   comboGroup.MaxSelect,
				IsRequired:  comboGroup.IsRequired,
				DisplayType: "radio", // 默认单选
			}

			if err := specGroup.SetComboConfig(groupConfig); err != nil {
				return fmt.Errorf("设置套餐组配置失败: %v", err)
			}

			if err := tx.Create(specGroup).Error; err != nil {
				return fmt.Errorf("创建套餐组失败: %v", err)
			}

			// 创建套餐选项（规格值）
			for _, option := range comboGroup.Options {
				specValue := &shop.ShopGoodsSpecValue{
					GoodsId:     goods.Id,
					SpecGroupId: specGroup.Id,
					Name:        option.Name,
					Value:       option.Value,
					SpecType:    shop.SpecValueTypeCombo, // 套餐选项类型
					Sort:        option.Sort,
					Status:      1,
				}

				// 设置套餐选项配置
				optionConfig := &shop.ComboOptionConfig{
					OptionGoodsId: &option.GoodsId,
					OptionSkuId:   option.SkuId,
					AddPrice:      option.AddPrice,
					Quantity:      option.Quantity,
					IsDefault:     option.IsDefault,
					Description:   option.Name,
				}

				if err := specValue.SetComboConfig(optionConfig); err != nil {
					return fmt.Errorf("设置套餐选项配置失败: %v", err)
				}

				if err := tx.Create(specValue).Error; err != nil {
					return fmt.Errorf("创建套餐选项失败: %v", err)
				}
			}
		}
	}

	// 套餐商品创建默认SKU
	skuCode := fmt.Sprintf("COMBO%d%d%d", goods.Id, time.Now().Unix(), rand.Int63n(1000000))

	// 套餐商品的价格优先从ComboSpecData.Price获取，其次从req.Price获取
	var basePrice float64 = 0
	if req.ComboSpecData != nil && req.ComboSpecData.Price > 0 {
		basePrice = req.ComboSpecData.Price
	} else {
		basePrice = req.Price
	}
	fmt.Printf("套餐商品更新 - 获取价格: req.Price=%f, ComboSpecData.Price=%f, basePrice=%f\n", req.Price,
		func() float64 {
			if req.ComboSpecData != nil {
				return req.ComboSpecData.Price
			}
			return 0
		}(), basePrice)

	sku := &shop.ShopGoodsSku{
		GoodsId:       goods.Id,
		SpecType:      shop.SpecGroupTypeCombo, // 设置为套餐规格类型
		SkuCode:       skuCode,
		Price:         basePrice, // 使用基础价格
		OriginalPrice: basePrice,
		CostPrice:     0,
		Stock:         999, // 套餐默认库存
		Weight:        0,
		Volume:        0,
		Status:        1,
		Sort:          0,
		Image:         req.CoverImg,
		CreateTime:    time.Now(),
		UpdateTime:    time.Now(),
	}

	if err := tx.Create(sku).Error; err != nil {
		return fmt.Errorf("创建套餐SKU失败: %v", err)
	}

	return nil
}
