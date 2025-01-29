package vo

// GoodsListVO 商品列表响应
type GoodsListVO struct {
	List  []GoodsVO `json:"list"`
	Total int64     `json:"total"`
}

// GoodsVO 商品信息
type GoodsVO struct {
	GoodsId       int64   `json:"goodsId"`
	CategoryId    int64   `json:"categoryId"`
	Type          int8    `json:"type"`
	Name          string  `json:"name"`
	Description   string  `json:"description"`
	Price         float64 `json:"price"`
	OriginalPrice float64 `json:"originalPrice"`
	Stock         int     `json:"stock"`
	Sales         int     `json:"sales"`
	Status        int8    `json:"status"`
	Sort          int     `json:"sort"`
}
