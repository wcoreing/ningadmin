package dto

// GoodsListReq 商品列表请求参数
type GoodsListReq struct {
	PageNum     int     `json:"pageNum" form:"pageNum"`       // 页码
	PageSize    int     `json:"pageSize" form:"pageSize"`     // 每页数量
	Name        string  `json:"name" form:"name"`             // 商品名称
	CategoryId  int64   `json:"categoryId" form:"categoryId"` // 分类ID
	CategoryIds []int64 `json:"-" form:"-"`                   // 分类ID列表(内部使用)
	SupplierId  int64   `json:"supplierId" form:"supplierId"` // 供应商ID
	BrandId     int64   `json:"brandId" form:"brandId"`       // 品牌ID
	TagId       int64   `json:"tagId" form:"tagId"`           // 标签ID
	GoodsId     int64   `json:"-" form:"-"`                   // 商品ID(内部使用)
	Status      string  `form:"status"`
}

// GoodsDTO 商品信息DTO
type GoodsDTO struct {
	GoodsId       int64    `json:"goodsId"`
	CategoryId    int64    `json:"categoryId"`
	CategoryName  string   `json:"categoryName"`
	SupplierId    int64    `json:"supplierId"`
	SupplierName  string   `json:"supplierName"`
	BrandId       int64    `json:"brandId"`
	BrandName     string   `json:"brandName"`
	TagId         int64    `json:"tagId"`
	TagName       string   `json:"tagName"`
	Name          string   `json:"name"`
	Price         float64  `json:"price"`
	OriginalPrice float64  `json:"originalPrice"`
	Stock         int      `json:"stock"`
	Sales         int      `json:"sales"`
	Status        int8     `json:"status"`
	Sort          int      `json:"sort"`
	TagIds        []int64  `json:"tagIds"`   // 标签ID列表
	TagNames      []string `json:"tagNames"` // 标签名称列表
	CoverImg      string   `json:"coverImg"`
}

// GoodsStatusReq 商品状态更新请求
type GoodsStatusReq struct {
	Id     int64  `json:"id" binding:"required"`
	Status string `json:"status" binding:"required"`
}

// GoodsImageUpdateItem 商品图片更新项
type GoodsImageUpdateItem struct {
	Url      string `json:"url"`      // 图片URL
	Position uint8  `json:"position"` // 图片位置
	Sort     int    `json:"sort"`     // 排序
}

// GoodsUpdateReq 商品更新请求
type GoodsUpdateReq struct {
	GoodsId       int64                   `json:"goodsId,omitempty"`
	GoodsName     string                  `json:"goodsName" binding:"required"`
	CategoryId    int64                   `json:"categoryId" binding:"required"`
	BrandId       int64                   `json:"brandId" binding:"required"`
	SupplierId    int64                   `json:"supplierId" binding:"required"`
	Price         float64                 `json:"price"`
	OriginalPrice float64                 `json:"originalPrice"`
	Stock         int                     `json:"stock"`
	Sort          int                     `json:"sort"`
	Status        string                  `json:"status"`
	Description   string                  `json:"description"`
	CoverImg      string                  `json:"coverImg"`
	TagIds        []int64                 `json:"tagIds"`
	Images        []*GoodsImageUpdateItem `json:"images"`
	SpecType      int                     `json:"specType"`   // 规格类型：1-单规格，2-多规格
	SpecGroups    []SpecGroup             `json:"specGroups"` // 规格组
	SkuList       []Sku                   `json:"skuList"`    // SKU列表
	DetailDesc    string                  `json:"detailDesc"` // 商品详情
}

// SpecGroup 规格组
type SpecGroup struct {
	ID           int64       `json:"id"`           // 规格组ID
	Name         string      `json:"name"`         // 规格组名称
	Values       []SpecValue `json:"values"`       // 规格值列表
	Sort         int         `json:"sort"`         // 排序
	InputVisible bool        `json:"inputVisible"` // 输入框是否可见
	InputValue   string      `json:"inputValue"`   // 输入框的值
}

// Sku SKU信息
type Sku struct {
	ID    int64                `json:"id,omitempty"`
	Specs map[string]SpecValue `json:"specs"` // 规格信息，key为规格组名称
	Price float64              `json:"price"` // 价格
	Stock int                  `json:"stock"` // 库存
	Image string               `json:"image"` // SKU图片
}

// SpecValue 规格值
type SpecValue struct {
	ID    int64  `json:"id"`    // 规格值ID
	Name  string `json:"name"`  // 规格值名称
	Value string `json:"value"` // 规格值
	Sort  int    `json:"sort"`  // 排序
}

// GoodsInfoResp 商品详情响应
type GoodsInfoResp struct {
	GoodsId       int64                   `json:"goodsId"`
	GoodsName     string                  `json:"goodsName"`
	CategoryId    int64                   `json:"categoryId"`
	BrandId       int64                   `json:"brandId"`
	SupplierId    int64                   `json:"supplierId"`
	Type          int8                    `json:"type"`
	Price         float64                 `json:"price"`
	OriginalPrice float64                 `json:"originalPrice"`
	CostPrice     float64                 `json:"costPrice"`
	Stock         int                     `json:"stock"`
	WarnStock     int                     `json:"warnStock"`
	Unit          string                  `json:"unit"`
	Weight        float64                 `json:"weight"`
	Volume        float64                 `json:"volume"`
	Sort          int                     `json:"sort"`
	Status        string                  `json:"status"`
	Description   string                  `json:"description"`
	CoverImg      string                  `json:"coverImg"`
	TagIds        []int64                 `json:"tagIds"`
	Images        []*GoodsImageUpdateItem `json:"images"`
	SpecType      int8                    `json:"specType"`   // 规格类型：1-单规格，2-多规格
	SpecGroups    []SpecGroup             `json:"specGroups"` // 规格组
	SkuList       []Sku                   `json:"skuList"`    // SKU列表
	DetailDesc    string                  `json:"detailDesc"` // 商品详情
}
