package shop

import (
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
	goods := &shop.GoodsDO{Id: id}
	if err := goods.First(); err != nil {
		return nil, err
	}

	// 获取商品主图
	var mainImage shop.GoodsImageDO
	if err := (&shop.GoodsImageDO{}).Query().
		Where("goods_id = ? AND is_main = 1", id).
		First(&mainImage).Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	// 获取商品标签关联
	var tagRels []shop.ShopGoodsTagRel
	var tagIds []int64
	if err := (&shop.ShopGoodsTagRel{}).Query().
		Where("goods_id = ?", id).
		Find(&tagRels).Error; err == nil {
		for _, rel := range tagRels {
			tagIds = append(tagIds, rel.TagId)
		}
	}

	return &dto.GoodsInfoResp{
		GoodsId:       goods.Id,
		GoodsName:     goods.Name,
		CategoryId:    goods.CategoryId,
		BrandId:       goods.BrandId,
		SupplierId:    goods.SupplierId,
		Price:         goods.Price,
		OriginalPrice: goods.OriginalPrice,
		Stock:         goods.Stock,
		Sort:          goods.Sort,
		Status:        fmt.Sprintf("%d", goods.Status),
		Description:   goods.Description,
		CoverImg:      mainImage.Url,
		TagIds:        tagIds,
	}, nil
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
	status, _ := strconv.ParseInt(req.Status, 10, 8)
	goods.Status = int8(status)

	if err := tx.Save(goods).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("更新商品失败: %v", err)
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
