package dto

// GoodsSpecGroup 商品规格组
type GoodsSpecGroup struct {
	ID     int64             `json:"id"`     // 规格组ID
	Name   string            `json:"name"`   // 规格组名称
	Sort   int               `json:"sort"`   // 排序
	Values []*GoodsSpecValue `json:"values"` // 规格值列表
}

// GoodsSpecValue 商品规格值
type GoodsSpecValue struct {
	ID    int64  `json:"id"`    // 规格值ID
	Value string `json:"value"` // 规格值
	Sort  int    `json:"sort"`  // 排序
}
