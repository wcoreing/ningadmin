package form

import (
	"encoding/json"
	"time"
)

// FormData 表单数据实体
type FormData struct {
	ID         int64           `gorm:"column:id;primaryKey"` // 数据ID
	FormID     int64           `gorm:"column:form_id"`       // 表单ID
	Data       json.RawMessage `gorm:"column:data"`          // 表单数据
	Submitter  json.RawMessage `gorm:"column:submitter"`     // 提交者信息
	IP         string          `gorm:"column:ip"`            // IP地址
	Status     string          `gorm:"column:status"`        // 状态
	CreateBy   string          `gorm:"column:create_by"`     // 创建人
	CreateTime time.Time       `gorm:"column:create_time"`   // 创建时间
}

// TableName 表名
func (FormData) TableName() string {
	return "form_data"
}
