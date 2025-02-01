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
	GoodsId       int64   `json:"goodsId" binding:"required"`     // 商品ID
	GoodsName     string  `json:"goodsName" binding:"required"`   // 商品名称
	CategoryId    int64   `json:"categoryId" binding:"required"`  // 分类ID
	BrandId       int64   `json:"brandId" binding:"required"`     // 品牌ID
	SupplierId    int64   `json:"supplierId" binding:"required"`  // 供应商ID
	Price         float64 `json:"price" binding:"required,min=0"` // 价格
	OriginalPrice float64 `json:"originalPrice"`                  // 原价
	Stock         int     `json:"stock" binding:"required,min=0"` // 库存
	Sort          int     `json:"sort" binding:"min=0"`           // 排序
	Status        string  `json:"status" binding:"required"`      // 状态
	Description   string  `json:"description"`                    // 商品描述
	TagIds        []int64 `json:"tagIds"`                         // 标签ID列表
	CoverImg      string  `json:"coverImg"`                       // 封面图
}

// GoodsInfoResp 商品详情响应
type GoodsInfoResp struct {
	GoodsId       int64   `json:"goodsId"`       // 商品ID
	GoodsName     string  `json:"goodsName"`     // 商品名称
	CategoryId    int64   `json:"categoryId"`    // 分类ID
	BrandId       int64   `json:"brandId"`       // 品牌ID
	SupplierId    int64   `json:"supplierId"`    // 供应商ID
	Price         float64 `json:"price"`         // 价格
	OriginalPrice float64 `json:"originalPrice"` // 原价
	Stock         int     `json:"stock"`         // 库存
	Sort          int     `json:"sort"`          // 排序
	Status        string  `json:"status"`        // 状态
	Description   string  `json:"description"`   // 商品描述
	CoverImg      string  `json:"coverImg"`      // 封面图
	TagIds        []int64 `json:"tagIds"`        // 标签ID列表
}
