package shop

import (
	"fmt"
	"wnadmin/app/model/do/shop"
	"wnadmin/app/model/dto"
)

// CategoryService 商品分类服务
type CategoryService struct{}

// GetCategoryList 获取商品分类列表
func (s *CategoryService) GetCategoryList() ([]dto.CategoryDTO, error) {
	var categories []shop.ShopCategory

	// 创建查询对象并执行查询
	if err := (&shop.ShopCategory{
		Status: 1, // 只查询启用状态的分类
	}).Find(&categories); err != nil {
		return nil, err
	}

	// 添加调试日志
	fmt.Printf("查询到的分类数据: %+v\n", categories)

	var result []dto.CategoryDTO
	for _, category := range categories {
		result = append(result, dto.CategoryDTO{
			ID:       category.ID,
			Name:     category.Name,
			ParentID: category.ParentID,
			Level:    category.Level,
			Sort:     category.Sort,
		})
	}

	// 添加调试日志
	fmt.Printf("返回的DTO数据: %+v\n", result)
	return result, nil
}
