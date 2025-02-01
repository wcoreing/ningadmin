package shop

import (
	"wnadmin/app/model/do/shop"
	"wnadmin/app/model/dto"
)

// BrandService 品牌服务
type BrandService struct{}

// GetBrandList 获取品牌列表
func (s *BrandService) GetBrandList() ([]dto.BrandDTO, error) {
	var brands []shop.ShopBrand

	// 创建查询对象并执行查询
	if err := (&shop.ShopBrand{
		Status: 1, // 只查询启用状态的品牌
	}).Find(&brands); err != nil {
		return nil, err
	}

	var result []dto.BrandDTO
	for _, brand := range brands {
		result = append(result, dto.BrandDTO{
			ID:          brand.ID,
			Name:        brand.Name,
			Logo:        brand.Logo,
			Description: brand.Description,
			Website:     brand.Website,
			Sort:        brand.Sort,
		})
	}
	return result, nil
}
