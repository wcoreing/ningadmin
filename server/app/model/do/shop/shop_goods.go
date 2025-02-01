package shop

import (
	"fmt"
	"time"
	"wnadmin/app/model/dto"
	"wnadmin/app/provider/db"

	"gorm.io/gorm"
)

// GoodsDO 商品信息
type GoodsDO struct {
	Id            int64     `json:"goodsId" gorm:"column:id;primaryKey;comment:商品ID"`
	CategoryId    int64     `json:"categoryId" gorm:"column:category_id;comment:分类ID"`
	SupplierId    int64     `json:"supplierId" gorm:"column:supplier_id;comment:供应商ID"`
	BrandId       int64     `json:"brandId" gorm:"column:brand_id;comment:品牌ID"`
	Type          int8      `json:"type" gorm:"column:type;default:1;comment:商品类型1普通商品2虚拟商品"`
	Name          string    `json:"name" gorm:"column:name;comment:商品名称"`
	Description   string    `json:"description" gorm:"column:description;comment:商品描述"`
	Detail        string    `json:"detail" gorm:"column:detail;comment:商品详情"`
	Price         float64   `json:"price" gorm:"column:price;comment:售价"`
	OriginalPrice float64   `json:"originalPrice" gorm:"column:original_price;comment:原价"`
	CostPrice     float64   `json:"costPrice" gorm:"column:cost_price;comment:成本价"`
	Stock         int       `json:"stock" gorm:"column:stock;default:0;comment:库存"`
	WarnStock     int       `json:"warnStock" gorm:"column:warn_stock;default:0;comment:库存预警值"`
	Unit          string    `json:"unit" gorm:"column:unit;comment:单位"`
	Weight        float64   `json:"weight" gorm:"column:weight;default:0;comment:重量(kg)"`
	Volume        float64   `json:"volume" gorm:"column:volume;default:0;comment:体积(m³)"`
	Sales         int       `json:"sales" gorm:"column:sales;default:0;comment:销量"`
	VirtualSales  int       `json:"virtualSales" gorm:"column:virtual_sales;default:0;comment:虚拟销量"`
	SpecType      int8      `json:"specType" gorm:"column:spec_type;default:1;comment:规格类型：1-单规格 2-多规格"`
	Status        int8      `json:"status" gorm:"column:status;default:0;comment:状态(0:下架 1:上架)"`
	Sort          int       `json:"sort" gorm:"column:sort;default:0;comment:排序"`
	IsDelete      int8      `json:"-" gorm:"column:is_delete;default:0;comment:是否删除"`
	CreateTime    time.Time `json:"createTime" gorm:"column:create_time;autoCreateTime;comment:创建时间"`
	UpdateTime    time.Time `json:"updateTime" gorm:"column:update_time;autoUpdateTime;comment:更新时间"`
}

// TableName 表名
func (GoodsDO) TableName() string {
	return "shop_goods"
}

// Query 获取查询构建器
func (g *GoodsDO) Query() *gorm.DB {
	query := db.DB.WnAdmin.Model(g).Where("is_delete = ?", 0)

	// 添加调试日志
	fmt.Printf("构建查询条件: %+v\n", g)

	// 根据结构体中的非零值构建查询条件
	if g.CategoryId > 0 {
		query = query.Where("category_id = ?", g.CategoryId)
	}
	if g.SupplierId > 0 {
		query = query.Where("supplier_id = ?", g.SupplierId)
	}
	if g.BrandId > 0 {
		query = query.Where("brand_id = ?", g.BrandId)
	}
	if g.Status >= 0 && g.Status != 0 {
		query = query.Where("status = ?", g.Status)
	}

	return query
}

// Create 创建商品
func (g *GoodsDO) Create() error {
	return g.Query().Create(g).Error
}

// Update 更新商品
func (g *GoodsDO) Update() error {
	return g.Query().Updates(g).Error
}

// Delete 删除商品
func (g *GoodsDO) Delete() error {
	return g.Query().Where("goods_id = ?", g.Id).Update("is_delete", 1).Error
}

// First 获取单个商品
func (g *GoodsDO) First() error {
	return g.Query().First(g).Error
}

// Find 查询商品列表
func (g *GoodsDO) Find(dest interface{}) error {
	return g.Query().Find(dest).Error
}

// UpdateStatus 更新商品状态
func (g *GoodsDO) UpdateStatus() error {
	return g.Query().Where("goods_id = ?", g.Id).Update("status", g.Status).Error
}

// GetGoodsList 获取商品列表
func (g *GoodsDO) GetGoodsList(query *dto.GoodsListReq) ([]GoodsDO, int64, error) {
	var goods []GoodsDO
	var total int64

	q := db.DB.WnAdmin.Model(g).Where("is_delete = ?", 0)

	// 标签查询
	if query.TagId > 0 {
		q = q.Joins("JOIN shop_goods_tag_relation rel ON rel.goods_id = shop_goods.id").
			Where("rel.tag_id = ?", query.TagId)
	}

	// 其他查询条件
	if query.CategoryId > 0 {
		q = q.Where("category_id = ?", query.CategoryId)
	}
	if query.SupplierId > 0 {
		q = q.Where("supplier_id = ?", query.SupplierId)
	}
	if query.BrandId > 0 {
		q = q.Where("brand_id = ?", query.BrandId)
	}
	if query.Status != "" {
		q = q.Where("status = ?", query.Status)
	}
	if query.Name != "" {
		q = q.Where("name LIKE ?", "%"+query.Name+"%")
	}

	// 统计总数
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err := q.Offset((query.PageNum - 1) * query.PageSize).
		Limit(query.PageSize).
		Order("create_time DESC").
		Find(&goods).Error

	return goods, total, err
}
