package dto

// PermissionQuery 权限查询参数
type PermissionQuery struct {
	Keyword string `form:"keyword"`                              // 搜索关键字
	Type    string `form:"type" binding:"omitempty"`             // 权限类型
	Status  string `form:"status" binding:"omitempty,oneof=0 1"` // 状态(0-正常,1-停用)
}

// CreatePermissionReq 创建权限请求
type CreatePermissionReq struct {
	ParentID      int64  `json:"parentId"`                              // 上级权限ID
	Type          string `json:"type" binding:"required"`               // 权限类型
	Name          string `json:"name" binding:"required,min=2,max=30"`  // 权限名称
	Icon          string `json:"icon"`                                  // 图标
	RoutePath     string `json:"routePath"`                             // 路由路径
	ComponentPath string `json:"componentPath"`                         // 组件路径
	Permission    string `json:"permission" binding:"required,max=100"` // 权限标识
	Sort          int    `json:"sort" binding:"gte=0"`                  // 排序
	Status        string `json:"status" binding:"omitempty,oneof=0 1"`  // 状态(0-正常,1-停用)
}

// UpdatePermissionReq 更新权限请求
type UpdatePermissionReq struct {
	ID            int64  `json:"id" binding:"required"`                 // 权限ID
	ParentID      int64  `json:"parentId"`                              // 上级权限ID
	Type          string `json:"type" binding:"required"`               // 权限类型
	Name          string `json:"name" binding:"required,min=2,max=30"`  // 权限名称
	Icon          string `json:"icon"`                                  // 图标
	RoutePath     string `json:"routePath"`                             // 路由路径
	ComponentPath string `json:"componentPath"`                         // 组件路径
	Permission    string `json:"permission" binding:"required,max=100"` // 权限标识
	Sort          int    `json:"sort" binding:"gte=0"`                  // 排序
	Status        string `json:"status" binding:"omitempty,oneof=0 1"`  // 状态(0-正常,1-停用)
}

// UpdatePermissionStatusReq 更新权限状态请求
type UpdatePermissionStatusReq struct {
	ID     int64  `json:"id" binding:"required"`               // 权限ID
	Status string `json:"status" binding:"required,oneof=0 1"` // 状态(0-正常,1-停用)
}

// PermissionInfoResp 权限信息响应
type PermissionInfoResp struct {
	ID            int64                `json:"id"`            // 权限ID
	ParentID      int64                `json:"parentId"`      // 上级权限ID
	Type          string               `json:"type"`          // 权限类型
	Name          string               `json:"name"`          // 权限名称
	Icon          string               `json:"icon"`          // 图标
	RoutePath     string               `json:"routePath"`     // 路由路径
	ComponentPath string               `json:"componentPath"` // 组件路径
	Permission    string               `json:"permission"`    // 权限标识
	Sort          int                  `json:"sort"`          // 排序
	Status        string               `json:"status"`        // 状态
	CreateTime    string               `json:"createTime"`    // 创建时间
	Children      []PermissionInfoResp `json:"children"`      // 子权限
}

// PermissionTreeResp 权限树响应
type PermissionTreeResp struct {
	Total int64                `json:"total"` // 总数
	Rows  []PermissionInfoResp `json:"rows"`  // 数据列表
}

// PermissionTreeNode 权限树节点
type PermissionTreeNode struct {
	ID       int64                `json:"id"`       // 权限ID
	ParentID int64                `json:"parentId"` // 上级权限ID
	Name     string               `json:"name"`     // 权限名称
	Children []PermissionTreeNode `json:"children"` // 子权限
}
