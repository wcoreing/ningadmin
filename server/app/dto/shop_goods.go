package dto

// 商品列表请求
type GoodsListReq struct {
	PageNum   int    `json:"pageNum"`
	PageSize  int    `json:"pageSize"`
	GoodsName string `json:"goodsName"`
	Status    string `json:"status"`
}

// 商品列表响应
type GoodsListResp struct {
	List  []GoodsInfo `json:"list"`
	Total int64       `json:"total"`
}

// 商品信息
type GoodsInfo struct {
	GoodsId     int64   `json:"goodsId"`
	CategoryId  int64   `json:"categoryId"`
	GoodsName   string  `json:"goodsName"`
	CoverImg    string  `json:"coverImg"`
	Brief       string  `json:"brief"`
	Price       float64 `json:"price"`
	MarketPrice float64 `json:"marketPrice"`
	Stock       int     `json:"stock"`
	Sales       int     `json:"sales"`
	Status      string  `json:"status"`
}

// 商品状态更新请求
type GoodsStatusReq struct {
	Id     int64  `json:"id"`
	Status string `json:"status"`
}
