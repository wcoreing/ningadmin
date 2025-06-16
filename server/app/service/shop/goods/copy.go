package goods

import (
	"fmt"
	"time"
	"wnadmin/app/model/do/shop"
	"wnadmin/app/provider/db"

	"gorm.io/gorm"
)

// CopyService 商品复制服务
type CopyService struct{}

// generateUniqueSkuCode 生成唯一的SKU编码
func (s *CopyService) generateUniqueSkuCode(tx *gorm.DB, goodsId int64, index int) (string, error) {
	maxRetries := 10
	for retry := 0; retry < maxRetries; retry++ {
		now := time.Now()
		// 使用商品ID + 时间戳 + 纳秒 + 索引 + 重试次数来确保唯一性
		microsecond := now.Nanosecond() / 1000 // 转换为微秒
		skuCode := fmt.Sprintf("SKU%d%d%d%d%d", goodsId, now.Unix(), microsecond, index+1, retry)

		// 检查SKU编码是否已存在
		var count int64
		if err := tx.Model(&shop.ShopGoodsSku{}).Where("sku_code = ?", skuCode).Count(&count).Error; err != nil {
			return "", fmt.Errorf("检查SKU编码唯一性失败: %v", err)
		}

		if count == 0 {
			return skuCode, nil
		}

		// 如果重复，稍微等待一下再重试
		time.Sleep(time.Millisecond * time.Duration(retry+1))
	}

	return "", fmt.Errorf("生成唯一SKU编码失败，已重试%d次", maxRetries)
}

// Execute 执行商品复制
func (s *CopyService) Execute(goodsId int64, count int) error {
	// 查询原商品信息
	originalGoods := &shop.GoodsDO{}
	if err := db.WnAdmin().First(originalGoods, goodsId).Error; err != nil {
		return fmt.Errorf("商品不存在")
	}

	// 执行复制操作
	for i := 0; i < count; i++ {
		if err := s.copyGoods(originalGoods); err != nil {
			return fmt.Errorf("复制第%d个商品失败: %v", i+1, err)
		}
	}

	return nil
}

// copyGoods 复制单个商品
func (s *CopyService) copyGoods(originalGoods *shop.GoodsDO) error {
	return db.WnAdmin().Transaction(func(tx *gorm.DB) error {
		// 1. 复制商品基本信息
		newGoods := &shop.GoodsDO{
			Name:        fmt.Sprintf("%s-复制", originalGoods.Name),
			CategoryId:  originalGoods.CategoryId,
			Brand:       originalGoods.Brand,
			Supplier:    originalGoods.Supplier,
			Type:        originalGoods.Type,
			Unit:        originalGoods.Unit,
			Sort:        originalGoods.Sort,
			Status:      0, // 复制的商品默认为下架状态
			Description: originalGoods.Description,
			SpecType:    originalGoods.SpecType,
			CreateTime:  time.Now(),
			UpdateTime:  time.Now(),
		}

		if err := tx.Create(newGoods).Error; err != nil {
			return fmt.Errorf("复制商品基本信息失败: %v", err)
		}

		// 2. 复制商品详情
		var originalDetail shop.ShopGoodsDetail
		if err := tx.Where("goods_id = ?", originalGoods.Id).First(&originalDetail).Error; err == nil {
			newDetail := &shop.ShopGoodsDetail{
				GoodsId:    int(newGoods.Id),
				DetailDesc: originalDetail.DetailDesc,
				CreateTime: time.Now(),
				UpdateTime: time.Now(),
			}
			if err := tx.Create(newDetail).Error; err != nil {
				return fmt.Errorf("复制商品详情失败: %v", err)
			}
		}

		// 3. 复制商品图片
		var originalImages []shop.ShopGoodsImage
		if err := tx.Where("goods_id = ?", originalGoods.Id).Find(&originalImages).Error; err == nil {
			for _, originalImage := range originalImages {
				newImage := &shop.ShopGoodsImage{
					GoodsId:    newGoods.Id,
					Position:   originalImage.Position,
					Url:        originalImage.Url,
					Sort:       originalImage.Sort,
					CreateTime: time.Now(),
					UpdateTime: time.Now(),
				}
				if err := tx.Create(newImage).Error; err != nil {
					return fmt.Errorf("复制商品图片失败: %v", err)
				}
			}
		}

		// 4. 复制标签关联
		var originalTagRels []shop.ShopGoodsTagRel
		if err := tx.Where("goods_id = ?", originalGoods.Id).Find(&originalTagRels).Error; err == nil {
			for _, originalTagRel := range originalTagRels {
				newTagRel := &shop.ShopGoodsTagRel{
					GoodsId:    newGoods.Id,
					TagName:    originalTagRel.TagName,
					CreateTime: time.Now(),
				}
				if err := tx.Create(newTagRel).Error; err != nil {
					return fmt.Errorf("复制商品标签失败: %v", err)
				}
			}
		}

		// 5. 复制规格信息（先检查是否有规格数据）
		var specGroupIdMap, specValueIdMap map[int64]int64
		var specGroupCount int64
		if err := tx.Model(&shop.ShopGoodsSpecGroup{}).Where("goods_id = ?", originalGoods.Id).Count(&specGroupCount).Error; err == nil && specGroupCount > 0 {
			var err error
			specGroupIdMap, specValueIdMap, err = s.copySpecGroups(tx, originalGoods.Id, newGoods.Id)
			if err != nil {
				return err
			}
		}

		// 6. 复制SKU和SKU规格关联
		if err := s.copySkusAndSpecs(tx, originalGoods.Id, newGoods.Id, specGroupIdMap, specValueIdMap); err != nil {
			return err
		}

		return nil
	})
}

// copySpecGroups 复制规格组和规格值
func (s *CopyService) copySpecGroups(tx *gorm.DB, originalGoodsId, newGoodsId int64) (map[int64]int64, map[int64]int64, error) {
	// 查询原商品的规格组
	var originalSpecGroups []shop.ShopGoodsSpecGroup
	if err := tx.Where("goods_id = ?", originalGoodsId).Find(&originalSpecGroups).Error; err != nil {
		return nil, nil, fmt.Errorf("查询原规格组失败: %v", err)
	}

	// 用于存储原规格组ID到新规格组ID的映射
	specGroupIdMap := make(map[int64]int64)
	specValueIdMap := make(map[int64]int64)

	// 复制规格组
	for _, originalGroup := range originalSpecGroups {
		newGroup := &shop.ShopGoodsSpecGroup{
			GoodsId:    newGoodsId,
			Name:       originalGroup.Name,
			Sort:       originalGroup.Sort,
			Status:     originalGroup.Status,
			CreateTime: time.Now(),
			UpdateTime: time.Now(),
		}

		if err := tx.Create(newGroup).Error; err != nil {
			return nil, nil, fmt.Errorf("复制规格组失败: %v", err)
		}

		// 记录ID映射
		specGroupIdMap[originalGroup.Id] = newGroup.Id

		// 查询并复制该规格组下的规格值
		var originalSpecValues []shop.ShopGoodsSpecValue
		if err := tx.Where("goods_id = ? AND spec_group_id = ?", originalGoodsId, originalGroup.Id).Find(&originalSpecValues).Error; err != nil {
			return nil, nil, fmt.Errorf("查询原规格值失败: %v", err)
		}

		for _, originalValue := range originalSpecValues {
			newValue := &shop.ShopGoodsSpecValue{
				GoodsId:     newGoodsId,
				SpecGroupId: newGroup.Id,
				Name:        originalValue.Name,
				Value:       originalValue.Value,
				Sort:        originalValue.Sort,
				Status:      originalValue.Status,
				CreateTime:  time.Now(),
				UpdateTime:  time.Now(),
			}

			if err := tx.Create(newValue).Error; err != nil {
				return nil, nil, fmt.Errorf("复制规格值失败: %v", err)
			}

			// 记录ID映射
			specValueIdMap[originalValue.Id] = newValue.Id
		}
	}

	return specGroupIdMap, specValueIdMap, nil
}

// copySkusAndSpecs 复制SKU和SKU规格关联
func (s *CopyService) copySkusAndSpecs(tx *gorm.DB, originalGoodsId, newGoodsId int64, specGroupIdMap, specValueIdMap map[int64]int64) error {
	// 查询原商品的SKU
	var originalSkus []shop.ShopGoodsSku
	if err := tx.Where("goods_id = ?", originalGoodsId).Find(&originalSkus).Error; err != nil {
		return fmt.Errorf("查询原SKU失败: %v", err)
	}

	// 规格映射关系已通过参数传入

	// 复制SKU
	for i, originalSku := range originalSkus {
		// 生成唯一的SKU编码
		skuCode, err := s.generateUniqueSkuCode(tx, newGoodsId, i)
		if err != nil {
			return fmt.Errorf("生成SKU编码失败: %v", err)
		}

		newSku := &shop.ShopGoodsSku{
			GoodsId:       newGoodsId,
			SpecType:      originalSku.SpecType, // 保持原有的规格类型
			SkuCode:       skuCode,
			Price:         originalSku.Price,
			OriginalPrice: originalSku.OriginalPrice,
			CostPrice:     originalSku.CostPrice,
			Stock:         originalSku.Stock,
			Weight:        originalSku.Weight,
			Volume:        originalSku.Volume,
			Status:        originalSku.Status,
			Sort:          originalSku.Sort,
			Image:         originalSku.Image,
			CreateTime:    time.Now(),
			UpdateTime:    time.Now(),
		}

		if err := tx.Create(newSku).Error; err != nil {
			return fmt.Errorf("复制SKU失败: %v", err)
		}

		// 复制SKU规格关联（如果是多规格商品）
		if len(specGroupIdMap) > 0 && len(specValueIdMap) > 0 {
			// 先尝试复制原有的SKU规格关联
			if err := s.copySkuSpecs(tx, originalSku.Id, newSku.Id, specGroupIdMap, specValueIdMap); err != nil {
				return fmt.Errorf("复制SKU规格关联失败: %v", err)
			}
		}
	}

	return nil
}

// copySkuSpecs 复制SKU规格关联
func (s *CopyService) copySkuSpecs(tx *gorm.DB, originalSkuId, newSkuId int64, specGroupIdMap, specValueIdMap map[int64]int64) error {
	// 查询原SKU的规格关联
	var originalSkuSpecs []shop.ShopGoodsSkuSpec
	if err := tx.Where("sku_id = ?", originalSkuId).Find(&originalSkuSpecs).Error; err != nil {
		return fmt.Errorf("查询原SKU规格关联失败: %v", err)
	}

	// 复制SKU规格关联
	for _, originalSkuSpec := range originalSkuSpecs {
		// 获取新的规格组ID和规格值ID
		newSpecGroupId, groupExists := specGroupIdMap[originalSkuSpec.SpecGroupId]
		newSpecValueId, valueExists := specValueIdMap[originalSkuSpec.SpecValueId]

		if !groupExists || !valueExists {
			// 如果找不到对应的映射，跳过这个关联
			continue
		}

		newSkuSpec := &shop.ShopGoodsSkuSpec{
			SkuId:       newSkuId,
			SpecGroupId: newSpecGroupId,
			SpecValueId: newSpecValueId,
			CreateTime:  time.Now(),
		}

		if err := tx.Create(newSkuSpec).Error; err != nil {
			return fmt.Errorf("创建SKU规格关联失败: %v", err)
		}

	}

	return nil
}
