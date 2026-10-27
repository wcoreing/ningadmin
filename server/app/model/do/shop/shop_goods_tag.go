package shop

import (
	"time"
	"wnadmin/app/provider/db"

	"gorm.io/gorm"
)

// ShopGoodsTag 商品标签
type ShopGoodsTag struct {
	ID         int64     `gorm:"primaryKey;column:id" json:"id"`
	Name       string    `gorm:"column:name" json:"name"`              // 标签名称
	Sort       int       `gorm:"column:sort" json:"sort"`              // 排序
	IsDelete   int8      `gorm:"column:is_delete" json:"isDelete"`     // 是否删除
	CreateTime time.Time `gorm:"column:create_time" json:"createTime"` // 创建时间
	UpdateTime time.Time `gorm:"column:update_time" json:"updateTime"` // 更新时间
}

// TableName 表名
func (s *ShopGoodsTag) TableName() string {
	return "shop_goods_tag"
}

// Query 获取查询构建器
func (s *ShopGoodsTag) Query() *gorm.DB {
	query := db.WnAdmin().Model(s).Where("is_delete = ?", 0)

	// 标签名称
	if s.Name != "" {
		query = query.Where("name LIKE ?", "%"+s.Name+"%")
	}

	return query.Order("sort DESC, id DESC")
}

// Find 查询标签列表
func (s *ShopGoodsTag) Find(dest interface{}) error {
	return s.Query().Find(dest).Error
}

// First 获取单个标签
func (s *ShopGoodsTag) First() error {
	return s.Query().First(s).Error
}
