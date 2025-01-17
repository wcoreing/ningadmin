package dto

// LoginReq 登录请求参数
type LoginReq struct {
	UserName string `json:"userName" binding:"required" example:"admin"`  // 用户名
	Password string `json:"password" binding:"required" example:"123456"` // 密码
}

// LoginResp 登录响应数据
type LoginResp struct {
	Token    string   `json:"token"`    // 访问令牌
	UserInfo UserInfo `json:"userInfo"` // 用户信息
}

// UserInfo 用户信息
type UserInfo struct {
	UserID      int64  `json:"userId"`      // 用户ID
	UserName    string `json:"userName"`    // 用户名
	NickName    string `json:"nickName"`    // 昵称
	Email       string `json:"email"`       // 邮箱
	Phonenumber string `json:"phonenumber"` // 手机号
	Sex         string `json:"sex"`         // 性别
	Avatar      string `json:"avatar"`      // 头像
	Status      string `json:"status"`      // 状态
	CreateTime  string `json:"createTime"`  // 创建时间
}
