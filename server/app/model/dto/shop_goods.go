package dto

// GoodsListReq 商品列表请求
type GoodsListReq struct {
	PageNum    int    `json:"pageNum" form:"pageNum"`
	PageSize   int    `json:"pageSize" form:"pageSize"`
	GoodsName  string `json:"goodsName" form:"goodsName"`
	CategoryId int64  `json:"categoryId" form:"categoryId"`
	SupplierId int64  `json:"supplierId" form:"supplierId"`
	BrandId    int64  `json:"brandId" form:"brandId"`
	TagId      int64  `json:"tagId" form:"tagId"`
	Status     string `json:"status" form:"status"`
}

// GoodsStatusReq 商品状态更新请求
type GoodsStatusReq struct {
	Id     int64  `json:"id" binding:"required"`
	Status string `json:"status" binding:"required"`
}
