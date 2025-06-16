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

// CreateService 商品创建服务
type CreateService struct{}

// Execute 执行商品创建
func (s *CreateService) Execute(req *dto.GoodsUpdateReq) error {
	// 开启事务
	tx := db.WnAdmin().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 1. 创建商品基本信息
	goods, err := s.createGoodsBasicInfo(tx, req)
	if err != nil {
		tx.Rollback()
		return err
	}

	// 2. 创建商品详情
	if err := s.createGoodsDetail(tx, goods.Id, req.DetailDesc); err != nil {
		tx.Rollback()
		return err
	}

	// 3. 处理商品图片
	if err := s.handleGoodsImages(tx, goods.Id, req); err != nil {
		tx.Rollback()
		return err
	}

	// 4. 处理商品标签
	if err := s.handleGoodsTags(tx, goods.Id, req.TagNames); err != nil {
		tx.Rollback()
		return err
	}

	// 5. 处理商品规格和SKU
	if err := s.handleGoodsSpecs(tx, goods, req); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// createGoodsBasicInfo 创建商品基本信息
func (s *CreateService) createGoodsBasicInfo(tx *gorm.DB, req *dto.GoodsUpdateReq) (*shop.GoodsDO, error) {
	goods := &shop.GoodsDO{
		Name:        req.GoodsName,
		CategoryId:  req.CategoryId,
		Brand:       req.Brand,
		Supplier:    req.Supplier,
		Type:        req.Type, // 商品类型：1-实物商品 2-虚拟商品 3-套餐商品
		Unit:        req.Unit, // 商品单位
		Sort:        req.Sort,
		Description: req.Description,
		SpecType:    int8(req.SpecType),
		Status:      int8(req.Status),
		CreateTime:  time.Now(),
		UpdateTime:  time.Now(),
	}

	if err := tx.Create(goods).Error; err != nil {
		return nil, fmt.Errorf("创建商品失败: %v", err)
	}

	return goods, nil
}

// createGoodsDetail 创建商品详情
func (s *CreateService) createGoodsDetail(tx *gorm.DB, goodsId int64, detailDesc string) error {
	detail := &shop.ShopGoodsDetail{
		GoodsId:    int(goodsId),
		DetailDesc: detailDesc,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}

	if err := tx.Create(detail).Error; err != nil {
		return fmt.Errorf("创建商品详情失败: %v", err)
	}

	return nil
}

// handleGoodsImages 处理商品图片
func (s *CreateService) handleGoodsImages(tx *gorm.DB, goodsId int64, req *dto.GoodsUpdateReq) error {
	// 添加商品主图
	if req.CoverImg != "" {
		// 去除域名，只保存相对路径
		cleanUrl := utils.RemoveDomainFromURL(req.CoverImg)
		if err := s.createGoodsImage(tx, goodsId, cleanUrl, shop.ImagePositionMain, 0); err != nil {
			return err
		}
	}

	// 添加商品轮播图
	for i, imageItem := range req.Images {
		if imageItem != nil && imageItem.Url != "" {
			// 去除域名，只保存相对路径
			cleanUrl := utils.RemoveDomainFromURL(imageItem.Url)
			if err := s.createGoodsImage(tx, goodsId, cleanUrl, shop.ImagePositionNormal, i+1); err != nil {
				return err
			}
		}
	}

	return nil
}

// createGoodsImage 创建单个商品图片记录
func (s *CreateService) createGoodsImage(tx *gorm.DB, goodsId int64, imageUrl string, position uint8, sort int) error {
	image := &shop.ShopGoodsImage{
		GoodsId:    goodsId,
		Position:   position,
		Url:        imageUrl,
		Sort:       sort,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}

	if err := tx.Create(image).Error; err != nil {
		return fmt.Errorf("添加商品图片失败: %v", err)
	}

	return nil
}

// handleGoodsTags 处理商品标签关联
func (s *CreateService) handleGoodsTags(tx *gorm.DB, goodsId int64, tagNames []string) error {
	if len(tagNames) == 0 {
		return nil
	}

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

	return nil
}

// ensureTagExists 确保标签存在，如果不存在则创建
func (s *CreateService) ensureTagExists(tx *gorm.DB, tagName string) error {
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

// handleGoodsSpecs 处理商品规格和SKU - 支持多规格数据同时存储
func (s *CreateService) handleGoodsSpecs(tx *gorm.DB, goods *shop.GoodsDO, req *dto.GoodsUpdateReq) error {
	fmt.Printf("CreateService - 开始处理规格创建\n")

	// 使用统一的多规格服务
	multiSpecService := &MultiSpecService{}
	if err := multiSpecService.SaveAllSpecData(tx, goods.Id, req); err != nil {
		return fmt.Errorf("保存多规格数据失败: %v", err)
	}

	fmt.Printf("CreateService - 规格创建完成\n")
	return nil
}

// handleMultiSpecs 处理多规格商品
func (s *CreateService) handleMultiSpecs(tx *gorm.DB, goods *shop.GoodsDO, req *dto.GoodsUpdateReq) error {
	if len(req.SpecGroups) == 0 {
		return fmt.Errorf("多规格商品必须有规格组")
	}

	// 1. 创建规格组和规格值
	specMaps, err := s.createSpecGroupsAndValues(tx, goods.Id, req.SpecGroups)
	if err != nil {
		return err
	}

	// 2. 创建SKU
	if len(req.SkuList) > 0 {
		return s.createMultiSpecSkus(tx, goods.Id, req.SkuList, specMaps)
	}

	return nil
}

// createSpecGroupsAndValues 创建规格组和规格值
func (s *CreateService) createSpecGroupsAndValues(tx *gorm.DB, goodsId int64, specGroups []dto.SpecGroup) (*SpecMaps, error) {
	specGroupMap := make(map[string]int64)
	specValueMap := make(map[string]int64)

	// 过滤有效的规格组
	validSpecGroups := s.filterValidSpecGroups(specGroups)

	// 创建规格组
	for _, group := range validSpecGroups {
		specGroup, err := s.createSpecGroup(tx, goodsId, group)
		if err != nil {
			return nil, err
		}

		specGroupMap[group.Name] = specGroup.Id

		// 创建规格值
		for _, value := range group.Values {
			specValue, err := s.createSpecValue(tx, goodsId, specGroup.Id, value)
			if err != nil {
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

// filterValidSpecGroups 过滤有效的规格组
func (s *CreateService) filterValidSpecGroups(specGroups []dto.SpecGroup) []dto.SpecGroup {
	validSpecGroups := make([]dto.SpecGroup, 0)
	for _, group := range specGroups {
		if group.Name != "" && len(group.Values) > 0 {
			validSpecGroups = append(validSpecGroups, group)
		}
	}
	return validSpecGroups
}

// createSpecGroup 创建规格组
func (s *CreateService) createSpecGroup(tx *gorm.DB, goodsId int64, group dto.SpecGroup) (*shop.ShopGoodsSpecGroup, error) {
	specGroup := &shop.ShopGoodsSpecGroup{
		GoodsId:    goodsId,
		Name:       group.Name,
		Sort:       group.Sort,
		Status:     1,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}

	if err := tx.Create(specGroup).Error; err != nil {
		return nil, fmt.Errorf("创建规格组失败: %v", err)
	}

	return specGroup, nil
}

// createSpecValue 创建规格值
func (s *CreateService) createSpecValue(tx *gorm.DB, goodsId, specGroupId int64, value dto.SpecValue) (*shop.ShopGoodsSpecValue, error) {
	specValue := &shop.ShopGoodsSpecValue{
		GoodsId:     goodsId,
		SpecGroupId: specGroupId,
		Name:        value.Name,
		Value:       value.Value,
		Sort:        value.Sort,
		Status:      1,
		CreateTime:  time.Now(),
		UpdateTime:  time.Now(),
	}

	if err := tx.Create(specValue).Error; err != nil {
		return nil, fmt.Errorf("创建规格值失败: %v", err)
	}

	return specValue, nil
}

// createMultiSpecSkus 创建多规格SKU
func (s *CreateService) createMultiSpecSkus(tx *gorm.DB, goodsId int64, skuList []dto.Sku, specMaps *SpecMaps) error {
	for _, sku := range skuList {
		// 创建SKU基本信息
		skuData, err := s.createSku(tx, goodsId, sku)
		if err != nil {
			return err
		}

		// 创建SKU规格关联
		if err := s.createSkuSpecRelations(tx, skuData.Id, sku.Specs, specMaps); err != nil {
			return err
		}
	}

	return nil
}

// createSku 创建SKU基本信息
func (s *CreateService) createSku(tx *gorm.DB, goodsId int64, sku dto.Sku) (*shop.ShopGoodsSku, error) {
	return s.createSkuWithSpecType(tx, goodsId, sku, 2) // 默认多规格类型
}

// createSkuWithSpecType 创建指定规格类型的SKU
func (s *CreateService) createSkuWithSpecType(tx *gorm.DB, goodsId int64, sku dto.Sku, specType int8) (*shop.ShopGoodsSku, error) {
	skuCode := s.generateSkuCode(goodsId)

	// 去除图片URL中的域名
	cleanImageUrl := utils.RemoveDomainFromURL(sku.Image)

	skuData := &shop.ShopGoodsSku{
		GoodsId:       goodsId,
		SpecType:      specType, // 使用传入的规格类型
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
		return nil, fmt.Errorf("保存SKU失败: %v", err)
	}

	return skuData, nil
}

// createSkuSpecRelations 创建SKU规格关联
func (s *CreateService) createSkuSpecRelations(tx *gorm.DB, skuId int64, specs map[string]dto.SpecValue, specMaps *SpecMaps) error {
	for specGroupName, specValue := range specs {
		specGroupId := specMaps.GroupMap[specGroupName]
		specValueKey := fmt.Sprintf("%s:%s", specGroupName, specValue.Value)
		specValueId := specMaps.ValueMap[specValueKey]

		if specGroupId > 0 && specValueId > 0 {
			skuSpec := &shop.ShopGoodsSkuSpec{
				SkuId:       skuId,
				SpecGroupId: specGroupId,
				SpecValueId: specValueId,
				CreateTime:  time.Now(),
			}

			if err := tx.Create(skuSpec).Error; err != nil {
				return fmt.Errorf("保存SKU规格关联失败: %v", err)
			}
		}
	}

	return nil
}

// handleSingleSpec 处理单规格商品
func (s *CreateService) handleSingleSpec(tx *gorm.DB, goods *shop.GoodsDO, req *dto.GoodsUpdateReq) error {
	// 获取单规格商品的价格库存信息
	skuInfo := s.getSingleSpecSkuInfo(req)

	// 创建默认SKU，使用单规格类型
	_, err := s.createSkuWithSpecType(tx, goods.Id, skuInfo, 1)
	return err
}

// getSingleSpecSkuInfo 获取单规格商品的SKU信息
func (s *CreateService) getSingleSpecSkuInfo(req *dto.GoodsUpdateReq) dto.Sku {
	// 如果有SKU列表，使用第一个SKU的信息
	if len(req.SkuList) > 0 {
		sku := req.SkuList[0]
		// 去除图片URL中的域名
		sku.Image = utils.RemoveDomainFromURL(sku.Image)
		return sku
	}

	// 如果没有SKU列表，使用前端传递的单规格字段
	return dto.Sku{
		Price:         req.Price,
		OriginalPrice: req.OriginalPrice,
		CostPrice:     req.CostPrice,
		Stock:         req.Stock,
		Weight:        req.Weight,
		Volume:        req.Volume,
		Image:         utils.RemoveDomainFromURL(req.CoverImg), // 使用封面图作为SKU图片，去除域名
	}
}

// generateSkuCode 生成SKU编码
func (s *CreateService) generateSkuCode(goodsId int64) string {
	return fmt.Sprintf("SKU%d%d%d", goodsId, time.Now().Unix(), rand.Int63n(1000000))
}

// handleComboSpecs 处理套餐规格商品
func (s *CreateService) handleComboSpecs(tx *gorm.DB, goods *shop.GoodsDO, req *dto.GoodsUpdateReq) error {
	fmt.Printf("创建套餐规格商品\n")

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

	// 套餐商品创建单一默认SKU（类似单规格商品）
	skuCode := fmt.Sprintf("COMBO%d%d%d", goods.Id, time.Now().Unix(), rand.Int63n(1000000))

	// 套餐商品的价格优先从ComboSpecData.Price获取，其次从req.Price获取
	var basePrice float64 = 0
	if req.ComboSpecData != nil && req.ComboSpecData.Price > 0 {
		basePrice = req.ComboSpecData.Price
	} else {
		basePrice = req.Price
	}
	fmt.Printf("套餐商品创建 - 获取价格: req.Price=%f, ComboSpecData.Price=%f, basePrice=%f\n", req.Price,
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
