package dto

// MenuDTO 菜单数据传输对象
type MenuDTO struct {
	ID            int64  `json:"id"`
	ParentID      int64  `json:"parentId"`
	Type          string `json:"type" binding:"required"`
	Name          string `json:"name" binding:"required"`
	Icon          string `json:"icon"`
	RoutePath     string `json:"routePath"`
	ComponentPath string `json:"componentPath"`
	Permission    string `json:"permission"`
	Sort          int    `json:"sort"`
	Status        string `json:"status" binding:"required"`
	Operator      string `json:"-"` // 操作人，从 context 中获取
}

// MenuListReq 菜单列表请求参数
type MenuListReq struct {
	Keyword string `form:"keyword"`
	Status  string `form:"status"`
}

// MenuRoleDTO 菜单角色关联
type MenuRoleDTO struct {
	ID      int64   `json:"id" binding:"required"`
	RoleIDs []int64 `json:"roleIds"`
}

// MenuStatusDTO 菜单状态
type MenuStatusDTO struct {
	ID     int64  `json:"id" binding:"required"`
	Status string `json:"status" binding:"required"`
}
