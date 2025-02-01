package shop

import (
	"wnadmin/app/model/do/shop"
	"wnadmin/app/model/dto"
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
			Where("goods_id = ?", g.GoodsId).
			Find(&tagRels).Error; err == nil {
			for _, rel := range tagRels {
				tagIds = append(tagIds, rel.TagId)
				if name, ok := tagMap[rel.TagId]; ok {
					tagNames = append(tagNames, name)
				}
			}
		}

		item := dto.GoodsDTO{
			GoodsId:       g.GoodsId,
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
		goodsIds = append(goodsIds, g.GoodsId)
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
	goodsDO := &shop.GoodsDO{GoodsId: id}
	return goodsDO.Delete()
}

// UpdateGoodsStatus 更新商品状态
func (s *GoodsService) UpdateGoodsStatus(id int64, status int8) error {
	goodsDO := &shop.GoodsDO{
		GoodsId: id,
		Status:  status,
	}
	return goodsDO.UpdateStatus()
}
