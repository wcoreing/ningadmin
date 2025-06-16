package common

import (
	"errors"
	"wnadmin/app/model/do/shop"
	"wnadmin/app/model/dto/common"
	"wnadmin/app/provider/db"
)

// CategoryService 通用分类服务
type CategoryService struct{}

// GetCategoryTree 获取分类树
func (s *CategoryService) GetCategoryTree(module string) ([]shop.Category, error) {
	category := &shop.Category{}
	return category.BuildTree(module)
}

// GetCategoryList 获取分类列表
func (s *CategoryService) GetCategoryList(module string) ([]shop.Category, error) {
	category := &shop.Category{}
	return category.GetByModule(module)
}

// GetCategoryOptions 获取分类选项（用于前端下拉框）
func (s *CategoryService) GetCategoryOptions(module string) ([]common.CategoryOption, error) {
	categories, err := s.GetCategoryTree(module)
	if err != nil {
		return nil, err
	}

	var options []common.CategoryOption
	for _, category := range categories {
		option := s.buildCategoryOption(category)
		options = append(options, option)
	}

	return options, nil
}

// CreateCategory 创建分类
func (s *CategoryService) CreateCategory(req *common.CategoryCreateReq) error {
	// 检查父级分类
	if req.ParentId > 0 {
		var parent shop.Category
		if err := db.WnAdmin().Where("module = ?", req.Module).First(&parent, req.ParentId).Error; err != nil {
			return errors.New("父级分类不存在")
		}
		if parent.Level >= 3 {
			return errors.New("最多支持三级分类")
		}
	}

	// 创建分类
	category := &shop.Category{
		ParentId: req.ParentId,
		Name:     req.Name,
		Icon:     req.Icon,
		Image:    req.Image,
		Level:    s.getLevel(req.ParentId, req.Module),
		Sort:     req.Sort,
		Status:   req.Status,
		Module:   req.Module,
	}

	return db.WnAdmin().Create(category).Error
}

// UpdateCategory 更新分类
func (s *CategoryService) UpdateCategory(req *common.CategoryUpdateReq) error {
	var category shop.Category
	if err := db.WnAdmin().Where("module = ?", req.Module).First(&category, req.Id).Error; err != nil {
		return errors.New("分类不存在")
	}

	// 检查父级分类
	if req.ParentId > 0 && req.ParentId != category.ParentId {
		var parent shop.Category
		if err := db.WnAdmin().Where("module = ?", req.Module).First(&parent, req.ParentId).Error; err != nil {
			return errors.New("父级分类不存在")
		}
		if parent.Level >= 3 {
			return errors.New("最多支持三级分类")
		}
		category.Level = s.getLevel(req.ParentId, req.Module)
	}

	category.ParentId = req.ParentId
	category.Name = req.Name
	category.Icon = req.Icon
	category.Image = req.Image
	category.Sort = req.Sort
	category.Status = req.Status

	return db.WnAdmin().Save(&category).Error
}

// DeleteCategory 删除分类
func (s *CategoryService) DeleteCategory(id int64) error {
	// 先查询分类信息获取模块
	var category shop.Category
	if err := db.WnAdmin().First(&category, id).Error; err != nil {
		return errors.New("分类不存在")
	}

	// 检查是否有子分类
	var count int64
	if err := db.WnAdmin().Model(&shop.Category{}).Where("parent_id = ? AND is_delete = ? AND module = ?", id, 0, category.Module).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errors.New("请先删除子分类")
	}

	// 根据模块检查是否有关联数据
	switch category.Module {
	case "goods":
		// 检查是否有关联商品
		if err := db.WnAdmin().Model(&shop.GoodsDO{}).Where("category_id = ? AND is_delete = ?", id, 0).Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			return errors.New("该分类下有商品，无法删除")
		}
	case "combo":
		// 检查是否有关联套餐
		if err := db.WnAdmin().Model(&shop.ShopGoodsComboDO{}).Where("category_id = ? AND is_delete = ?", id, 0).Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			return errors.New("该分类下有套餐，无法删除")
		}
	}

	return db.WnAdmin().Model(&shop.Category{}).Where("id = ?", id).Update("is_delete", 1).Error
}

// buildCategoryOption 构建分类选项
func (s *CategoryService) buildCategoryOption(category shop.Category) common.CategoryOption {
	option := common.CategoryOption{
		ID:       category.ID,
		Name:     category.Name,
		ParentId: category.ParentId,
		Level:    category.Level,
		Sort:     category.Sort,
		Status:   category.Status,
		Icon:     category.Icon,
		Image:    category.Image,
	}

	// 递归构建子分类
	if len(category.Children) > 0 {
		for _, child := range category.Children {
			childOption := s.buildCategoryOption(child)
			option.Children = append(option.Children, childOption)
		}
	}

	return option
}

// getLevel 获取分类层级
func (s *CategoryService) getLevel(parentId int64, module string) int {
	if parentId == 0 {
		return 1
	}
	var parent shop.Category
	if err := db.WnAdmin().Where("module = ?", module).First(&parent, parentId).Error; err != nil {
		return 1
	}
	return parent.Level + 1
}
