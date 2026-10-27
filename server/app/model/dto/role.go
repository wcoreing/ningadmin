package dto

// RoleQuery 角色查询参数
type RoleQuery struct {
	PageNum  int    `form:"pageNum" binding:"required,min=1"`          // 页码
	PageSize int    `form:"pageSize" binding:"required,min=1,max=100"` // 每页数量
	RoleName string `form:"roleName"`                                  // 角色名称
	RoleKey  string `form:"roleKey"`                                   // 角色标识
	Status   string `form:"status" binding:"omitempty,oneof=0 1"`      // 状态(0-正常,1-停用)
}

// CreateRoleReq 创建角色请求
type CreateRoleReq struct {
	RoleName string `json:"roleName" binding:"required,min=2,max=30"` // 角色名称
	RoleKey  string `json:"roleKey" binding:"required,min=2,max=100"` // 角色标识
	Status   string `json:"status" binding:"omitempty,oneof=0 1"`     // 状态(0-正常,1-停用)
}

// UpdateRoleReq 更新角色请求
type UpdateRoleReq struct {
	RoleID   int64  `json:"roleId" binding:"required,min=1"`          // 角色ID
	RoleName string `json:"roleName" binding:"required,min=2,max=30"` // 角色名称
	RoleKey  string `json:"roleKey" binding:"required,min=2,max=100"` // 角色标识
	Status   string `json:"status" binding:"omitempty,oneof=0 1"`     // 状态(0-正常,1-停用)
}

// UpdateRoleStatusReq 更新角色状态请求
type UpdateRoleStatusReq struct {
	RoleID int64  `json:"roleId" binding:"required,min=1"`     // 角色ID
	Status string `json:"status" binding:"required,oneof=0 1"` // 状态(0-正常,1-停用)
}

// RolePermissionReq 角色权限请求
type RolePermissionReq struct {
	RoleID        int64   `json:"roleId" binding:"required,min=1"`   // 角色ID
	PermissionIds []int64 `json:"permissionIds" binding:"omitempty"` // 权限ID列表
}

// RoleInfoResp 角色信息响应
type RoleInfoResp struct {
	RoleID     int64  `json:"roleId"`     // 角色ID
	RoleName   string `json:"roleName"`   // 角色名称
	RoleKey    string `json:"roleKey"`    // 角色标识
	Status     string `json:"status"`     // 状态
	CreateTime string `json:"createTime"` // 创建时间
	CreateBy   string `json:"createBy"`   // 创建者
	UpdateTime string `json:"updateTime"` // 更新时间
	UpdateBy   string `json:"updateBy"`   // 更新者
}

// RoleListResp 角色列表响应
type RoleListResp struct {
	Total int64          `json:"total"` // 总数
	Rows  []RoleInfoResp `json:"rows"`  // 数据列表
}

// RolePermissionResp 角色权限响应
type RolePermissionResp struct {
	RoleID        int64   `json:"roleId"`        // 角色ID
	PermissionIds []int64 `json:"permissionIds"` // 权限ID列表
}
