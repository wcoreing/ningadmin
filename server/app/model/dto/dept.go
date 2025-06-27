package dto

// DeptQuery 部门查询参数
type DeptQuery struct {
	PageNum  int    `form:"pageNum"`
	PageSize int    `form:"pageSize"`
	Keyword  string `form:"keyword"`
	Status   string `form:"status"`
	ParentId *int64 `form:"parentId"`
}

// CreateDeptReq 创建部门请求
type CreateDeptReq struct {
	ParentId int64  `json:"parentId"`
	DeptName string `json:"deptName" binding:"required"`
	OrderNum int    `json:"orderNum"`
	Leader   string `json:"leader"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Status   string `json:"status"`
}

// UpdateDeptReq 更新部门请求
type UpdateDeptReq struct {
	DeptId   int64  `json:"deptId" binding:"required"`
	ParentId int64  `json:"parentId"`
	DeptName string `json:"deptName" binding:"required"`
	OrderNum int    `json:"orderNum"`
	Leader   string `json:"leader"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Status   string `json:"status"`
}

// UpdateDeptStatusReq 更新部门状态请求
type UpdateDeptStatusReq struct {
	DeptId int64  `json:"deptId" binding:"required"`
	Status string `json:"status" binding:"required"`
}

// DeptOption 部门选项
type DeptOption struct {
	Value int64  `json:"value"` // 部门ID
	Label string `json:"label"` // 部门名称
}
