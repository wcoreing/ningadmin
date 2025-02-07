package shop

import (
	"time"
	"wnadmin/app/provider/db"

	"gorm.io/gorm"
)

// ShopCategory 商品分类
type ShopCategory struct {
	ID         int64          `gorm:"column:id;primaryKey;autoIncrement" json:"id"`         // 分类ID
	ParentId   int64          `gorm:"column:parent_id;not null;default:0" json:"parent_id"` // 父级ID
	Name       string         `gorm:"column:name;not null" json:"name"`                     // 分类名称
	Icon       string         `gorm:"column:icon" json:"icon"`                              // 分类图标
	Image      string         `gorm:"column:image" json:"image"`                            // 分类图片
	Level      int            `gorm:"column:level;not null;default:1" json:"level"`         // 层级
	Sort       int            `gorm:"column:sort;not null;default:0" json:"sort"`           // 排序
	Status     int            `gorm:"column:status;not null;default:1" json:"status"`       // 状态
	IsDelete   int            `gorm:"column:is_delete;not null;default:0" json:"is_delete"` // 是否删除
	CreateTime time.Time      `gorm:"column:create_time;autoCreateTime" json:"create_time"` // 创建时间
	UpdateTime time.Time      `gorm:"column:update_time;autoUpdateTime" json:"update_time"` // 更新时间
	Children   []ShopCategory `json:"children" gorm:"-"`                                    // 子分类列表
}

// TableName 表名
func (s *ShopCategory) TableName() string {
	return "shop_category"
}

// Query 查询构造器
func (s *ShopCategory) Query() *gorm.DB {
	return db.DB.WnAdmin.Model(s)
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
