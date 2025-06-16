package dto

// SystemUserQuery 系统用户查询参数
type SystemUserQuery struct {
	PageNum  int    `form:"pageNum"`                              // 页码
	PageSize int    `form:"pageSize"`                             // 每页数量
	Keyword  string `form:"keyword"`                              // 搜索关键字
	Status   string `form:"status" binding:"omitempty,oneof=0 1"` // 状态(0-正常,1-停用)
}

// CreateSystemUserReq 创建系统用户请求
type CreateSystemUserReq struct {
	UserName    string `json:"userName" binding:"required"` // 用户名
	Password    string `json:"password" binding:"required"` // 密码
	NickName    string `json:"nickName"`                    // 昵称
	Email       string `json:"email"`                       // 邮箱
	Phonenumber string `json:"phonenumber"`                 // 手机号
	Sex         string `json:"sex"`                         // 性别
	Status      string `json:"status"`                      // 状态
	RoleID      int64  `json:"roleId" binding:"required"`   // 角色ID
}

// UpdateSystemUserReq 更新系统用户请求
type UpdateSystemUserReq struct {
	ID          int64  `json:"userId" binding:"required"`   // 用户ID
	UserName    string `json:"userName" binding:"required"` // 用户名
	NickName    string `json:"nickName"`                    // 昵称
	Email       string `json:"email"`                       // 邮箱
	Phonenumber string `json:"phonenumber"`                 // 手机号
	Sex         string `json:"sex"`                         // 性别
	Status      string `json:"status"`                      // 状态
	RoleID      int64  `json:"roleId" binding:"required"`   // 角色ID
}

// UpdateSystemUserStatusReq 更新系统用户状态请求
type UpdateSystemUserStatusReq struct {
	UserID int64  `json:"userId" binding:"required"` // 用户ID
	Status string `json:"status" binding:"required"` // 状态
}

// UpdateSystemUserPasswordReq 更新系统用户密码请求
type UpdateSystemUserPasswordReq struct {
	UserID   int64  `json:"userId" binding:"required"`   // 用户ID
	Password string `json:"password" binding:"required"` // 密码
}

// SystemUserInfoResp 系统用户信息响应
type SystemUserInfoResp struct {
	ID          int64     `json:"userId"`      // 用户ID
	UserName    string    `json:"userName"`    // 用户名
	NickName    string    `json:"nickName"`    // 昵称
	Email       string    `json:"email"`       // 邮箱
	Phonenumber string    `json:"phonenumber"` // 手机号
	Sex         string    `json:"sex"`         // 性别
	Status      string    `json:"status"`      // 状态
	CreateTime  string    `json:"createTime"`  // 创建时间
	CreateBy    string    `json:"createBy"`    // 创建者
	UpdateTime  string    `json:"updateTime"`  // 更新时间
	UpdateBy    string    `json:"updateBy"`    // 更新者
	Role        *RoleInfo `json:"role"`        // 角色信息
}

// SystemUserListResp 系统用户列表响应
type SystemUserListResp struct {
	Total int64                `json:"total"` // 总数
	Rows  []SystemUserInfoResp `json:"rows"`  // 数据列表
}
