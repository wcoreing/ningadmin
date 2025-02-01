package shop

import (
	"time"
	"wnadmin/app/provider/db"

	"gorm.io/gorm"
)

// ShopBrand 商品品牌
type ShopBrand struct {
	ID          int64     `gorm:"primaryKey;column:id" json:"id"`
	Name        string    `gorm:"column:name" json:"name"`               // 品牌名称
	Logo        string    `gorm:"column:logo" json:"logo"`               // 品牌logo
	Description string    `gorm:"column:description" json:"description"` // 品牌描述
	Website     string    `gorm:"column:website" json:"website"`         // 品牌官网
	Sort        int       `gorm:"column:sort" json:"sort"`               // 排序
	Status      int8      `gorm:"column:status" json:"status"`           // 状态(0:禁用 1:启用)
	IsDelete    int8      `gorm:"column:is_delete" json:"isDelete"`      // 是否删除
	CreateTime  time.Time `gorm:"column:create_time" json:"createTime"`  // 创建时间
	UpdateTime  time.Time `gorm:"column:update_time" json:"updateTime"`  // 更新时间
}

// TableName 表名
func (s *ShopBrand) TableName() string {
	return "shop_brand"
}

// Query 获取查询构建器
func (s *ShopBrand) Query() *gorm.DB {
	query := db.DB.WnAdmin.Model(s).Where("is_delete = ?", 0)

	// 品牌名称
	if s.Name != "" {
		query = query.Where("name LIKE ?", "%"+s.Name+"%")
	}

	// 状态
	if s.Status > 0 {
		query = query.Where("status = ?", s.Status)
	}

	return query.Order("sort DESC, id DESC")
}

// Find 查询品牌列表
func (s *ShopBrand) Find(dest interface{}) error {
	return s.Query().Find(dest).Error
}

// First 获取单个品牌
func (s *ShopBrand) First() error {
	return s.Query().First(s).Error
}
