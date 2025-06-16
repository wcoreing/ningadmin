package common

// CategoryOption 分类选项
type CategoryOption struct {
	ID       int64            `json:"id"`        // 分类ID
	Name     string           `json:"name"`      // 分类名称
	ParentId int64            `json:"parent_id"` // 父级ID
	Level    int              `json:"level"`     // 层级
	Sort     int              `json:"sort"`      // 排序
	Status   int              `json:"status"`    // 状态
	Icon     string           `json:"icon"`      // 图标
	Image    string           `json:"image"`     // 图片
	Children []CategoryOption `json:"children"`  // 子分类
}

// CategoryListReq 分类列表请求
type CategoryListReq struct {
	Module   string `json:"module" form:"module"`     // 模块标识
	ParentId int64  `json:"parentId" form:"parentId"` // 父级ID
	Status   string `json:"status" form:"status"`     // 状态
	Name     string `json:"name" form:"name"`         // 分类名称
}

// CategoryCreateReq 创建分类请求参数
type CategoryCreateReq struct {
	ParentId int64  `json:"parent_id"`                 // 父级ID
	Name     string `json:"name" binding:"required"`   // 分类名称
	Icon     string `json:"icon"`                      // 分类图标
	Image    string `json:"image"`                     // 分类图片
	Sort     int    `json:"sort"`                      // 排序
	Status   int    `json:"status"`                    // 状态
	Module   string `json:"module" binding:"required"` // 模块标识
}

// CategoryUpdateReq 更新分类请求参数
type CategoryUpdateReq struct {
	Id       int64  `json:"id" binding:"required"`     // 分类ID
	ParentId int64  `json:"parent_id"`                 // 父级ID
	Name     string `json:"name" binding:"required"`   // 分类名称
	Icon     string `json:"icon"`                      // 分类图标
	Image    string `json:"image"`                     // 分类图片
	Sort     int    `json:"sort"`                      // 排序
	Status   int    `json:"status"`                    // 状态
	Module   string `json:"module" binding:"required"` // 模块标识
}
