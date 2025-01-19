package SystemDB

import (
	"time"
	"wnadmin/app/provider/db"

	"gorm.io/gorm"
)

// ConfigDO 系统配置
type ConfigDO struct {
	ID          int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`           // 配置ID
	ModuleKey   string    `gorm:"column:module_key;size:100;not null" json:"moduleKey"`   // 模块标识
	ModuleName  string    `gorm:"column:module_name;size:100;not null" json:"moduleName"` // 模块标识
	ConfigName  string    `gorm:"column:config_name;size:100;not null" json:"configName"` // 配置名称
	ConfigKey   string    `gorm:"column:config_key;size:100;not null" json:"configKey"`   // 配置键名
	ConfigType  string    `gorm:"column:config_type;size:10;not null" json:"configType"`  // 配置类型
	ConfigValue string    `gorm:"column:config_value;type:text" json:"configValue"`       // 配置值
	Remark      string    `gorm:"column:remark;size:500" json:"remark"`                   // 备注
	Status      string    `gorm:"column:status;size:1;default:0" json:"status"`           // 状态（0正常 1停用）
	DelFlag     string    `gorm:"column:del_flag;size:1;default:0" json:"delFlag"`        // 删除标志（0代表存在 2代表删除）
	CreateTime  time.Time `gorm:"column:create_time;autoCreateTime" json:"createTime"`    // 创建时间
	UpdateTime  time.Time `gorm:"column:update_time;autoUpdateTime" json:"updateTime"`    // 更新时间
}

// TableName 表名
func (ConfigDO) TableName() string {
	return "sys_config"
}

// Query 查询构造器
func (c *ConfigDO) Query() *gorm.DB {
	return db.DB.WnAdmin.Model(c).Where("del_flag = ?", "0")
}

// First 获取第一条记录
func (c *ConfigDO) First() error {
	return c.Query().Where("id = ?", c.ID).First(c).Error
}

// Create 创建记录
func (c *ConfigDO) Create() error {
	return db.DB.WnAdmin.Create(c).Error
}

// Update 更新记录
func (c *ConfigDO) Update() error {
	return db.DB.WnAdmin.Model(c).Updates(c).Error
}

// Delete 删除记录
func (c *ConfigDO) Delete() error {
	return db.DB.WnAdmin.Model(c).Where("id = ?", c.ID).Update("del_flag", "2").Error
}
