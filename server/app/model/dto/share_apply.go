package dto

// ShareApplyListReq 分享申请列表请求
type ShareApplyListReq struct {
	PageNum  int    `form:"pageNum"`
	PageSize int    `form:"pageSize"`
	RealName string `form:"realName"`
	Phone    string `form:"phone"`
	Status   *int   `form:"status"`
}

// ShareApplyRejectReq 拒绝申请请求
type ShareApplyRejectReq struct {
	ID     uint   `json:"id" binding:"required"`
	Reason string `json:"reason" binding:"required"`
}

// ShareApplyVO 分享申请列表响应
type ShareApplyVO struct {
	ID           uint   `json:"id"`
	UserID       uint   `json:"userId"`
	Username     string `json:"username"` // 用户名(从user表读取)
	Nickname     string `json:"nickname"` // 昵称(从user表读取)
	Avatar       string `json:"avatar"`   // 头像(从user表读取)
	RealName     string `json:"realName"`
	Phone        string `json:"phone"`
	IDCard       string `json:"idCard"`
	Status       int    `json:"status"`
	ApplyTime    string `json:"applyTime"`
	AuditTime    string `json:"auditTime"`
	AuditUser    string `json:"auditUser"`
	RejectReason string `json:"rejectReason"`
}
