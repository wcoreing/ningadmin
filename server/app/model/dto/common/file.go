package common

import "wnadmin/app/model/dto/base"

// FileQuery 文件查询参数
type FileQuery struct {
	base.PageQuery
	CategoryId int64  `form:"category_id" json:"category_id"`
	Keyword    string `form:"keyword" json:"keyword"`
}

// FileVO 文件信息VO
type FileVO struct {
	Id           int64  `json:"id"`
	CategoryId   int64  `json:"categoryId"`
	Name         string `json:"name"`
	OriginalName string `json:"originalName"`
	Url          string `json:"url"`
	Size         int64  `json:"size"`
	Type         string `json:"type"`
	CreateTime   string `json:"createTime"`
}
