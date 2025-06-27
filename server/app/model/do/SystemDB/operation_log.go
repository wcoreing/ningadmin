package SystemDB

import (
	"wnadmin/app/common/utils"
	"wnadmin/app/provider/db"

	"gorm.io/gorm"
)

// OperationLogDO 操作日志数据对象
type OperationLogDO struct {
	ID            int64          `gorm:"column:id;primaryKey;autoIncrement" json:"id"`                 // 日志ID
	TraceID       string         `gorm:"column:trace_id;size:64" json:"traceId"`                       // 链路追踪ID
	UserID        int64          `gorm:"column:user_id" json:"userId"`                                 // 用户ID
	Username      string         `gorm:"column:username;size:64" json:"username"`                      // 用户名
	UserType      string         `gorm:"column:user_type;size:20;default:admin" json:"userType"`       // 用户类型
	OperationType string         `gorm:"column:operation_type;size:20;not null" json:"operationType"`  // 操作类型
	Module        string         `gorm:"column:module;size:50;not null" json:"module"`                 // 操作模块
	OperationName string         `gorm:"column:operation_name;size:100;not null" json:"operationName"` // 操作名称
	Method        string         `gorm:"column:method;size:10" json:"method"`                          // 请求方法
	RequestURL    string         `gorm:"column:request_url;size:500" json:"requestUrl"`                // 请求URL
	RequestParams string         `gorm:"column:request_params;type:text" json:"requestParams"`         // 请求参数
	ResponseData  string         `gorm:"column:response_data;type:text" json:"responseData"`           // 响应数据
	IPAddress     string         `gorm:"column:ip_address;size:45" json:"ipAddress"`                   // IP地址
	UserAgent     string         `gorm:"column:user_agent;size:500" json:"userAgent"`                  // 用户代理
	Location      string         `gorm:"column:location;size:100" json:"location"`                     // 操作地点
	Status        string         `gorm:"column:status;size:1;default:1" json:"status"`                 // 操作状态
	ErrorMsg      string         `gorm:"column:error_msg;type:text" json:"errorMsg"`                   // 错误信息
	CostTime      int64          `gorm:"column:cost_time" json:"costTime"`                             // 耗时(毫秒)
	CreateTime    utils.JSONTime `gorm:"column:create_time;autoCreateTime" json:"createTime"`          // 创建时间
}

// TableName 表名
func (OperationLogDO) TableName() string {
	return "sys_operation_log"
}

// Query 查询构造器
func (o *OperationLogDO) Query() *gorm.DB {
	return db.WnAdmin().Model(o)
}

// Create 创建记录
func (o *OperationLogDO) Create() error {
	return db.WnAdmin().Create(o).Error
}

// Find 查询多条记录
func (o *OperationLogDO) Find(dest interface{}) error {
	return o.Query().Find(dest).Error
}

// GetLogList 获取日志列表（分页）
func (o *OperationLogDO) GetLogList(page, pageSize int, conditions map[string]interface{}) ([]OperationLogDO, int64, error) {
	var logs []OperationLogDO
	var total int64

	query := o.Query()

	// 添加查询条件
	if username, ok := conditions["username"]; ok && username != "" {
		query = query.Where("username LIKE ?", "%"+username.(string)+"%")
	}
	if module, ok := conditions["module"]; ok && module != "" {
		query = query.Where("module = ?", module)
	}
	if operationType, ok := conditions["operation_type"]; ok && operationType != "" {
		query = query.Where("operation_type = ?", operationType)
	}
	if startTime, ok := conditions["start_time"]; ok && startTime != "" {
		query = query.Where("create_time >= ?", startTime)
	}
	if endTime, ok := conditions["end_time"]; ok && endTime != "" {
		query = query.Where("create_time <= ?", endTime)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err := query.Order("create_time desc").Offset(offset).Limit(pageSize).Find(&logs).Error

	return logs, total, err
}
