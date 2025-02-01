package shop

import (
	"time"
	"wnadmin/app/provider/db"

	"gorm.io/gorm"
)

// ShopGoodsTagRel 商品标签关联
type ShopGoodsTagRel struct {
	ID         int64     `gorm:"primaryKey;column:id" json:"id"`
	GoodsId    int64     `gorm:"column:goods_id" json:"goodsId"`       // 商品ID
	TagId      int64     `gorm:"column:tag_id" json:"tagId"`           // 标签ID
	CreateTime time.Time `gorm:"column:create_time" json:"createTime"` // 创建时间
}

// TableName 表名
func (ShopGoodsTagRel) TableName() string {
	return "shop_goods_tag_relation"
}

// Query 获取查询构建器
func (s *ShopGoodsTagRel) Query() *gorm.DB {
	return db.DB.WnAdmin.Model(s)
}
