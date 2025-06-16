package do

import (
	"time"
	"wnadmin/app/provider/db"

	"gorm.io/gorm"
)

// 常量定义
const (
	// 变动数据类型
	ChangeTypeBalance = 1 // 余额
	ChangeTypePoints  = 2 // 积分

	// 变动行为
	ActionTypeIncome = 1 // 收入/增加
	ActionTypeExpend = 2 // 支出/减少

	// 来源类型
	SourceTypeAdmin      = 1 // 管理员调整
	SourceTypeRecharge   = 2 // 充值/签到
	SourceTypeConsume    = 3 // 消费
	SourceTypeRefund     = 4 // 退款/兑换
	SourceTypeCommission = 5 // 分销收益/过期
	SourceTypeWithdraw   = 6 // 提现
)

// UserChangeLog 用户变动记录模型
type UserChangeLog struct {
	ID           int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserID       int64      `gorm:"column:user_id;not null" json:"userId"`
	ChangeType   int        `gorm:"column:change_type;type:tinyint(1);not null" json:"changeType"`      // 变动数据类型(1-余额 2-积分)
	ActionType   int        `gorm:"column:action_type;type:tinyint(1);not null" json:"actionType"`      // 变动行为(1-收入/增加 2-支出/减少)
	Amount       float64    `gorm:"column:amount;type:decimal(10,2);not null" json:"amount"`            // 变动数量
	BeforeValue  float64    `gorm:"column:before_value;type:decimal(10,2);not null" json:"beforeValue"` // 变动前数值
	AfterValue   float64    `gorm:"column:after_value;type:decimal(10,2);not null" json:"afterValue"`   // 变动后数值
	SourceType   int        `gorm:"column:source_type;type:tinyint(1);not null" json:"sourceType"`      // 来源类型
	SourceID     *int64     `gorm:"column:source_id" json:"sourceId"`                                   // 来源ID
	Remark       string     `gorm:"column:remark;size:255" json:"remark"`                               // 备注
	OperatorID   *int64     `gorm:"column:operator_id" json:"operatorId"`                               // 操作人ID
	OperatorName string     `gorm:"column:operator_name;size:50" json:"operatorName"`                   // 操作人姓名
	ExpireTime   *time.Time `gorm:"column:expire_time" json:"expireTime"`                               // 过期时间(仅积分有效)
	CreateTime   time.Time  `gorm:"column:create_time;default:CURRENT_TIMESTAMP" json:"createTime"`     // 创建时间
}

// TableName 指定表名
func (u *UserChangeLog) TableName() string {
	return "user_change_log"
}

// Query 获取查询构建器
func (u *UserChangeLog) Query() *gorm.DB {
	return db.WnAdmin().Model(u).Where("1 = 1")
}

// GetByID 根据ID获取记录
func (u *UserChangeLog) GetByID(id int64) error {
	return u.Query().Where("id = ?", id).First(u).Error
}

// Create 创建记录
func (u *UserChangeLog) Create() error {
	return db.WnAdmin().Create(u).Error
}

// GetList 获取变动记录列表
func (u *UserChangeLog) GetList(offset, limit int, conditions map[string]interface{}) ([]UserChangeLog, int64, error) {
	var logs []UserChangeLog
	var total int64

	query := u.Query()

	// 添加查询条件
	for key, value := range conditions {
		if value != nil && value != "" {
			switch key {
			case "userId":
				query = query.Where("user_id = ?", value)
			case "changeType":
				query = query.Where("change_type = ?", value)
			case "actionType":
				query = query.Where("action_type = ?", value)
			case "sourceType":
				query = query.Where("source_type = ?", value)
			case "startTime":
				query = query.Where("create_time >= ?", value)
			case "endTime":
				query = query.Where("create_time <= ?", value)
			}
		}
	}

	// 获取总数
	query.Count(&total)

	// 获取列表
	err := query.Offset(offset).Limit(limit).Order("create_time DESC").Find(&logs).Error

	return logs, total, err
}

// GetChangeTypeText 获取变动数据类型文本
func (u *UserChangeLog) GetChangeTypeText() string {
	changeTypeMap := map[int]string{
		ChangeTypeBalance: "余额",
		ChangeTypePoints:  "积分",
	}
	return changeTypeMap[u.ChangeType]
}

// GetActionTypeText 获取变动行为文本
func (u *UserChangeLog) GetActionTypeText() string {
	if u.ChangeType == ChangeTypeBalance {
		// 余额
		actionTypeMap := map[int]string{
			ActionTypeIncome: "收入",
			ActionTypeExpend: "支出",
		}
		return actionTypeMap[u.ActionType]
	} else {
		// 积分
		actionTypeMap := map[int]string{
			ActionTypeIncome: "增加",
			ActionTypeExpend: "减少",
		}
		return actionTypeMap[u.ActionType]
	}
}

// GetSourceTypeText 获取来源类型文本
func (u *UserChangeLog) GetSourceTypeText() string {
	if u.ChangeType == ChangeTypeBalance {
		// 余额来源类型
		sourceTypeMap := map[int]string{
			SourceTypeAdmin:      "管理员调整",
			SourceTypeRecharge:   "充值",
			SourceTypeConsume:    "消费",
			SourceTypeRefund:     "退款",
			SourceTypeCommission: "分销收益",
			SourceTypeWithdraw:   "提现",
		}
		return sourceTypeMap[u.SourceType]
	} else {
		// 积分来源类型
		sourceTypeMap := map[int]string{
			SourceTypeAdmin:      "管理员调整",
			SourceTypeRecharge:   "签到",
			SourceTypeConsume:    "消费",
			SourceTypeRefund:     "兑换",
			SourceTypeCommission: "过期",
			SourceTypeWithdraw:   "退款",
		}
		return sourceTypeMap[u.SourceType]
	}
}
