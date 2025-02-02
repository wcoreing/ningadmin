package shop

import (
	"time"
	"wnadmin/app/provider/db"

	"gorm.io/gorm"
)

// GoodsImageDO 商品图片
type GoodsImageDO struct {
	Id         int64     `json:"id" gorm:"column:id;primaryKey;comment:图片ID"`
	GoodsId    int64     `json:"goodsId" gorm:"column:goods_id;comment:商品ID"`
	IsMain     int8      `json:"isMain" gorm:"column:is_main;default:0;comment:是否主图(0:否 1:是)"`
	Url        string    `json:"url" gorm:"column:url;comment:图片地址"`
	Type       int8      `json:"type" gorm:"column:type;comment:图片类型(1:商品图片 2:详情图片)"`
	Sort       int       `json:"sort" gorm:"column:sort;default:0;comment:排序"`
	CreateTime time.Time `json:"createTime" gorm:"column:create_time;autoCreateTime;comment:创建时间"`
	UpdateTime time.Time `json:"updateTime" gorm:"column:update_time;autoUpdateTime;comment:更新时间"`
}

// TableName 表名
func (GoodsImageDO) TableName() string {
	return "shop_goods_image"
}

// Query 获取查询构建器
func (g *GoodsImageDO) Query() *gorm.DB {
	return db.DB.WnAdmin.Model(g)
}
