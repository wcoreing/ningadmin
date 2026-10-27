package base

// PageQuery 分页查询参数
type PageQuery struct {
	Page     int `form:"page" json:"page"`          // 页码
	PageSize int `form:"page_size" json:"pageSize"` // 每页数量
}
