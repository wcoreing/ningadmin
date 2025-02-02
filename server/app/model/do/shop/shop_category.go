package shop

import (
	"time"
	"wnadmin/app/provider/db"

	"gorm.io/gorm"
)

// ShopCategory 商品分类
type ShopCategory struct {
	ID         int64          `json:"id" gorm:"column:id;primaryKey;comment:分类ID"`
	ParentId   int64          `json:"parentId" gorm:"column:parent_id;default:0;comment:父级ID"`
	Name       string         `json:"name" gorm:"column:name;comment:分类名称"`
	Icon       string         `gorm:"column:icon" json:"icon"`   // 分类图标
	Image      string         `gorm:"column:image" json:"image"` // 分类图片
	Level      int8           `gorm:"column:level" json:"level"` // 层级
	Sort       int            `json:"sort" gorm:"column:sort;default:0;comment:排序"`
	Status     int8           `json:"status" gorm:"column:status;default:1;comment:状态(0:禁用 1:启用)"`
	IsDelete   int8           `json:"-" gorm:"column:is_delete;default:0;comment:是否删除"`
	CreateTime time.Time      `json:"createTime" gorm:"column:create_time;autoCreateTime;comment:创建时间"`
	UpdateTime time.Time      `json:"updateTime" gorm:"column:update_time;autoUpdateTime;comment:更新时间"`
	Children   []ShopCategory `json:"children" gorm:"-"` // 子分类列表
}

// TableName 表名
func (ShopCategory) TableName() string {
	return "shop_category"
}

// Query 获取查询构建器
func (c *ShopCategory) Query() *gorm.DB {
	return db.DB.WnAdmin.Model(c).Where("is_delete = ?", 0)
}

// BuildTree 构建分类树
func BuildCategoryTree(categories []ShopCategory) []ShopCategory {
	categoryMap := make(map[int64]*ShopCategory)
	var roots []ShopCategory

	// 先建立映射关系
	for i := range categories {
		categoryMap[categories[i].ID] = &categories[i]
	}

	// 构建树形结构
	for i := range categories {
		if categories[i].ParentId == 0 {
			// 这是根节点
			roots = append(roots, categories[i])
		} else {
			// 找到父节点，添加到其子节点列表中
			if parent, ok := categoryMap[categories[i].ParentId]; ok {
				parent.Children = append(parent.Children, categories[i])
			}
		}
	}

	return roots
}

// Find 查询分类列表
func (s *ShopCategory) Find(dest interface{}) error {
	return s.Query().Find(dest).Error
}

// First 获取单个分类
func (s *ShopCategory) First() error {
	return s.Query().First(s).Error
}
