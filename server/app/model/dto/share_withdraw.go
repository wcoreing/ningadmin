package dto

// ShareWithdrawListReq 提现列表请求
type ShareWithdrawListReq struct {
	PageNum  int    `form:"pageNum"`
	PageSize int    `form:"pageSize"`
	RealName string `form:"realName"`
	Phone    string `form:"phone"`
	Status   *int   `form:"status"`
}

// ShareWithdrawRejectReq 拒绝提现请求
type ShareWithdrawRejectReq struct {
	ID     uint   `json:"id" binding:"required"`
	Remark string `json:"remark" binding:"required"`
}

// ShareWithdrawVO 提现列表响应
type ShareWithdrawVO struct {
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

// ShareWithdrawPayReq 确认打款请求
type ShareWithdrawPayReq struct {
	ID     uint   `json:"id" binding:"required"`
	Remark string `json:"remark" binding:"required"`
}
