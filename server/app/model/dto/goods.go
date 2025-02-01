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
