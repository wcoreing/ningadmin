package shop

import (
	"wnadmin/app/model/do/shop"
	"wnadmin/app/model/dto"
)

// TagService 商品标签服务
type TagService struct{}

// GetTagList 获取商品标签列表
func (s *TagService) GetTagList() ([]dto.TagDTO, error) {
	var tags []shop.ShopGoodsTag

	// 创建查询对象并执行查询
	if err := (&shop.ShopGoodsTag{}).Find(&tags); err != nil {
		return nil, err
	}

	var result []dto.TagDTO
	for _, tag := range tags {
		result = append(result, dto.TagDTO{
			ID:   tag.ID,
			Name: tag.Name,
			Sort: tag.Sort,
		})
	}
	return result, nil
}
