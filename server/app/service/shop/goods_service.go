package shop

import (
	"encoding/json"
	"fmt"
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
		Detail:        goods.Detail,
		SpecType:      goods.SpecType,
		TagIds:        goods.TagIds,
	}

	// 设置封面图
	if len(goods.Images) > 0 {
		for _, img := range goods.Images {
			if img.IsMain == 1 {
				resp.CoverImg = img.Url
				break
			}
		}
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
			var specs map[string]dto.SpecValue
			if err := json.Unmarshal([]byte(sku.SpecsJson), &specs); err == nil {
				resp.SkuList[i] = dto.Sku{
					ID:    sku.Id,
					Specs: specs,
					Price: sku.Price,
					Stock: sku.Stock,
				}
			}
		}
	}

	return resp, nil
}

// UpdateGoods 更新商品
func (s *GoodsService) UpdateGoods(req *dto.GoodsUpdateReq) error {
	goods := &shop.GoodsDO{Id: req.GoodsId}
	if err := goods.First(); err != nil {
		return fmt.Errorf("商品不存在")
	}

	// 开启事务
	tx := db.DB.WnAdmin.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 更新商品基本信息
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
		tx.Rollback()
		return fmt.Errorf("更新商品失败: %v", err)
	}

	// 处理规格信息
	if req.SpecType == 1 && len(req.SpecGroups) > 0 {
		// 删除原有规格和规格值
		if err := tx.Where("goods_id = ?", goods.Id).Delete(&shop.ShopGoodsSpec{}).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("删除原有规格失败: %v", err)
		}
		if err := tx.Where("goods_id = ?", goods.Id).Delete(&shop.ShopGoodsSpecValue{}).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("删除原有规格值失败: %v", err)
		}
		// 删除原有SKU
		if err := tx.Where("goods_id = ?", goods.Id).Delete(&shop.ShopGoodsSku{}).Error; err != nil {
			tx.Rollback()
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
				tx.Rollback()
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
					tx.Rollback()
					return fmt.Errorf("保存规格值失败: %v", err)
				}
			}
		}

		// 保存SKU列表
		if len(req.SkuList) > 0 {
			for _, skuData := range req.SkuList {
				// 将规格信息转换为JSON字符串
				specsJSON, err := json.Marshal(skuData.Specs)
				if err != nil {
					return fmt.Errorf("序列化SKU规格信息失败: %v", err)
				}

				sku := &shop.ShopGoodsSku{
					GoodsId:       int64(goods.Id),
					SkuCode:       fmt.Sprintf("SKU%d%d", goods.Id, time.Now().UnixNano()),
					SpecsJson:     string(specsJSON),
					Price:         skuData.Price,
					OriginalPrice: goods.OriginalPrice,
					CostPrice:     goods.CostPrice,
					Stock:         skuData.Stock,
					Status:        1,
					Sort:          0,
					CreateTime:    time.Now(),
					UpdateTime:    time.Now(),
				}

				if err := tx.Create(sku).Error; err != nil {
					return fmt.Errorf("保存SKU失败: %v", err)
				}
			}
		}
	}

	// 更新商品主图
	if req.CoverImg != "" {
		// 删除旧的主图
		if err := tx.Where("goods_id = ? AND is_main = ?", goods.Id, 1).Delete(&shop.GoodsImageDO{}).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("删除旧主图失败: %v", err)
		}

		// 添加新的主图
		mainImage := &shop.GoodsImageDO{
			GoodsId:    goods.Id,
			IsMain:     1,
			Url:        req.CoverImg,
			Type:       1,
			Sort:       0,
			CreateTime: time.Now(),
			UpdateTime: time.Now(),
		}
		if err := tx.Create(mainImage).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("添加主图失败: %v", err)
		}
	}

	// 更新标签关联
	if len(req.TagIds) > 0 {
		// 删除旧的标签关联
		if err := tx.Where("goods_id = ?", goods.Id).Delete(&shop.ShopGoodsTagRel{}).Error; err != nil {
			tx.Rollback()
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
				tx.Rollback()
				return fmt.Errorf("添加标签关联失败: %v", err)
			}
		}
	}

	return tx.Commit().Error
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

	// 在事务开始时启用 SQL 日志
	tx = tx.Debug()

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
	}

	if err := tx.Create(goods).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("创建商品失败: %v", err)
	}

	// 添加商品主图
	if req.CoverImg != "" {
		mainImage := &shop.GoodsImageDO{
			GoodsId:    goods.Id,
			IsMain:     1,
			Url:        req.CoverImg,
			Type:       1,
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

	return tx.Commit().Error
}

// SaveGoods 保存商品
func (s *GoodsService) SaveGoods(req *dto.GoodsUpdateReq) error {
	// 打印接收到的参数
	reqJSON, _ := json.MarshalIndent(req, "", "  ")
	fmt.Printf("接收到的商品保存请求参数:\n%s\n", string(reqJSON))

	return db.DB.WnAdmin.Transaction(func(tx *gorm.DB) error {
		// 1. 保存商品基本信息
		status, _ := strconv.ParseInt(req.Status, 10, 8)
		goods := &shop.GoodsDO{
			Id:            req.GoodsId,
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
		}

		if err := tx.Save(goods).Error; err != nil {
			return fmt.Errorf("保存商品基本信息失败: %v", err)
		}

		// 2. 处理规格信息
		if req.SpecType == 1 && len(req.SpecGroups) > 0 {
			fmt.Printf("规格组数量: %d\n", len(req.SpecGroups))
			for i, group := range req.SpecGroups {
				fmt.Printf("规格组[%d]: ID=%d, Name=%s, 规格值数量=%d\n",
					i, group.ID, group.Name, len(group.Values))

				// 打印规格值
				for j, value := range group.Values {
					fmt.Printf("  规格值[%d]: ID=%d, Name=%s, Value=%s, Sort=%d\n",
						j, value.ID, value.Name, value.Value, value.Sort)
				}
			}

			// 打印SKU信息
			fmt.Printf("SKU列表数量: %d\n", len(req.SkuList))
			for i, sku := range req.SkuList {
				specsJSON, _ := json.MarshalIndent(sku.Specs, "", "  ")
				fmt.Printf("SKU[%d]: ID=%d, Price=%.2f, Stock=%d\n规格信息:\n%s\n",
					i, sku.ID, sku.Price, sku.Stock, string(specsJSON))
			}

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
				// 打印规格组信息
				fmt.Printf("正在保存规格组: Name=%s\n", group.Name)

				spec := &shop.ShopGoodsSpec{
					GoodsId:    goods.Id,
					Name:       group.Name,
					Sort:       group.Sort,
					IsDelete:   0,
					CreateTime: time.Now(),
					UpdateTime: time.Now(),
				}

				if err := tx.Create(spec).Error; err != nil {
					fmt.Printf("保存规格组失败: %v\n", err)
					return fmt.Errorf("保存规格失败: %v", err)
				}
				fmt.Printf("规格组保存成功，获得ID: %d\n", spec.Id)

				// 保存规格值
				for _, value := range group.Values {
					fmt.Printf("正在保存规格值: Name=%s, SpecId=%d\n", value.Name, spec.Id)

					specValue := &shop.ShopGoodsSpecValue{
						GoodsId:   goods.Id,
						SpecId:    spec.Id,
						Value:     value.Name,
						Sort:      value.Sort,
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					}

					if err := tx.Create(specValue).Error; err != nil {
						fmt.Printf("SQL: %v\n", tx.ToSQL(func(tx *gorm.DB) *gorm.DB {
							return tx.Create(specValue)
						}))
						fmt.Printf("Error: %v\n", err)
						return fmt.Errorf("保存规格值失败: %v", err)
					}
					fmt.Printf("规格值保存成功\n")
				}
			}

			// 保存SKU列表
			if len(req.SkuList) > 0 {
				for _, skuData := range req.SkuList {
					// 将规格信息转换为JSON字符串
					specsJSON, err := json.Marshal(skuData.Specs)
					if err != nil {
						return fmt.Errorf("序列化SKU规格信息失败: %v", err)
					}

					sku := &shop.ShopGoodsSku{
						GoodsId:       int64(goods.Id),
						SkuCode:       fmt.Sprintf("SKU%d%d", goods.Id, time.Now().UnixNano()),
						SpecsJson:     string(specsJSON),
						Price:         skuData.Price,
						OriginalPrice: goods.OriginalPrice,
						CostPrice:     goods.CostPrice,
						Stock:         skuData.Stock,
						Status:        1,
						Sort:          0,
						CreateTime:    time.Now(),
						UpdateTime:    time.Now(),
					}

					if err := tx.Create(sku).Error; err != nil {
						return fmt.Errorf("保存SKU失败: %v", err)
					}
				}
			}
		}

		// 3. 处理商品图片
		if req.CoverImg != "" {
			// 删除旧的主图
			if err := tx.Where("goods_id = ? AND is_main = ?", goods.Id, 1).Delete(&shop.GoodsImageDO{}).Error; err != nil {
				return fmt.Errorf("删除旧主图失败: %v", err)
			}

			// 添加新的主图
			mainImage := &shop.GoodsImageDO{
				GoodsId:    goods.Id,
				IsMain:     1,
				Url:        req.CoverImg,
				Type:       1,
				Sort:       0,
				CreateTime: time.Now(),
				UpdateTime: time.Now(),
			}
			if err := tx.Create(mainImage).Error; err != nil {
				return fmt.Errorf("添加主图失败: %v", err)
			}
		}

		// 4. 处理商品标签
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

		return nil
	})
}

func (s *GoodsService) GetByID(id int64) (*shop.GoodsDO, error) {
	goods := &shop.GoodsDO{}

	// 查询商品基本信息
	if err := db.DB.WnAdmin.First(goods, id).Error; err != nil {
		return nil, err
	}

	// 如果是多规格商品，查询规格信息
	if goods.SpecType == 1 {
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
	var images []*shop.GoodsImageDO
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
