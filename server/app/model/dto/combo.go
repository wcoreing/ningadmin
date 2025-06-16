package dto

// ComboListReq 套餐列表请求参数
type ComboListReq struct {
	PageNum    int    `json:"pageNum" form:"pageNum"`       // 页码
	PageSize   int    `json:"pageSize" form:"pageSize"`     // 每页数量
	Name       string `json:"name" form:"name"`             // 套餐名称
	CategoryId int64  `json:"categoryId" form:"categoryId"` // 分类ID
	Status     string `json:"status" form:"status"`         // 状态
}

// ComboDTO 套餐信息DTO
type ComboDTO struct {
	Id            int64   `json:"id"`
	Name          string  `json:"name"`
	Description   string  `json:"description"`
	Image         string  `json:"image"`
	CategoryId    int64   `json:"categoryId"`
	CategoryName  string  `json:"categoryName"`
	BasePrice     float64 `json:"basePrice"`
	OriginalPrice float64 `json:"originalPrice"`
	GroupCount    int     `json:"groupCount"` // 套餐组数量
	Sort          int     `json:"sort"`
	Status        int8    `json:"status"`
	CreateTime    string  `json:"createTime"`
}

// ComboStatusReq 套餐状态修改请求
type ComboStatusReq struct {
	Id     int64 `json:"id" binding:"required"`
	Status int8  `json:"status" binding:"gte=0,lte=1"`
}

// ComboImageItem 套餐图片项
type ComboImageItem struct {
	Url      string `json:"url"`      // 图片URL
	Position int8   `json:"position"` // 图片位置：0-详情图，1-封面图
	Sort     int    `json:"sort"`     // 排序
}

// ComboGroupGoods 套餐组商品
type ComboGroupGoods struct {
	Id         int64   `json:"id,omitempty"`
	GoodsId    int64   `json:"goodsId"`
	GoodsName  string  `json:"goodsName"`
	GoodsImage string  `json:"goodsImage"`
	SkuId      *int64  `json:"skuId"`
	SkuInfo    *string `json:"skuInfo"`
	AddPrice   float64 `json:"addPrice"`
	Quantity   int     `json:"quantity"`
	IsDefault  bool    `json:"isDefault"`
	Sort       int     `json:"sort"`
}

// ComboGroup 套餐组
type ComboGroup struct {
	Id          int64              `json:"id,omitempty"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	MinSelect   int                `json:"minSelect"`
	MaxSelect   int                `json:"maxSelect"`
	IsRequired  bool               `json:"isRequired"`
	Sort        int                `json:"sort"`
	Goods       []*ComboGroupGoods `json:"goods"`
}

// ComboCreateReq 套餐创建请求
type ComboCreateReq struct {
	Name          string            `json:"name" binding:"required"`
	Description   string            `json:"description"`
	Image         string            `json:"image" binding:"required"`
	CategoryId    int64             `json:"categoryId" binding:"required"`
	BasePrice     float64           `json:"basePrice" binding:"gte=0"`
	OriginalPrice float64           `json:"originalPrice" binding:"gte=0"`
	Sort          int               `json:"sort"`
	Status        int8              `json:"status" binding:"oneof=0 1"`
	Images        []*ComboImageItem `json:"images"`
	Groups        []*ComboGroup     `json:"groups"`
	DetailDesc    string            `json:"detailDesc"`
}

// ComboUpdateReq 套餐更新请求
type ComboUpdateReq struct {
	Id            int64             `json:"id" binding:"required"`
	Name          string            `json:"name" binding:"required"`
	Description   string            `json:"description"`
	Image         string            `json:"image" binding:"required"`
	CategoryId    int64             `json:"categoryId" binding:"required"`
	BasePrice     float64           `json:"basePrice" binding:"gte=0"`
	OriginalPrice float64           `json:"originalPrice" binding:"gte=0"`
	Sort          int               `json:"sort"`
	Status        int8              `json:"status" binding:"oneof=0 1"`
	Images        []*ComboImageItem `json:"images"`
	Groups        []*ComboGroup     `json:"groups"`
	DetailDesc    string            `json:"detailDesc"`
}

// ComboInfoResp 套餐详情响应
type ComboInfoResp struct {
	Id            int64             `json:"id"`
	Name          string            `json:"name"`
	Description   string            `json:"description"`
	Image         string            `json:"image"`
	CategoryId    int64             `json:"categoryId"`
	BasePrice     float64           `json:"basePrice"`
	OriginalPrice float64           `json:"originalPrice"`
	Sort          int               `json:"sort"`
	Status        int8              `json:"status"`
	Images        []*ComboImageItem `json:"images"`
	Groups        []*ComboGroup     `json:"groups"`
	DetailDesc    string            `json:"detailDesc"`
}

// ComboBatchDeleteReq 套餐批量删除请求
type ComboBatchDeleteReq struct {
	Ids []int64 `json:"ids" binding:"required"`
}
