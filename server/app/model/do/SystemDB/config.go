package SystemDB

import (
	"time"
	"wnadmin/app/provider/db"

	"gorm.io/gorm"
)

// ConfigDO 系统配置
type ConfigDO struct {
	ID          int64     `json:"id" gorm:"primaryKey"`
	ConfigName  string    `json:"configName" gorm:"column:config_name;comment:配置名称"`
	ConfigKey   string    `json:"configKey" gorm:"column:config_key;comment:配置键"`
	ConfigValue string    `json:"configValue" gorm:"column:config_value;comment:配置值"`
	ConfigType  string    `json:"configType" gorm:"column:config_type;comment:配置类型"`
	ModuleKey   string    `json:"moduleKey" gorm:"column:module_key;comment:模块标识"`
	ModuleName  string    `json:"moduleName,omitempty" gorm:"column:module_name;comment:模块名称"`
	Remark      string    `json:"remark" gorm:"comment:备注"`
	CreateTime  time.Time `json:"createTime" gorm:"column:create_time;autoCreateTime"`
	UpdateTime  time.Time `json:"updateTime" gorm:"column:update_time;autoUpdateTime"`
}

// TableName 表名
func (ConfigDO) TableName() string {
	return "sys_config"
}

// Query 获取查询构建器
func (c *ConfigDO) Query() *gorm.DB {
	return db.DB.WnAdmin.Model(c)
}

// Create 创建记录
func (c *ConfigDO) Create() error {
	return c.Query().Create(c).Error
}

// Update 更新记录
func (c *ConfigDO) Update() error {
	return c.Query().Updates(c).Error
}

// Delete 删除记录
func (c *ConfigDO) Delete() error {
	return c.Query().Delete(c).Error
}

// First 获取第一条记录
func (c *ConfigDO) First() error {
	return c.Query().First(c).Error
}

// Find 查询多条记录
func (c *ConfigDO) Find(dest interface{}) error {
	return c.Query().Find(dest).Error
}
