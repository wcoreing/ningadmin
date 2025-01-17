package dto

// UserQuery 用户查询参数
type UserQuery struct {
	PageNum  int    `form:"pageNum"`
	PageSize int    `form:"pageSize"`
	Keyword  string `form:"keyword"`
	Status   string `form:"status"`
}

// CreateUserReq 创建用户请求
type CreateUserReq struct {
	UserName    string `json:"userName" binding:"required"`
	Password    string `json:"password" binding:"required"`
	NickName    string `json:"nickName"`
	Email       string `json:"email"`
	Phonenumber string `json:"phonenumber"`
	Sex         string `json:"sex"`
	Status      string `json:"status"`
}

// UpdateUserReq 更新用户请求
type UpdateUserReq struct {
	ID          int64  `json:"userId" binding:"required"`
	UserName    string `json:"userName" binding:"required"`
	NickName    string `json:"nickName"`
	Email       string `json:"email"`
	Phonenumber string `json:"phonenumber"`
	Sex         string `json:"sex"`
	Status      string `json:"status"`
}

// UpdateUserStatusReq 更新用户状态请求
type UpdateUserStatusReq struct {
	UserID int64  `json:"userId" binding:"required"`
	Status string `json:"status" binding:"required"`
}

// UpdateUserPasswordReq 更新用户密码请求
type UpdateUserPasswordReq struct {
	UserID   int64  `json:"userId" binding:"required"`
	Password string `json:"password" binding:"required"`
}
