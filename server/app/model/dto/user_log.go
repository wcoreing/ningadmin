package dto

import "time"

// UserChangeLogQuery 用户变动记录查询参数
type UserChangeLogQuery struct {
	PageNum    int    `form:"pageNum"`    // 页码
	PageSize   int    `form:"pageSize"`   // 每页数量
	UserID     *int64 `form:"userId"`     // 用户ID
	ChangeType *int   `form:"changeType"` // 变动数据类型(1-余额 2-积分)
	ActionType *int   `form:"actionType"` // 变动行为(1-收入/增加 2-支出/减少)
	SourceType *int   `form:"sourceType"` // 来源类型
	StartTime  string `form:"startTime"`  // 开始时间
	EndTime    string `form:"endTime"`    // 结束时间
}

// UserChangeLogVO 用户变动记录视图对象
type UserChangeLogVO struct {
	ID             int64      `json:"id"`
	UserID         int64      `json:"userId"`
	Username       string     `json:"username"`       // 用户名
	Nickname       string     `json:"nickname"`       // 昵称
	ChangeType     int        `json:"changeType"`     // 变动数据类型
	ChangeTypeText string     `json:"changeTypeText"` // 变动数据类型文本
	ActionType     int        `json:"actionType"`     // 变动行为
	ActionTypeText string     `json:"actionTypeText"` // 变动行为文本
	Amount         float64    `json:"amount"`         // 变动数量
	BeforeValue    float64    `json:"beforeValue"`    // 变动前数值
	AfterValue     float64    `json:"afterValue"`     // 变动后数值
	SourceType     int        `json:"sourceType"`     // 来源类型
	SourceTypeText string     `json:"sourceTypeText"` // 来源类型文本
	SourceID       *int64     `json:"sourceId"`       // 来源ID
	Remark         string     `json:"remark"`         // 备注
	OperatorID     *int64     `json:"operatorId"`     // 操作人ID
	OperatorName   string     `json:"operatorName"`   // 操作人姓名
	ExpireTime     *time.Time `json:"expireTime"`     // 过期时间
	CreateTime     time.Time  `json:"createTime"`     // 创建时间
}

// UserChangeLogListResp 用户变动记录列表响应
type UserChangeLogListResp struct {
	List  []UserChangeLogVO `json:"list"`
	Total int64             `json:"total"`
}
