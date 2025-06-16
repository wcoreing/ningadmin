package shop

import (
	"time"
	"wnadmin/app/provider/db"

	"gorm.io/gorm"
)

// Category 通用分类表
type Category struct {
	ID         int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`         // 分类ID
	ParentId   int64      `gorm:"column:parent_id;not null;default:0" json:"parent_id"` // 父级ID
	Name       string     `gorm:"column:name;not null" json:"name"`                     // 分类名称
	Icon       string     `gorm:"column:icon" json:"icon"`                              // 分类图标
	Image      string     `gorm:"column:image" json:"image"`                            // 分类图片
	Level      int        `gorm:"column:level;not null;default:1" json:"level"`         // 层级
	Sort       int        `gorm:"column:sort;not null;default:0" json:"sort"`           // 排序
	Status     int        `gorm:"column:status;not null;default:1" json:"status"`       // 状态
	Module     string     `gorm:"column:module;not null;default:goods" json:"module"`   // 模块标识
	IsDelete   int        `gorm:"column:is_delete;not null;default:0" json:"is_delete"` // 是否删除
	CreateTime time.Time  `gorm:"column:create_time;autoCreateTime" json:"create_time"` // 创建时间
	UpdateTime time.Time  `gorm:"column:update_time;autoUpdateTime" json:"update_time"` // 更新时间
	Children   []Category `json:"children" gorm:"-"`                                    // 子分类列表
}

// TableName 表名
func (c *Category) TableName() string {
	return "category"
}

// Query 查询构造器
func (c *Category) Query() *gorm.DB {
	return db.WnAdmin().Model(c).Where("is_delete = 0")
}

// BuildTree 构建分类树
func (c *Category) BuildTree(module string) ([]Category, error) {
	var categories []Category

	// 查询指定模块的所有分类
	if err := c.Query().Where("module = ?", module).Order("sort DESC, id ASC").Find(&categories).Error; err != nil {
		return nil, err
	}

	// 构建树形结构
	categoryMap := make(map[int64]*Category)
	var rootCategories []Category

	// 先将所有分类放入map
	for i := range categories {
		categoryMap[categories[i].ID] = &categories[i]
	}

	// 构建父子关系
	for i := range categories {
		if categories[i].ParentId == 0 {
			// 根分类
			rootCategories = append(rootCategories, categories[i])
		} else {
			// 子分类
			if parent, exists := categoryMap[categories[i].ParentId]; exists {
				parent.Children = append(parent.Children, categories[i])
			}
		}
	}

	return rootCategories, nil
}

// GetByModule 根据模块获取分类列表
func (c *Category) GetByModule(module string) ([]Category, error) {
	var categories []Category
	err := c.Query().Where("module = ?", module).Order("sort DESC, id ASC").Find(&categories).Error
	return categories, err
}
