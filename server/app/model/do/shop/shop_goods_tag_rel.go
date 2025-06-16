package shop

import (
	"time"
	"wnadmin/app/provider/db"

	"gorm.io/gorm"
)

// ShopGoodsTagRel 商品标签关联
type ShopGoodsTagRel struct {
	Id         int64     `json:"id" gorm:"column:id;primaryKey;comment:ID"`
	GoodsId    int64     `json:"goodsId" gorm:"column:goods_id;comment:商品ID"`
	TagName    string    `json:"tagName" gorm:"column:tag_name;comment:标签名称"`
	CreateTime time.Time `json:"createTime" gorm:"column:create_time;autoCreateTime;comment:创建时间"`
}

// TableName 表名
func (ShopGoodsTagRel) TableName() string {
	return "shop_goods_tag_relation"
}

// Query 获取查询构建器
func (g *ShopGoodsTagRel) Query() *gorm.DB {
	return db.WnAdmin().Model(g)
}
