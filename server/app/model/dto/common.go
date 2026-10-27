package dto

// PageResult 分页结果
type PageResult struct {
	List     interface{} `json:"list"`
	Total    int64       `json:"total"`
	PageNum  int         `json:"pageNum"`
	PageSize int         `json:"pageSize"`
}
