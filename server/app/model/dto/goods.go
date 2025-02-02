package dto

// GoodsListReq 商品列表请求参数
type GoodsListReq struct {
	PageNum    int    `form:"pageNum"`
	PageSize   int    `form:"pageSize"`
	Name       string `form:"name"`
	CategoryId int64  `form:"categoryId"`
	SupplierId int64  `form:"supplierId"`
	BrandId    int64  `form:"brandId"`
	TagId      int64  `form:"tagId"`
	Status     string `form:"status"`
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
}

// GoodsStatusReq 商品状态更新请求
type GoodsStatusReq struct {
	Id     int64  `json:"id" binding:"required"`
	Status string `json:"status" binding:"required"`
}

// GoodsUpdateReq 商品更新请求
type GoodsUpdateReq struct {
	GoodsId       int64       `json:"goodsId"`       // 商品ID
	GoodsName     string      `json:"goodsName"`     // 商品名称
	CategoryId    int64       `json:"categoryId"`    // 分类ID
	BrandId       int64       `json:"brandId"`       // 品牌ID
	SupplierId    int64       `json:"supplierId"`    // 供应商ID
	Price         float64     `json:"price"`         // 销售价
	OriginalPrice float64     `json:"originalPrice"` // 市场价
	Stock         int         `json:"stock"`         // 库存
	Sort          int         `json:"sort"`          // 排序
	Status        string      `json:"status"`        // 商品状态 0-上架 1-下架
	Description   string      `json:"description"`   // 商品描述
	CoverImg      string      `json:"coverImg"`      // 封面图
	TagIds        []int64     `json:"tagIds"`        // 标签ID列表
	Images        []string    `json:"images"`        // 商品图片列表
	Detail        string      `json:"detail"`        // 商品详情
	SpecType      int         `json:"specType"`      // 规格类型 0-单规格 1-多规格
	SpecGroups    []SpecGroup `json:"specGroups"`    // 规格组
	SkuList       []Sku       `json:"skuList"`       // SKU列表
	IsDraft       bool        `json:"isDraft"`       // 是否为草稿
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

// SpecValue 规格值
type SpecValue struct {
	ID    int64  `json:"id"`    // 规格值ID
	Name  string `json:"name"`  // 规格值名称
	Value string `json:"value"` // 规格值
	Sort  int    `json:"sort"`  // 排序
}

// Sku SKU信息
type Sku struct {
	ID    int64                `json:"id"`    // SKU ID
	Specs map[string]SpecValue `json:"specs"` // 规格信息，key是规格名称，value是规格值对象
	Price float64              `json:"price"` // 价格
	Stock int                  `json:"stock"` // 库存
}

// GoodsInfoResp 商品详情响应
type GoodsInfoResp struct {
	GoodsId       int64       `json:"goodsId"`       // 商品ID
	GoodsName     string      `json:"goodsName"`     // 商品名称
	CategoryId    int64       `json:"categoryId"`    // 分类ID
	BrandId       int64       `json:"brandId"`       // 品牌ID
	SupplierId    int64       `json:"supplierId"`    // 供应商ID
	Type          int8        `json:"type"`          // 商品类型
	Price         float64     `json:"price"`         // 价格
	OriginalPrice float64     `json:"originalPrice"` // 原价
	CostPrice     float64     `json:"costPrice"`     // 成本价
	Stock         int         `json:"stock"`         // 库存
	WarnStock     int         `json:"warnStock"`     // 库存预警值
	Unit          string      `json:"unit"`          // 单位
	Weight        float64     `json:"weight"`        // 重量
	Volume        float64     `json:"volume"`        // 体积
	Sort          int         `json:"sort"`          // 排序
	Status        string      `json:"status"`        // 状态
	Description   string      `json:"description"`   // 商品描述
	Detail        string      `json:"detail"`        // 商品详情
	CoverImg      string      `json:"coverImg"`      // 封面图
	TagIds        []int64     `json:"tagIds"`        // 标签ID列表
	SpecType      int8        `json:"specType"`      // 规格类型
	SpecGroups    []SpecGroup `json:"specGroups"`    // 规格组
	SkuList       []Sku       `json:"skuList"`       // SKU列表
}
