package SystemDB

import (
	"encoding/json"
	"time"
	"wnadmin/app/provider/db"

	"gorm.io/gorm"
)

// PopupDO 弹窗数据对象
type PopupDO struct {
	ID           int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Title        string     `gorm:"column:title;type:varchar(100);default:''" json:"title"`
	Content      string     `gorm:"column:content;type:text;not null" json:"content"`
	TitleStyle   string     `gorm:"column:title_style;type:json" json:"title_style"`
	Buttons      string     `gorm:"column:buttons;type:json" json:"buttons"`
	TriggerType  string     `gorm:"column:trigger_type;type:varchar(20);not null;default:'manual'" json:"trigger_type"`
	TriggerRule  string     `gorm:"column:trigger_rule;type:json" json:"trigger_rule"`
	Priority     int        `gorm:"column:priority;type:int;not null;default:0" json:"priority"`
	AutoClose    int        `gorm:"column:auto_close;type:int;not null;default:0" json:"auto_close"`
	Delay        int        `gorm:"column:delay;type:int;not null;default:0" json:"delay"`
	MaskClosable int        `gorm:"column:mask_closable;type:tinyint(1);not null;default:1" json:"mask_closable"`
	StartTime    *time.Time `gorm:"column:start_time;type:datetime" json:"start_time"`
	EndTime      *time.Time `gorm:"column:end_time;type:datetime" json:"end_time"`
	Status       int        `gorm:"column:status;type:tinyint(1);not null;default:1" json:"status"`
	ShowCount    int64      `gorm:"column:show_count;type:bigint;not null;default:0" json:"show_count"`
	ClickCount   int64      `gorm:"column:click_count;type:bigint;not null;default:0" json:"click_count"`
	CreateTime   time.Time  `gorm:"column:create_time;autoCreateTime" json:"create_time"`
	UpdateTime   time.Time  `gorm:"column:update_time;autoUpdateTime" json:"update_time"`
	IsDelete     int        `gorm:"column:is_delete;type:tinyint(1);not null;default:0" json:"is_delete"`
}

// TableName 表名
func (PopupDO) TableName() string {
	return "popup"
}

// Query 查询构造器
func (p *PopupDO) Query() *gorm.DB {
	return db.WnAdmin().Model(p).Where("is_delete = ?", 0)
}

// Create 创建
func (p *PopupDO) Create() error {
	return db.WnAdmin().Create(p).Error
}

// Update 更新
func (p *PopupDO) Update() error {
	return db.WnAdmin().Model(p).Updates(p).Error
}

// Delete 删除
func (p *PopupDO) Delete() error {
	return db.WnAdmin().Model(p).Where("id = ?", p.ID).Update("is_delete", 1).Error
}

// First 查询单条记录
func (p *PopupDO) First() error {
	return p.Query().Where("id = ?", p.ID).First(p).Error
}

// GetTitleStyleMap 获取标题样式映射
func (p *PopupDO) GetTitleStyleMap() map[string]interface{} {
	if p.TitleStyle == "" {
		return nil
	}
	var result map[string]interface{}
	json.Unmarshal([]byte(p.TitleStyle), &result)
	return result
}

// GetButtonsMap 获取按钮配置映射
func (p *PopupDO) GetButtonsMap() []map[string]interface{} {
	if p.Buttons == "" {
		return nil
	}
	var result []map[string]interface{}
	json.Unmarshal([]byte(p.Buttons), &result)
	return result
}

// GetTriggerRuleMap 获取触发规则映射
func (p *PopupDO) GetTriggerRuleMap() map[string]interface{} {
	if p.TriggerRule == "" {
		return nil
	}
	var result map[string]interface{}
	json.Unmarshal([]byte(p.TriggerRule), &result)
	return result
}

// SetTitleStyle 设置标题样式
func (p *PopupDO) SetTitleStyle(style map[string]interface{}) error {
	if style == nil {
		p.TitleStyle = ""
		return nil
	}
	data, err := json.Marshal(style)
	if err != nil {
		return err
	}
	p.TitleStyle = string(data)
	return nil
}

// SetButtons 设置按钮配置
func (p *PopupDO) SetButtons(buttons []map[string]interface{}) error {
	if buttons == nil {
		p.Buttons = ""
		return nil
	}
	data, err := json.Marshal(buttons)
	if err != nil {
		return err
	}
	p.Buttons = string(data)
	return nil
}

// SetTriggerRule 设置触发规则
func (p *PopupDO) SetTriggerRule(rule map[string]interface{}) error {
	if rule == nil {
		p.TriggerRule = ""
		return nil
	}
	data, err := json.Marshal(rule)
	if err != nil {
		return err
	}
	p.TriggerRule = string(data)
	return nil
}
