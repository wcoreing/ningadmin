package dto

// DistributionApplyListReq 分销员申请列表请求
type DistributionApplyListReq struct {
	PageNum  int    `form:"pageNum"`
	PageSize int    `form:"pageSize"`
	RealName string `form:"realName"`
	Phone    string `form:"phone"`
	Status   *int   `form:"status"`
}

// DistributionApplyRejectReq 拒绝申请请求
type DistributionApplyRejectReq struct {
	ID     uint   `json:"id" binding:"required"`
	Reason string `json:"reason" binding:"required"`
}

// DistributionApplyVO 分销员申请列表响应
type DistributionApplyVO struct {
	ID           uint   `json:"id"`
	UserID       uint   `json:"userId"`
	RealName     string `json:"realName"`
	Phone        string `json:"phone"`
	IDCard       string `json:"idCard"`
	Status       int    `json:"status"`
	ApplyTime    string `json:"applyTime"`
	AuditTime    string `json:"auditTime"`
	AuditUser    string `json:"auditUser"`
	RejectReason string `json:"rejectReason"`
}
