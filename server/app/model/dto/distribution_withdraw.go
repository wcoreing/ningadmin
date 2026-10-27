package dto

// WithdrawListReq 提现列表请求
type WithdrawListReq struct {
	PageNum  int    `form:"pageNum"`
	PageSize int    `form:"pageSize"`
	RealName string `form:"realName"`
	Phone    string `form:"phone"`
	Status   *int   `form:"status"`
}

// WithdrawRejectReq 拒绝提现请求
type WithdrawRejectReq struct {
	ID     uint   `json:"id" binding:"required"`
	Remark string `json:"remark" binding:"required"`
}

// WithdrawVO 提现列表响应
type WithdrawVO struct {
	ID          uint    `json:"id"`
	UserID      uint    `json:"userId"`
	Amount      float64 `json:"amount"`
	Method      string  `json:"method"`
	Account     string  `json:"account"`
	RealName    string  `json:"realName"`
	Status      int     `json:"status"`
	CreateTime  string  `json:"createTime"`
	AuditTime   string  `json:"auditTime"`
	AuditRemark string  `json:"auditRemark"`
	PayTime     string  `json:"payTime"`
	PayRemark   string  `json:"payRemark"`
}

// WithdrawPayReq 确认打款请求
type WithdrawPayReq struct {
	ID     uint   `json:"id" binding:"required"`
	Remark string `json:"remark" binding:"required"`
}
