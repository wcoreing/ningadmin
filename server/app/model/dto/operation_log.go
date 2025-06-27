package dto

// OperationLogData 操作日志数据
type OperationLogData struct {
	TraceID       string      `json:"traceId"`       // 链路追踪ID
	UserID        int64       `json:"userId"`        // 用户ID
	Username      string      `json:"username"`      // 用户名
	UserType      string      `json:"userType"`      // 用户类型
	OperationType string      `json:"operationType"` // 操作类型
	Module        string      `json:"module"`        // 操作模块
	OperationName string      `json:"operationName"` // 操作名称
	Method        string      `json:"method"`        // 请求方法
	RequestURL    string      `json:"requestUrl"`    // 请求URL
	RequestParams interface{} `json:"requestParams"` // 请求参数
	ResponseData  interface{} `json:"responseData"`  // 响应数据
	IPAddress     string      `json:"ipAddress"`     // IP地址
	UserAgent     string      `json:"userAgent"`     // 用户代理
	Location      string      `json:"location"`      // 操作地点
	Status        string      `json:"status"`        // 操作状态
	ErrorMsg      string      `json:"errorMsg"`      // 错误信息
	CostTime      int64       `json:"costTime"`      // 耗时(毫秒)
}

// OperationLogListReq 操作日志列表请求
type OperationLogListReq struct {
	Page          int    `form:"page" json:"page"`          // 页码
	PageSize      int    `form:"page_size" json:"pageSize"` // 每页数量
	Username      string `form:"username"`                  // 用户名
	Module        string `form:"module"`                    // 模块
	OperationType string `form:"operationType"`             // 操作类型
	StartTime     string `form:"startTime"`                 // 开始时间
	EndTime       string `form:"endTime"`                   // 结束时间
}

// OperationInfo 操作信息
type OperationInfo struct {
	Module string `json:"module"` // 模块名称
	Name   string `json:"name"`   // 操作名称
	Type   string `json:"type"`   // 操作类型
}

// OperationLogBatchDeleteReq 批量删除操作日志请求
type OperationLogBatchDeleteReq struct {
	Ids []int64 `json:"ids" binding:"required"` // 要删除的日志ID列表
}

// OperationLogClearReq 清空操作日志请求
type OperationLogClearReq struct {
	Days int `json:"days" binding:"required,min=1,max=365"` // 保留天数
}
