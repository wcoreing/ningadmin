package shop

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"time"
	"wnadmin/app/model/do/shop"
	"wnadmin/app/model/dto"
	"wnadmin/app/provider/db"

	"gorm.io/gorm"
)

// GoodsService 商品服务
type GoodsService struct{}

// GetGoodsList 获取商品列表
func (s *GoodsService) GetGoodsList(req *dto.GoodsListReq) (*dto.PageResult, error) {
	var (
		goods  []shop.GoodsDO
		total  int64
		result = make([]dto.GoodsDTO, 0)
	)

	// 查询商品列表
	if req.Name != "" {
		// 尝试将name转换为ID
		if id, err := strconv.ParseInt(req.Name, 10, 64); err == nil {
			// 如果可以转换为数字，添加ID查询条件
			req.GoodsId = id
		}
	}

	goods, total, err := (&shop.GoodsDO{}).GetGoodsList(req)
	if err != nil {
		return nil, err
	}

	// 获取关联数据
	categoryMap := s.getCategoryMap(goods)
	supplierMap := s.getSupplierMap(goods)
	brandMap := s.getBrandMap(goods)
	tagMap := s.getTagMap(goods)

	// 组装数据
	for _, g := range goods {
		// 获取商品的标签关联
		var tagRels []shop.ShopGoodsTagRel
		var tagIds []int64
		var tagNames []string

		if err := (&shop.ShopGoodsTagRel{}).Query().
			Where("goods_id = ?", g.Id).
			Find(&tagRels).Error; err == nil {
			for _, rel := range tagRels {
				tagIds = append(tagIds, rel.TagId)
				if name, ok := tagMap[rel.TagId]; ok {
					tagNames = append(tagNames, name)
				}
			}
		}

		item := dto.GoodsDTO{
			GoodsId:       g.Id,
			CategoryId:    g.CategoryId,
			CategoryName:  categoryMap[g.CategoryId],
			SupplierId:    g.SupplierId,
			SupplierName:  supplierMap[g.SupplierId],
			BrandId:       g.BrandId,
			BrandName:     brandMap[g.BrandId],
			Name:          g.Name,
			Price:         g.Price,
			OriginalPrice: g.OriginalPrice,
			Stock:         g.Stock,
			Sales:         g.Sales,
			Status:        g.Status,
			Sort:          g.Sort,
			TagIds:        tagIds,   // 添加标签ID列表
			TagNames:      tagNames, // 添加标签名称列表
		}
		result = append(result, item)
	}

	// 查询商品图片
	for i := range result {
		var images []*shop.ShopGoodsImage
		if err := db.DB.WnAdmin.Where("goods_id = ? AND position = ?", result[i].GoodsId, shop.ImagePositionMain).
			First(&images).Error; err == nil && len(images) > 0 {
			result[i].CoverImg = images[0].Url
		}
	}

	return &dto.PageResult{
		List:     result,
		Total:    total,
		PageNum:  req.PageNum,
		PageSize: req.PageSize,
	}, nil
}

// 获取分类Map
func (s *GoodsService) getCategoryMap(goods []shop.GoodsDO) map[int64]string {
	result := make(map[int64]string)
	var categoryIds []int64
	for _, g := range goods {
		if g.CategoryId > 0 {
			categoryIds = append(categoryIds, g.CategoryId)
		}
	}
	if len(categoryIds) > 0 {
		var categories []shop.ShopCategory
		if err := (&shop.ShopCategory{}).Query().Where("id IN ?", categoryIds).Find(&categories).Error; err == nil {
			for _, c := range categories {
				result[c.ID] = c.Name
			}
		}
	}
	return result
}

// 获取供应商Map
func (s *GoodsService) getSupplierMap(goods []shop.GoodsDO) map[int64]string {
	result := make(map[int64]string)
	var supplierIds []int64
	for _, g := range goods {
		if g.SupplierId > 0 {
			supplierIds = append(supplierIds, g.SupplierId)
		}
	}
	if len(supplierIds) > 0 {
		var suppliers []shop.ShopSupplier
		if err := (&shop.ShopSupplier{}).Query().Where("id IN ?", supplierIds).Find(&suppliers).Error; err == nil {
			for _, s := range suppliers {
				result[s.ID] = s.Name
			}
		}
	}
	return result
}

// 获取品牌Map
func (s *GoodsService) getBrandMap(goods []shop.GoodsDO) map[int64]string {
	result := make(map[int64]string)
	var brandIds []int64
	for _, g := range goods {
		if g.BrandId > 0 {
			brandIds = append(brandIds, g.BrandId)
		}
	}
	if len(brandIds) > 0 {
		var brands []shop.ShopBrand
		if err := (&shop.ShopBrand{}).Query().Where("id IN ?", brandIds).Find(&brands).Error; err == nil {
			for _, b := range brands {
				result[b.ID] = b.Name
			}
		}
	}
	return result
}

// 获取标签Map
func (s *GoodsService) getTagMap(goods []shop.GoodsDO) map[int64]string {
	result := make(map[int64]string)
	var goodsIds []int64
	for _, g := range goods {
		goodsIds = append(goodsIds, g.Id)
	}

	if len(goodsIds) > 0 {
		// 查询商品标签关联
		var tagRels []shop.ShopGoodsTagRel
		if err := (&shop.ShopGoodsTagRel{}).Query().
			Where("goods_id IN ?", goodsIds).
			Find(&tagRels).Error; err == nil {

			// 获取标签信息
			var tagIds []int64
			for _, rel := range tagRels {
				tagIds = append(tagIds, rel.TagId)
			}

			var tags []shop.ShopGoodsTag
			if err := (&shop.ShopGoodsTag{}).Query().
				Where("id IN ?", tagIds).
				Find(&tags).Error; err == nil {
				for _, t := range tags {
					result[t.ID] = t.Name
				}
			}
		}
	}
	return result
}

// DeleteGoods 删除商品
func (s *GoodsService) DeleteGoods(id int64) error {
	goodsDO := &shop.GoodsDO{Id: id}
	return goodsDO.Delete()
}

// UpdateGoodsStatus 更新商品状态
func (s *GoodsService) UpdateGoodsStatus(id int64, status int8) error {
	goodsDO := &shop.GoodsDO{
		Id:     id,
		Status: status,
	}
	return goodsDO.UpdateStatus()
}

// GetGoodsInfo 获取商品详情
func (s *GoodsService) GetGoodsInfo(id int64) (*dto.GoodsInfoResp, error) {
	goods, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	// 查询商品详情
	var detail shop.ShopGoodsDetail
	if err := db.DB.WnAdmin.Where("goods_id = ?", id).First(&detail).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return nil, err
		}
	}

	fmt.Printf("商品图片数据: %+v\n", goods.Images)

	// 转换为响应对象
	resp := &dto.GoodsInfoResp{
		GoodsId:       goods.Id,
		GoodsName:     goods.Name,
		CategoryId:    goods.CategoryId,
		BrandId:       goods.BrandId,
		SupplierId:    goods.SupplierId,
		Type:          goods.Type,
		Price:         goods.Price,
		OriginalPrice: goods.OriginalPrice,
		CostPrice:     goods.CostPrice,
		Stock:         goods.Stock,
		WarnStock:     goods.WarnStock,
		Unit:          goods.Unit,
		Weight:        goods.Weight,
		Volume:        goods.Volume,
		Sort:          goods.Sort,
		Status:        strconv.FormatInt(int64(goods.Status), 10),
		Description:   goods.Description,
		DetailDesc:    detail.DetailDesc,
		SpecType:      goods.SpecType,
		TagIds:        goods.TagIds,
	}

	// 处理图片
	if len(goods.Images) > 0 {
		// 设置主图
		for _, img := range goods.Images {
			if img.Position == shop.ImagePositionMain {
				resp.CoverImg = img.Url
				break
			}
		}

		// 设置其他图片列表
		var images []*dto.GoodsImageUpdateItem
		for _, img := range goods.Images {
			if img.Position != shop.ImagePositionMain {
				images = append(images, &dto.GoodsImageUpdateItem{
					Url:      img.Url,
					Position: img.Position,
					Sort:     img.Sort,
				})
			}
		}
		resp.Images = images
	} else {
		resp.Images = make([]*dto.GoodsImageUpdateItem, 0) // 确保返回空数组而不是 nil
	}

	// 转换规格组
	if goods.SpecGroups != nil {
		resp.SpecGroups = make([]dto.SpecGroup, len(goods.SpecGroups))
		for i, group := range goods.SpecGroups {
			specGroup := dto.SpecGroup{
				ID:   group.Id,
				Name: group.Name,
				Sort: group.Sort,
			}
			// 转换规格值
			if group.Values != nil {
				specGroup.Values = make([]dto.SpecValue, len(group.Values))
				for j, value := range group.Values {
					specGroup.Values[j] = dto.SpecValue{
						ID:    value.Id,
						Name:  value.Value,
						Value: value.Value,
						Sort:  value.Sort,
					}
				}
			}
			resp.SpecGroups[i] = specGroup
		}
	}

	// 转换SKU列表
	if goods.SkuList != nil {
		resp.SkuList = make([]dto.Sku, len(goods.SkuList))
		for i, sku := range goods.SkuList {
			// 打印调试信息
			fmt.Printf("原始 SKU 数据: %+v\n", sku)
			fmt.Printf("规格 JSON: %s\n", sku.SpecsJson)

			var specs map[string]dto.SpecValue
			if err := json.Unmarshal([]byte(sku.SpecsJson), &specs); err != nil {
				fmt.Printf("解析 SKU 规格失败: %v\n", err)
				continue
			}

			resp.SkuList[i] = dto.Sku{
				ID:    sku.Id,
				Specs: specs,
				Price: sku.Price,
				Stock: sku.Stock,
				Image: sku.Image,
			}

			// 打印转换后的数据
			fmt.Printf("转换后的 SKU 数据: %+v\n", resp.SkuList[i])
		}
	}

	// 添加调试日志
	fmt.Printf("返回的图片数据: CoverImg=%s, Images=%+v\n", resp.CoverImg, resp.Images)

	return resp, nil
}

// UpdateGoods 更新商品
func (s *GoodsService) UpdateGoods(req *dto.GoodsUpdateReq) error {
	return db.DB.WnAdmin.Transaction(func(tx *gorm.DB) error {
		// 1. 查询商品
		goods := &shop.GoodsDO{Id: req.GoodsId}
		if err := tx.First(goods).Error; err != nil {
			return fmt.Errorf("商品不存在")
		}

		// 2. 更新商品基本信息
		goods.Name = req.GoodsName
		goods.CategoryId = req.CategoryId
		goods.BrandId = req.BrandId
		goods.SupplierId = req.SupplierId
		goods.Price = req.Price
		goods.OriginalPrice = req.OriginalPrice
		goods.Stock = req.Stock
		goods.Sort = req.Sort
		goods.Description = req.Description
		goods.SpecType = int8(req.SpecType)
		status, _ := strconv.ParseInt(req.Status, 10, 8)
		goods.Status = int8(status)

		if err := tx.Save(goods).Error; err != nil {
			return fmt.Errorf("更新商品失败: %v", err)
		}

		// 更新商品详情
		var detail shop.ShopGoodsDetail
		if err := tx.Where("goods_id = ?", goods.Id).First(&detail).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				detail = shop.ShopGoodsDetail{
					GoodsId: int(goods.Id),
				}
			} else {
				return fmt.Errorf("查询商品详情失败: %v", err)
			}
		}
		detail.DetailDesc = req.DetailDesc
		if err := tx.Save(&detail).Error; err != nil {
			return fmt.Errorf("保存商品详情失败: %v", err)
		}

		// 3. 处理规格信息
		if req.SpecType == 2 { // 多规格商品
			// 删除原有规格和规格值
			if err := tx.Where("goods_id = ?", goods.Id).Delete(&shop.ShopGoodsSpec{}).Error; err != nil {
				return fmt.Errorf("删除原有规格失败: %v", err)
			}
			if err := tx.Where("goods_id = ?", goods.Id).Delete(&shop.ShopGoodsSpecValue{}).Error; err != nil {
				return fmt.Errorf("删除原有规格值失败: %v", err)
			}
			// 删除原有SKU
			if err := tx.Where("goods_id = ?", goods.Id).Delete(&shop.ShopGoodsSku{}).Error; err != nil {
				return fmt.Errorf("删除原有SKU失败: %v", err)
			}

			// 保存新的规格和规格值
			for _, group := range req.SpecGroups {
				spec := &shop.ShopGoodsSpec{
					GoodsId:    goods.Id,
					Name:       group.Name,
					Sort:       group.Sort,
					IsDelete:   0,
					CreateTime: time.Now(),
					UpdateTime: time.Now(),
				}

				if err := tx.Create(spec).Error; err != nil {
					return fmt.Errorf("保存规格失败: %v", err)
				}

				// 保存规格值
				for _, value := range group.Values {
					specValue := &shop.ShopGoodsSpecValue{
						GoodsId:   goods.Id,
						SpecId:    spec.Id,
						Value:     value.Name,
						Sort:      value.Sort,
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					}

					if err := tx.Create(specValue).Error; err != nil {
						return fmt.Errorf("保存规格值失败: %v", err)
					}
				}
			}

			// 添加调试日志
			fmt.Printf("接收到的商品数据: %+v\n", req)
			fmt.Printf("SKU列表数据: %+v\n", req.SkuList)

			// 保存SKU列表
			for _, sku := range req.SkuList {
				specsJson, err := json.Marshal(sku.Specs)
				if err != nil {
					tx.Rollback()
					return fmt.Errorf("序列化规格数据失败: %v", err)
				}

				skuCode := fmt.Sprintf("SKU%d%d%d", goods.Id, time.Now().Unix(), rand.Int63n(1000000))
				skuData := &shop.ShopGoodsSku{
					GoodsId:       goods.Id,
					SkuCode:       skuCode,
					SpecsJson:     string(specsJson),
					Price:         sku.Price,
					OriginalPrice: goods.OriginalPrice,
					CostPrice:     goods.CostPrice,
					Stock:         sku.Stock,
					Status:        1,
					Sort:          0,
					Image:         sku.Image,
					CreateTime:    time.Now(),
					UpdateTime:    time.Now(),
				}

				if err := tx.Create(skuData).Error; err != nil {
					return fmt.Errorf("保存SKU失败: %v", err)
				}
			}
		} else if req.SpecType == 1 { // 单规格商品
			fmt.Printf("处理单规格商品\n")
			// 删除原有规格和规格值
			if err := tx.Where("goods_id = ?", goods.Id).Delete(&shop.ShopGoodsSpec{}).Error; err != nil {
				return fmt.Errorf("删除原有规格失败: %v", err)
			}
			if err := tx.Where("goods_id = ?", goods.Id).Delete(&shop.ShopGoodsSpecValue{}).Error; err != nil {
				return fmt.Errorf("删除原有规格值失败: %v", err)
			}
			// 删除原有SKU
			if err := tx.Where("goods_id = ?", goods.Id).Delete(&shop.ShopGoodsSku{}).Error; err != nil {
				return fmt.Errorf("删除原有SKU失败: %v", err)
			}

			// 创建默认SKU
			skuCode := fmt.Sprintf("SKU%d%d%d", goods.Id, time.Now().Unix(), rand.Int63n(1000000))
			sku := &shop.ShopGoodsSku{
				GoodsId:       int64(goods.Id),
				SkuCode:       skuCode,
				SpecsJson:     "{}", // 单规格商品，规格为空
				Price:         goods.Price,
				OriginalPrice: goods.OriginalPrice,
				CostPrice:     goods.CostPrice,
				Stock:         goods.Stock,
				Status:        1,
				Sort:          0,
				CreateTime:    time.Now(),
				UpdateTime:    time.Now(),
			}

			if err := tx.Create(sku).Error; err != nil {
				return fmt.Errorf("保存SKU失败: %v", err)
			}
		}

		// 4. 处理图片
		// 3.1 处理主图
		// 先删除所有主图
		if err := tx.Where("goods_id = ? AND position = ?", goods.Id, shop.ImagePositionMain).Delete(&shop.ShopGoodsImage{}).Error; err != nil {
			return fmt.Errorf("删除旧主图失败: %v", err)
		}

		// 只有当 CoverImg 不为空时才添加新的主图
		if req.CoverImg != "" {
			mainImage := &shop.ShopGoodsImage{
				GoodsId:    goods.Id,
				Position:   shop.ImagePositionMain,
				Url:        req.CoverImg,
				Sort:       0,
				IsDelete:   0,
				CreateTime: time.Now(),
				UpdateTime: time.Now(),
			}
			if err := tx.Create(mainImage).Error; err != nil {
				return fmt.Errorf("添加主图失败: %v", err)
			}
		}

		// 3.2 处理商品图片列表
		// 先删除所有非主图的图片
		if err := tx.Where("goods_id = ? AND position != ?", goods.Id, shop.ImagePositionMain).Delete(&shop.ShopGoodsImage{}).Error; err != nil {
			return fmt.Errorf("删除旧图片失败: %v", err)
		}

		// 添加新的图片列表，确保不添加主图
		if len(req.Images) > 0 {
			for _, img := range req.Images {
				// 跳过主图位置的图片
				if img.Position == shop.ImagePositionMain {
					continue
				}
				image := &shop.ShopGoodsImage{
					GoodsId:    goods.Id,
					Position:   img.Position,
					Url:        img.Url,
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

		// 更新标签关联
		if len(req.TagIds) > 0 {
			// 删除旧的标签关联
			if err := tx.Where("goods_id = ?", goods.Id).Delete(&shop.ShopGoodsTagRel{}).Error; err != nil {
				return fmt.Errorf("删除旧标签失败: %v", err)
			}

			// 添加新的标签关联
			for _, tagId := range req.TagIds {
				rel := &shop.ShopGoodsTagRel{
					GoodsId:    goods.Id,
					TagId:      tagId,
					CreateTime: time.Now(),
				}
				if err := tx.Create(rel).Error; err != nil {
					return fmt.Errorf("添加标签关联失败: %v", err)
				}
			}
		}

		return nil // 让事务自动提交
	})
}

// CreateGoods 创建商品
func (s *GoodsService) CreateGoods(req *dto.GoodsUpdateReq) error {
	// 开启事务
	tx := db.DB.WnAdmin.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 创建商品基本信息
	status, _ := strconv.ParseInt(req.Status, 10, 8)
	goods := &shop.GoodsDO{
		Name:          req.GoodsName,
		CategoryId:    req.CategoryId,
		BrandId:       req.BrandId,
		SupplierId:    req.SupplierId,
		Price:         req.Price,
		OriginalPrice: req.OriginalPrice,
		Stock:         req.Stock,
		Sort:          req.Sort,
		Description:   req.Description,
		SpecType:      int8(req.SpecType),
		Status:        int8(status),
		CreateTime:    time.Now(),
		UpdateTime:    time.Now(),
	}

	if err := tx.Create(goods).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("创建商品失败: %v", err)
	}

	// 创建商品详情
	detail := &shop.ShopGoodsDetail{
		GoodsId:    int(goods.Id),
		DetailDesc: req.DetailDesc,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}
	if err := tx.Create(detail).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("创建商品详情失败: %v", err)
	}

	// 添加商品主图
	if req.CoverImg != "" {
		mainImage := &shop.ShopGoodsImage{
			GoodsId:    goods.Id,
			Position:   shop.ImagePositionMain,
			Url:        req.CoverImg,
			Sort:       0,
			CreateTime: time.Now(),
			UpdateTime: time.Now(),
		}
		if err := tx.Create(mainImage).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("添加主图失败: %v", err)
		}
	}

	// 添加标签关联
	if len(req.TagIds) > 0 {
		for _, tagId := range req.TagIds {
			rel := &shop.ShopGoodsTagRel{
				GoodsId:    goods.Id,
				TagId:      tagId,
				CreateTime: time.Now(),
			}
			if err := tx.Create(rel).Error; err != nil {
				tx.Rollback()
				return fmt.Errorf("添加标签关联失败: %v", err)
			}
		}
	}

	// 如果是多规格商品，保存规格信息
	if req.SpecType == 2 && len(req.SpecGroups) > 0 {
		// 保存规格组
		for _, group := range req.SpecGroups {
			specGroup := &shop.ShopGoodsSpec{
				GoodsId:    goods.Id,
				Name:       group.Name,
				Sort:       group.Sort,
				CreateTime: time.Now(),
				UpdateTime: time.Now(),
			}
			if err := tx.Create(specGroup).Error; err != nil {
				tx.Rollback()
				return fmt.Errorf("创建规格组失败: %v", err)
			}

			// 保存规格值
			for _, value := range group.Values {
				specValue := &shop.ShopGoodsSpecValue{
					GoodsId:   goods.Id,
					SpecId:    specGroup.Id,
					Value:     value.Value,
					Sort:      value.Sort,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}
				if err := tx.Create(specValue).Error; err != nil {
					tx.Rollback()
					return fmt.Errorf("创建规格值失败: %v", err)
				}
			}
		}

		// 保存 SKU 列表
		if len(req.SkuList) > 0 {
			for _, sku := range req.SkuList {
				specsJson, err := json.Marshal(sku.Specs)
				if err != nil {
					tx.Rollback()
					return fmt.Errorf("序列化规格数据失败: %v", err)
				}

				skuCode := fmt.Sprintf("SKU%d%d%d", goods.Id, time.Now().Unix(), rand.Int63n(1000000))
				skuData := &shop.ShopGoodsSku{
					GoodsId:       goods.Id,
					SkuCode:       skuCode,
					SpecsJson:     string(specsJson),
					Price:         sku.Price,
					OriginalPrice: goods.OriginalPrice,
					CostPrice:     goods.CostPrice,
					Stock:         sku.Stock,
					Status:        1,
					Sort:          0,
					Image:         sku.Image,
					CreateTime:    time.Now(),
					UpdateTime:    time.Now(),
				}

				if err := tx.Create(skuData).Error; err != nil {
					return fmt.Errorf("保存SKU失败: %v", err)
				}
			}
		}
	}

	return tx.Commit().Error
}

func (s *GoodsService) GetByID(id int64) (*shop.GoodsDO, error) {
	goods := &shop.GoodsDO{}

	// 查询商品基本信息
	if err := db.DB.WnAdmin.First(goods, id).Error; err != nil {
		return nil, err
	}

	// 如果是多规格商品，查询规格信息
	if goods.SpecType == 2 {
		// 查询规格组
		var specGroups []*shop.ShopGoodsSpec
		if err := db.DB.WnAdmin.Where("goods_id = ? AND is_delete = 0", id).
			Order("sort ASC").Find(&specGroups).Error; err != nil {
			return nil, err
		}

		// 查询规格值
		for _, group := range specGroups {
			var values []*shop.ShopGoodsSpecValue
			if err := db.DB.WnAdmin.Where("goods_id = ? AND spec_id = ?", id, group.Id).
				Order("sort ASC").Find(&values).Error; err != nil {
				return nil, err
			}
			// 转换为前端需要的格式
			for _, value := range values {
				group.Values = append(group.Values, &shop.ShopGoodsSpecValue{
					Id:        value.Id,
					GoodsId:   value.GoodsId,
					SpecId:    value.SpecId,
					Value:     value.Value,
					Sort:      value.Sort,
					CreatedAt: value.CreatedAt,
					UpdatedAt: value.UpdatedAt,
				})
			}
		}
		goods.SpecGroups = specGroups

		// 查询SKU信息
		var skuList []*shop.ShopGoodsSku
		if err := db.DB.WnAdmin.Where("goods_id = ? AND status = 1", id).
			Order("sort ASC").Find(&skuList).Error; err != nil {
			return nil, err
		}
		goods.SkuList = skuList
	}

	// 查询商品图片
	var images []*shop.ShopGoodsImage
	if err := db.DB.WnAdmin.Where("goods_id = ?", id).Find(&images).Error; err != nil {
		return nil, err
	}
	goods.Images = images

	// 查询商品标签
	var tagRels []*shop.ShopGoodsTagRel
	if err := db.DB.WnAdmin.Where("goods_id = ?", id).Find(&tagRels).Error; err != nil {
		return nil, err
	}
	var tagIds []int64
	for _, rel := range tagRels {
		tagIds = append(tagIds, rel.TagId)
	}
	goods.TagIds = tagIds

	return goods, nil
}
