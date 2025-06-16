package db

import (
	"wnadmin/app/config"
	"wnadmin/app/provider/db/driver"

	"gorm.io/gorm"
)

type DBModel struct {
	WnAdmin *gorm.DB
}

var DB DBModel

func init() {
	// 初始化数据库连接
	DB.WnAdmin = driver.GormMysqlByConfig(config.DB.Mysql["default"])
}

// WnAdmin 获取数据库实例
func WnAdmin() *gorm.DB {
	return driver.GetDB()
}

// IsConnected 检查数据库连接状态
func IsConnected() bool {
	return driver.IsConnected()
}

// ReconnectDB 重新连接数据库
func ReconnectDB() error {
	return driver.ReconnectDB()
}

// CloseDB 关闭数据库连接
func CloseDB() error {
	return driver.CloseDB()
}
