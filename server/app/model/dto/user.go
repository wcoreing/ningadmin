package dto

import "time"

// UserQuery 用户查询参数
type UserQuery struct {
	PageNum  int    `form:"pageNum"`  // 页码
	PageSize int    `form:"pageSize"` // 每页数量
	Keyword  string `form:"keyword"`  // 关键字(用户名、昵称、手机号、邮箱)
	Status   *int   `form:"status"`   // 状态(0:禁用 1:正常)
	Level    *int   `form:"level"`    // 会员等级
	Gender   *int   `form:"gender"`   // 性别
}

// CreateUserReq 创建用户请求
type CreateUserReq struct {
	Username string `json:"username" binding:"required" comment:"用户名"`
	Password string `json:"password" binding:"required,min=6" comment:"密码"`
	Nickname string `json:"nickname" comment:"昵称"`
	Gender   *int   `json:"gender" comment:"性别"`
	Avatar   string `json:"avatar" comment:"头像"`
	Phone    string `json:"phone" comment:"手机号"`
	Email    string `json:"email" comment:"邮箱"`
	Birthday string `json:"birthday" comment:"生日"`
	Level    int    `json:"level" comment:"会员等级"`
	Status   int    `json:"status" comment:"状态"`
}

// UpdateUserReq 更新用户请求
type UpdateUserReq struct {
	ID       int64  `json:"id" binding:"required" comment:"用户ID"`
	Username string `json:"username" binding:"required" comment:"用户名"`
	Nickname string `json:"nickname" comment:"昵称"`
	Gender   *int   `json:"gender" comment:"性别"`
	Avatar   string `json:"avatar" comment:"头像"`
	Phone    string `json:"phone" comment:"手机号"`
	Email    string `json:"email" comment:"邮箱"`
	Birthday string `json:"birthday" comment:"生日"`
	Level    int    `json:"level" comment:"会员等级"`
	Status   int    `json:"status" comment:"状态"`
}

// UpdateUserStatusReq 更新用户状态请求
type UpdateUserStatusReq struct {
	ID     int64 `json:"id" binding:"required" comment:"用户ID"`
	Status *int  `json:"status" binding:"required" comment:"状态"`
}

// UpdateUserPasswordReq 更新用户密码请求
type UpdateUserPasswordReq struct {
	ID       int64  `json:"id" binding:"required" comment:"用户ID"`
	Password string `json:"password" binding:"required,min=6" comment:"新密码"`
}

// UpdateUserBalanceReq 更新用户余额请求
type UpdateUserBalanceReq struct {
	ID      int64    `json:"id" binding:"required" comment:"用户ID"`
	Balance *float64 `json:"balance" binding:"required" comment:"余额"`
	Remark  string   `json:"remark" comment:"备注"`
}

// UpdateUserPointsReq 更新用户积分请求
type UpdateUserPointsReq struct {
	ID     int64  `json:"id" binding:"required" comment:"用户ID"`
	Points *int   `json:"points" binding:"required" comment:"积分"`
	Remark string `json:"remark" comment:"备注"`
}

// UserVO 用户信息视图对象
type UserVO struct {
	ID            int64      `json:"id"`
	Openid        string     `json:"openid"`
	Username      string     `json:"username"`
	Nickname      string     `json:"nickname"`
	Gender        *int       `json:"gender"`
	GenderText    string     `json:"genderText"`
	Avatar        string     `json:"avatar"`
	Phone         string     `json:"phone"`
	Email         string     `json:"email"`
	Birthday      string     `json:"birthday"`
	Balance       float64    `json:"balance"`
	Points        int        `json:"points"`
	Level         int        `json:"level"`
	LevelText     string     `json:"levelText"`
	TotalAmount   float64    `json:"totalAmount"`
	LastLoginTime *time.Time `json:"lastLoginTime"`
	LastLoginIP   string     `json:"lastLoginIp"`
	Status        int        `json:"status"`
	StatusText    string     `json:"statusText"`
	CreateTime    time.Time  `json:"createTime"`
	UpdateTime    time.Time  `json:"updateTime"`
}

// UserListResp 用户列表响应
type UserListResp struct {
	List  []UserVO `json:"list"`
	Total int64    `json:"total"`
}

// GetGenderText 获取性别文本
func GetGenderText(gender *int) string {
	if gender == nil {
		return "未知"
	}
	switch *gender {
	case 0:
		return "女"
	case 1:
		return "男"
	default:
		return "未知"
	}
}

// GetLevelText 获取会员等级文本
func GetLevelText(level int) string {
	switch level {
	case 1:
		return "普通会员"
	case 2:
		return "黄金会员"
	case 3:
		return "钻石会员"
	default:
		return "普通会员"
	}
}

// GetStatusText 获取状态文本
func GetStatusText(status int) string {
	switch status {
	case 0:
		return "禁用"
	case 1:
		return "正常"
	default:
		return "未知"
	}
}
