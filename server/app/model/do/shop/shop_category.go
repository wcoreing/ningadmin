package shop

import (
	"time"
	"wnadmin/app/provider/db"

	"gorm.io/gorm"
)

// ShopCategory 商品分类
type ShopCategory struct {
	ID         int64     `gorm:"primaryKey;column:id" json:"id"`
	ParentID   int64     `gorm:"column:parent_id" json:"parentId"`     // 父级ID
	Name       string    `gorm:"column:name" json:"name"`              // 分类名称
	Icon       string    `gorm:"column:icon" json:"icon"`              // 分类图标
	Image      string    `gorm:"column:image" json:"image"`            // 分类图片
	Level      int8      `gorm:"column:level" json:"level"`            // 层级
	Sort       int       `gorm:"column:sort" json:"sort"`              // 排序
	Status     int8      `gorm:"column:status" json:"status"`          // 状态
	IsDelete   int8      `gorm:"column:is_delete" json:"isDelete"`     // 是否删除
	CreateTime time.Time `gorm:"column:create_time" json:"createTime"` // 创建时间
	UpdateTime time.Time `gorm:"column:update_time" json:"updateTime"` // 更新时间
}

// TableName 表名
func (s *ShopCategory) TableName() string {
	return "shop_category"
}

// Query 获取查询构建器
func (s *ShopCategory) Query() *gorm.DB {
	query := db.DB.WnAdmin.Model(s).Where("is_delete = ?", 0)

	// 分类名称
	if s.Name != "" {
		query = query.Where("name LIKE ?", "%"+s.Name+"%")
	}

	// 父级ID
	if s.ParentID > 0 {
		query = query.Where("parent_id = ?", s.ParentID)
	}

	// 层级
	if s.Level > 0 {
		query = query.Where("level = ?", s.Level)
	}

	// 状态
	if s.Status > 0 {
		query = query.Where("status = ?", s.Status)
	}

	return query.Order("sort DESC, id DESC")
}

// Find 查询分类列表
func (s *ShopCategory) Find(dest interface{}) error {
	return s.Query().Find(dest).Error
}

// First 获取单个分类
func (s *ShopCategory) First() error {
	return s.Query().First(s).Error
}
