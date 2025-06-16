package shop

import (
	"fmt"
	"wnadmin/app/model/dto"
	"wnadmin/app/provider/db"

	"time"

	"gorm.io/gorm"
)

// 商品类型常量
const (
	GoodsTypeNormal  = 1 // 实物商品
	GoodsTypeVirtual = 2 // 虚拟商品
	GoodsTypeCombo   = 3 // 套餐商品
)

// 规格类型常量
const (
	SpecTypeSingle = 1 // 单规格
	SpecTypeMulti  = 2 // 多规格
	SpecTypeCombo  = 3 // 套餐规格
)

// GoodsDO 商品表
type GoodsDO struct {
	Id          int64     `gorm:"column:id" json:"id"`                   // 商品ID
	Name        string    `gorm:"column:name" json:"name"`               // 商品名称
	CategoryId  int64     `gorm:"column:category_id" json:"categoryId"`  // 分类ID
	Brand       string    `gorm:"column:brand" json:"brand"`             // 品牌名称
	Supplier    string    `gorm:"column:supplier" json:"supplier"`       // 供应商名称
	Type        int8      `gorm:"column:type" json:"type"`               // 商品类型：1-实物商品 2-虚拟商品 3-套餐商品
	Unit        string    `gorm:"column:unit" json:"unit"`               // 单位
	Sort        int       `gorm:"column:sort" json:"sort"`               // 排序
	Status      int8      `gorm:"column:status" json:"status"`           // 状态
	Description string    `gorm:"column:description" json:"description"` // 描述
	SpecType    int8      `gorm:"column:spec_type" json:"specType"`      // 规格类型
	IsDelete    int8      `gorm:"column:is_delete" json:"isDelete"`      // 是否删除
	CreateTime  time.Time `gorm:"column:create_time" json:"createTime"`  // 创建时间
	UpdateTime  time.Time `gorm:"column:update_time" json:"updateTime"`  // 更新时间

	CoverImg string `gorm:"-" json:"coverImg"`

	// 关联字段，不映射到数据库
	Images     []*ShopGoodsImage `gorm:"-" json:"images"`     // 商品图片列表
	SpecGroups []*ShopGoodsSpec  `gorm:"-" json:"specGroups"` // 规格组列表
	SkuList    []*ShopGoodsSku   `gorm:"-" json:"skuList"`    // SKU列表
	TagNames   []string          `gorm:"-" json:"tagNames"`   // 标签名称列表
}

// TableName 表名
func (GoodsDO) TableName() string {
	return "shop_goods"
}

// IsNormalGoods 判断是否为普通商品
func (g *GoodsDO) IsNormalGoods() bool {
	return g.Type == GoodsTypeNormal
}

// IsVirtualGoods 判断是否为虚拟商品
func (g *GoodsDO) IsVirtualGoods() bool {
	return g.Type == GoodsTypeVirtual
}

// Query 获取查询构建器
func (g *GoodsDO) Query() *gorm.DB {
	query := db.WnAdmin().Model(g).Where("is_delete = ?", 0)

	// 添加调试日志
	fmt.Printf("构建查询条件: %+v\n", g)

	// 根据结构体中的非零值构建查询条件
	if g.CategoryId > 0 {
		query = query.Where("category_id = ?", g.CategoryId)
	}

	if g.Status >= 0 && g.Status != 0 {
		query = query.Where("status = ?", g.Status)
	}

	// 商品类型筛选
	if g.Type > 0 {
		query = query.Where("type = ?", g.Type)
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
func (g *GoodsDO) GetGoodsList(req *dto.GoodsListReq) ([]GoodsDO, int64, error) {
	var (
		list  []GoodsDO
		total int64
		db    = g.Query()
	)

	// 构建查询条件
	if req.Name != "" {
		db = db.Where("name LIKE ?", "%"+req.Name+"%")
	}
	if len(req.CategoryIds) > 0 {
		db = db.Where("category_id IN ?", req.CategoryIds)
	}

	if req.Status != "" {
		db = db.Where("status = ?", req.Status)
	}

	// 商品类型筛选
	if req.Type > 0 {
		db = db.Where("type = ?", req.Type)
	}

	// 统计总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err := db.Offset((req.PageNum - 1) * req.PageSize).
		Limit(req.PageSize).
		Order("create_time DESC").
		Find(&list).Error

	return list, total, err
}

type ShopGoods struct {
	Id       uint    `json:"id" gorm:"primaryKey"`
	Name     string  `json:"name"`
	CoverImg string  `json:"coverImg" gorm:"column:url"`
	Price    float64 `json:"price"`
}
