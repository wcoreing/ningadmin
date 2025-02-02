package shop

import (
	"fmt"
	"wnadmin/app/model/dto"
	"wnadmin/app/provider/db"

	"time"

	"gorm.io/gorm"
)

// GoodsDO 商品表
type GoodsDO struct {
	Id            int64     `gorm:"column:id" json:"id"`                        // 商品ID
	Name          string    `gorm:"column:name" json:"name"`                    // 商品名称
	CategoryId    int64     `gorm:"column:category_id" json:"categoryId"`       // 分类ID
	BrandId       int64     `gorm:"column:brand_id" json:"brandId"`             // 品牌ID
	SupplierId    int64     `gorm:"column:supplier_id" json:"supplierId"`       // 供应商ID
	Type          int8      `gorm:"column:type" json:"type"`                    // 商品类型1普通商品2虚拟商品
	Price         float64   `gorm:"column:price" json:"price"`                  // 销售价
	OriginalPrice float64   `gorm:"column:original_price" json:"originalPrice"` // 市场价
	CostPrice     float64   `gorm:"column:cost_price" json:"costPrice"`         // 成本价
	Stock         int       `gorm:"column:stock" json:"stock"`                  // 库存
	WarnStock     int       `gorm:"column:warn_stock" json:"warnStock"`         // 库存预警值
	Unit          string    `gorm:"column:unit" json:"unit"`                    // 单位
	Weight        float64   `gorm:"column:weight" json:"weight"`                // 重量(kg)
	Volume        float64   `gorm:"column:volume" json:"volume"`                // 体积(m³)
	Sales         int       `gorm:"column:sales" json:"sales"`                  // 销量
	VirtualSales  int       `gorm:"column:virtual_sales" json:"virtualSales"`   // 虚拟销量
	Sort          int       `gorm:"column:sort" json:"sort"`                    // 排序
	Status        int8      `gorm:"column:status" json:"status"`                // 状态
	Description   string    `gorm:"column:description" json:"description"`      // 描述
	Detail        string    `gorm:"column:detail" json:"detail"`                // 商品详情
	SpecType      int8      `gorm:"column:spec_type" json:"specType"`           // 规格类型
	IsDelete      int8      `gorm:"column:is_delete" json:"isDelete"`           // 是否删除
	CreateTime    time.Time `gorm:"column:create_time" json:"createTime"`       // 创建时间
	UpdateTime    time.Time `gorm:"column:update_time" json:"updateTime"`       // 更新时间

	// 关联字段，不映射到数据库
	Images     []*GoodsImageDO  `gorm:"-" json:"images"`     // 商品图片列表
	SpecGroups []*ShopGoodsSpec `gorm:"-" json:"specGroups"` // 规格组列表
	SkuList    []*ShopGoodsSku  `gorm:"-" json:"skuList"`    // SKU列表
	TagIds     []int64          `gorm:"-" json:"tagIds"`     // 标签ID列表
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
	return g.Query().Where("id = ?", g.Id).Update("is_delete", 1).Error
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
	return g.Query().Where("id = ?", g.Id).Update("status", g.Status).Error
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
