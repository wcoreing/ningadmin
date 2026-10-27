package common

// FileCategoryTreeVO 文件分类树
type FileCategoryTreeVO struct {
	Id       int64                `json:"id"`
	ParentId int64                `json:"parentId"`
	Name     string               `json:"name"`
	Code     string               `json:"code"`
	Icon     string               `json:"icon"`
	Sort     int                  `json:"sort"`
	Count    int                  `json:"count"`
	Children []FileCategoryTreeVO `json:"children"`
}
