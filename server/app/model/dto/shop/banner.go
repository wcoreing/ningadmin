package shop

// BannerListReq Banner列表请求参数
type BannerListReq struct {
	PageNum  int    `form:"pageNum" binding:"required"`
	PageSize int    `form:"pageSize" binding:"required"`
	Title    string `form:"title"`
	Position int    `form:"position"`
	Status   *int   `form:"status"`
}

type BannerCreateReq struct {
	Title    string `json:"title" binding:"required"`
	Image    string `json:"image" binding:"required"`
	Url      string `json:"url"`
	Position int    `json:"position" binding:"required"`
	Type     int    `json:"type" binding:"required"`
	TargetId int64  `json:"targetId"`
	Sort     int    `json:"sort"`
	Status   int    `json:"status" binding:"oneof=0 1"`
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
	Status   int    `json:"status" binding:"oneof=0 1"`
}
