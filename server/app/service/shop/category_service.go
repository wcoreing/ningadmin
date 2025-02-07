package shop

import (
	"errors"
	shop_do "wnadmin/app/model/do/shop"
	shop_dto "wnadmin/app/model/dto/shop"
	"wnadmin/app/provider/db"
)

// CategoryService 商品分类服务
type CategoryService struct {
	shopCategory *shop_do.ShopCategory
}

func NewCategoryService() *CategoryService {
	return &CategoryService{
		shopCategory: &shop_do.ShopCategory{},
	}
}

// GetList 获取分类列表(树形结构)
func (s *CategoryService) GetList(req *shop_dto.CategoryListReq) ([]shop_dto.CategoryVO, error) {
	var categories []shop_do.ShopCategory

	// 构建查询条件
	query := db.DB.WnAdmin.Model(&shop_do.ShopCategory{}).Where("is_delete = ?", 0)
	if req.Name != "" {
		query = query.Where("name LIKE ?", "%"+req.Name+"%")
	}
	if req.Status != nil {
		query = query.Where("status = ?", *req.Status)
	}
	if req.ParentId != nil {
		query = query.Where("parent_id = ?", *req.ParentId)
	}

	// 获取所有分类
	if err := query.Order("sort DESC, id ASC").Find(&categories).Error; err != nil {
		return nil, err
	}

	// 转换为树形结构
	return s.buildTree(categories, 0), nil
}

// Create 创建分类
func (s *CategoryService) Create(req *shop_dto.CategoryCreateReq) error {
	// 检查父级分类
	if req.ParentId > 0 {
		var parent shop_do.ShopCategory
		if err := db.DB.WnAdmin.First(&parent, req.ParentId).Error; err != nil {
			return errors.New("父级分类不存在")
		}
		if parent.Level >= 3 {
			return errors.New("最多支持三级分类")
		}
	}

	// 创建分类
	category := &shop_do.ShopCategory{
		ParentId: req.ParentId,
		Name:     req.Name,
		Icon:     req.Icon,
		Image:    req.Image,
		Level:    s.getLevel(req.ParentId),
		Sort:     req.Sort,
		Status:   req.Status,
	}

	return db.DB.WnAdmin.Create(category).Error
}

// Update 更新分类
func (s *CategoryService) Update(req *shop_dto.CategoryUpdateReq) error {
	var category shop_do.ShopCategory
	if err := db.DB.WnAdmin.First(&category, req.Id).Error; err != nil {
		return errors.New("分类不存在")
	}

	// 检查父级分类
	if req.ParentId > 0 && req.ParentId != category.ParentId {
		var parent shop_do.ShopCategory
		if err := db.DB.WnAdmin.First(&parent, req.ParentId).Error; err != nil {
			return errors.New("父级分类不存在")
		}
		if parent.Level >= 3 {
			return errors.New("最多支持三级分类")
		}
		category.Level = s.getLevel(req.ParentId)
	}

	category.ParentId = req.ParentId
	category.Name = req.Name
	category.Icon = req.Icon
	category.Image = req.Image
	category.Sort = req.Sort
	category.Status = req.Status

	return db.DB.WnAdmin.Save(&category).Error
}

// Delete 删除分类
func (s *CategoryService) Delete(id int64) error {
	// 检查是否有子分类
	var count int64
	if err := db.DB.WnAdmin.Model(&shop_do.ShopCategory{}).Where("parent_id = ? AND is_delete = ?", id, 0).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errors.New("请先删除子分类")
	}

	// 检查是否有关联商品
	if err := db.DB.WnAdmin.Model(&shop_do.GoodsDO{}).Where("category_id = ? AND is_delete = ?", id, 0).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errors.New("该分类下有商品，无法删除")
	}

	return db.DB.WnAdmin.Model(&shop_do.ShopCategory{}).Where("id = ?", id).Update("is_delete", 1).Error
}

// 构建树形结构
func (s *CategoryService) buildTree(categories []shop_do.ShopCategory, parentId int64) []shop_dto.CategoryVO {
	var result []shop_dto.CategoryVO
	for _, category := range categories {
		if category.ParentId == parentId {
			vo := shop_dto.CategoryVO{
				Id:         category.ID,
				ParentId:   category.ParentId,
				Name:       category.Name,
				Icon:       category.Icon,
				Image:      category.Image,
				Level:      category.Level,
				Sort:       category.Sort,
				Status:     category.Status,
				CreateTime: category.CreateTime.Format("2006-01-02 15:04:05"),
				UpdateTime: category.UpdateTime.Format("2006-01-02 15:04:05"),
			}
			vo.Children = s.buildTree(categories, category.ID)
			result = append(result, vo)
		}
	}
	return result
}

// 获取分类层级
func (s *CategoryService) getLevel(parentId int64) int {
	if parentId == 0 {
		return 1
	}
	var parent shop_do.ShopCategory
	if err := db.DB.WnAdmin.First(&parent, parentId).Error; err != nil {
		return 1
	}
	return parent.Level + 1
}
