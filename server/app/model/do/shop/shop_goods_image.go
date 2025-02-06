package shop

import (
	"time"
	"wnadmin/app/provider/db"

	"gorm.io/gorm"
)

// 图片位置常量
const (
	ImagePositionNormal = 0 // 普通图片
	ImagePositionMain   = 1 // 主图
	ImagePositionList   = 2 // 列表图
	ImagePositionDetail = 3 // 详情主图
	ImagePositionShare  = 4 // 分享图
)

// ShopGoodsImage 商品图片表
type ShopGoodsImage struct {
	Id         int64     `gorm:"column:id;primary_key;auto_increment" json:"id"`      // 图片ID
	GoodsId    int64     `gorm:"column:goods_id;not null" json:"goodsId"`             // 商品ID
	Position   uint8     `gorm:"column:position;not null;default:0" json:"position"`  // 图片位置
	Url        string    `gorm:"column:url;not null" json:"url"`                      // 图片地址
	Sort       int       `gorm:"column:sort;not null;default:0" json:"sort"`          // 排序
	IsDelete   uint8     `gorm:"column:is_delete;not null;default:0" json:"isDelete"` // 是否删除
	CreateTime time.Time `gorm:"column:create_time;autoCreateTime" json:"createTime"` // 创建时间
	UpdateTime time.Time `gorm:"column:update_time;autoUpdateTime" json:"updateTime"` // 更新时间
}

// TableName 表名
func (s *ShopGoodsImage) TableName() string {
	return "shop_goods_image"
}

// Query 获取查询器
func (s *ShopGoodsImage) Query() *gorm.DB {
	return db.DB.WnAdmin.Model(s)
}
