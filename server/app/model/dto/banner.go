package dto

// BannerInfoResp Banner详情响应
type BannerInfoResp struct {
	Id        int64          `json:"id"`
	Title     string         `json:"title"`
	Image     string         `json:"image"`
	Url       string         `json:"url"`
	Position  int            `json:"position"`
	Type      int            `json:"type"`
	TargetId  int64          `json:"targetId"`
	Sort      int            `json:"sort"`
	Status    int            `json:"status"`
	GoodsInfo *GoodsBaseInfo `json:"goodsInfo,omitempty"`
}

type BannerListReq struct {
	Position int    `form:"position"`
	Status   *int   `form:"status"`
	Title    string `form:"title"`
	PageNum  int    `form:"pageNum,default=1"`
	PageSize int    `form:"pageSize,default=10"`
}

// BannerListResp Banner列表响应
type BannerListResp struct {
	Id        int64          `json:"id"`
	Title     string         `json:"title"`
	Image     string         `json:"image"`
	Position  int            `json:"position"`
	Type      int            `json:"type"`
	TargetId  int64          `json:"targetId"`
	Url       string         `json:"url"`
	Sort      int            `json:"sort"`
	Status    int            `json:"status"`
	GoodsInfo *GoodsBaseInfo `json:"goodsInfo,omitempty"`
}

// GoodsBaseInfo 商品基本信息
type GoodsBaseInfo struct {
	GoodsId  uint    `json:"goodsId"`
	Name     string  `json:"name"`
	CoverImg string  `json:"coverImg"`
	Price    float64 `json:"price"`
}

type BannerCreateReq struct {
	Title    string `json:"title" binding:"required"`
	Image    string `json:"image" binding:"required"`
	Url      string `json:"url"`
	Position int    `json:"position" binding:"required"`
	Type     int    `json:"type" binding:"required"`
	TargetId int64  `json:"targetId"`
	Sort     int    `json:"sort"`
	Status   int    `json:"status"`
}

type BannerUpdateReq struct {
	Id       int64  `json:"id" binding:"required"`
	Title    string `json:"title" binding:"required"`
	Image    string `json:"image" binding:"required"`
	Url      string `json:"url"`
	Position int    `json:"position" binding:"required"`
	Type     int    `json:"type" binding:"required"`
	TargetId int64  `json:"targetId"`
	Sort     int    `json:"sort"`
	Status   int    `json:"status"`
}
