package dto

// GoodsBatchMoveReq 商品批量移动请求
type GoodsBatchMoveReq struct {
	GoodsIds   []int64 `json:"goodsIds" binding:"required"`
	CategoryId int64   `json:"categoryId" binding:"required"`
}

// GoodsBatchDeleteReq 商品批量删除请求
type GoodsBatchDeleteReq struct {
	GoodsIds []int64 `json:"goodsIds" binding:"required"`
}
