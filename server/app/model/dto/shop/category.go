package shop

// CategoryListReq 分类列表请求参数
type CategoryListReq struct {
	Name     string `form:"name"`      // 分类名称
	Status   *int   `form:"status"`    // 状态(0:禁用 1:启用)
	ParentId *int64 `form:"parent_id"` // 父级ID
}

// CategoryCreateReq 创建分类请求参数
type CategoryCreateReq struct {
	ParentId int64  `json:"parent_id"`               // 父级ID
	Name     string `json:"name" binding:"required"` // 分类名称
	Icon     string `json:"icon"`                    // 分类图标
	Image    string `json:"image"`                   // 分类图片
	Sort     int    `json:"sort"`                    // 排序
	Status   int    `json:"status"`                  // 状态
}

// CategoryUpdateReq 更新分类请求参数
type CategoryUpdateReq struct {
	Id       int64  `json:"id" binding:"required"`   // 分类ID
	ParentId int64  `json:"parent_id"`               // 父级ID
	Name     string `json:"name" binding:"required"` // 分类名称
	Icon     string `json:"icon"`                    // 分类图标
	Image    string `json:"image"`                   // 分类图片
	Sort     int    `json:"sort"`                    // 排序
	Status   int    `json:"status"`                  // 状态
}

// CategoryVO 分类信息
type CategoryVO struct {
	Id         int64        `json:"id"`          // 分类ID
	ParentId   int64        `json:"parent_id"`   // 父级ID
	Name       string       `json:"name"`        // 分类名称
	Icon       string       `json:"icon"`        // 分类图标
	Image      string       `json:"image"`       // 分类图片
	Level      int          `json:"level"`       // 层级
	Sort       int          `json:"sort"`        // 排序
	Status     int          `json:"status"`      // 状态
	Children   []CategoryVO `json:"children"`    // 子分类
	CreateTime string       `json:"create_time"` // 创建时间
	UpdateTime string       `json:"update_time"` // 更新时间
}
