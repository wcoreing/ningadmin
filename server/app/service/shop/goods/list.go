package goods

import (
	"strconv"
	"wnadmin/app/common/utils"
	"wnadmin/app/model/do/shop"
	"wnadmin/app/model/dto"
	"wnadmin/app/provider/db"
)

// ListService 商品列表查询服务
type ListService struct{}

// Execute 执行商品列表查询
func (s *ListService) Execute(req *dto.GoodsListReq) (*dto.PageResult, error) {
	var (
		goods  []shop.GoodsDO
		total  int64
		result = make([]dto.GoodsDTO, 0)
	)

	// 处理分类查询
	if err := s.processCategoryFilter(req); err != nil {
		return nil, err
	}

	// 处理名称查询
	s.processNameFilter(req)

	// 查询商品列表
	goods, total, err := (&shop.GoodsDO{}).GetGoodsList(req)
	if err != nil {
		return nil, err
	}

	// 获取关联数据
	categoryMap := s.getCategoryMap(goods)
	tagMap := s.getTagMap(goods)

	// 组装数据
	for _, g := range goods {
		item := s.buildGoodsDTO(g, categoryMap, tagMap)
		result = append(result, item)
	}

	// 查询商品图片
	s.loadGoodsImages(result)

	return &dto.PageResult{
		List:     result,
		Total:    total,
		PageNum:  req.PageNum,
		PageSize: req.PageSize,
	}, nil
}

// processCategoryFilter 处理分类过滤
func (s *ListService) processCategoryFilter(req *dto.GoodsListReq) error {
	if req.CategoryId <= 0 {
		return nil
	}

	var categoryIds []int64
	// 先添加当前分类ID
	categoryIds = append(categoryIds, req.CategoryId)

	// 查询所有子分类
	var categories []shop.Category
	if err := (&shop.Category{}).Query().Where("parent_id = ? AND module = ?", req.CategoryId, "goods").Find(&categories).Error; err == nil {
		for _, category := range categories {
			categoryIds = append(categoryIds, category.ID)
			// 如果还有下一级，继续查询
			var subCategories []shop.Category
			if err := (&shop.Category{}).Query().Where("parent_id = ? AND module = ?", category.ID, "goods").Find(&subCategories).Error; err == nil {
				for _, subCategory := range subCategories {
					categoryIds = append(categoryIds, subCategory.ID)
				}
			}
		}
	}

	// 修改查询条件，使用 IN 查询
	req.CategoryIds = categoryIds
	return nil
}

// processNameFilter 处理名称过滤
func (s *ListService) processNameFilter(req *dto.GoodsListReq) {
	if req.Name == "" {
		return
	}

	// 尝试将name转换为ID
	if id, err := strconv.ParseInt(req.Name, 10, 64); err == nil {
		// 如果可以转换为数字，添加ID查询条件
		req.GoodsId = id
	}
}

// getCategoryMap 获取分类Map
func (s *ListService) getCategoryMap(goods []shop.GoodsDO) map[int64]string {
	result := make(map[int64]string)
	var categoryIds []int64
	for _, g := range goods {
		if g.CategoryId > 0 {
			categoryIds = append(categoryIds, g.CategoryId)
		}
	}
	if len(categoryIds) > 0 {
		var categories []shop.Category
		if err := (&shop.Category{}).Query().Where("id IN ? AND module = ?", categoryIds, "goods").Find(&categories).Error; err == nil {
			for _, c := range categories {
				result[c.ID] = c.Name
			}
		}
	}
	return result
}

// getTagMap 获取标签Map
func (s *ListService) getTagMap(goods []shop.GoodsDO) map[int64][]string {
	result := make(map[int64][]string)
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

			// 按商品ID分组标签名称
			for _, rel := range tagRels {
				if rel.TagName != "" {
					result[rel.GoodsId] = append(result[rel.GoodsId], rel.TagName)
				}
			}
		}
	}
	return result
}

// buildGoodsDTO 构建商品DTO
func (s *ListService) buildGoodsDTO(g shop.GoodsDO, categoryMap map[int64]string, tagMap map[int64][]string) dto.GoodsDTO {
	// 获取商品的标签名称
	var tagNames []string
	if tags, ok := tagMap[g.Id]; ok {
		tagNames = tags
	}

	// 获取商品的SKU信息来获取价格和库存
	var skus []shop.ShopGoodsSku
	var minPrice float64 = 0
	var maxPrice float64 = 0
	var totalStock int = 0

	if err := db.WnAdmin().Where("goods_id = ? AND status = 1", g.Id).Find(&skus).Error; err == nil && len(skus) > 0 {
		// 初始化价格范围
		minPrice = skus[0].Price
		maxPrice = skus[0].Price

		// 计算价格范围和总库存
		for _, sku := range skus {
			if sku.Price < minPrice {
				minPrice = sku.Price
			}
			if sku.Price > maxPrice {
				maxPrice = sku.Price
			}
			totalStock += sku.Stock
		}
	}

	return dto.GoodsDTO{
		GoodsId:      g.Id,
		CategoryId:   g.CategoryId,
		CategoryName: categoryMap[g.CategoryId],
		Brand:        g.Brand,
		Name:         g.Name,
		Status:       g.Status,
		Sort:         g.Sort,
		TagNames:     tagNames,   // 标签名称列表
		MinPrice:     minPrice,   // 最低价格
		MaxPrice:     maxPrice,   // 最高价格
		TotalStock:   totalStock, // 总库存
	}
}

// loadGoodsImages 加载商品图片
func (s *ListService) loadGoodsImages(result []dto.GoodsDTO) {
	for i := range result {
		var images []*shop.ShopGoodsImage
		if err := db.WnAdmin().Where("goods_id = ? AND position = ?", result[i].GoodsId, shop.ImagePositionMain).
			First(&images).Error; err == nil && len(images) > 0 {
			result[i].CoverImg = utils.AddDomainToURL(images[0].Url)
		}
	}
}
