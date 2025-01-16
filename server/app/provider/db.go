package provider

import (
	"fmt"
	"log"
	"sync"
	"time"
	"wnadmin/app/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var (
	db   *gorm.DB
	once sync.Once
)

// initDB 初始化数据库连接
func initDB() {
	mysqlConfig := config.DB.Mysql["default"]
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s",
		mysqlConfig.Username,
		mysqlConfig.Password,
		mysqlConfig.Path,
		mysqlConfig.Port,
		mysqlConfig.Dbname,
		mysqlConfig.Config,
	)

	var logLevel logger.LogLevel
	switch mysqlConfig.LogMode {
	case "silent":
		logLevel = logger.Silent
	case "error":
		logLevel = logger.Error
	case "warn":
		logLevel = logger.Warn
	case "info":
		logLevel = logger.Info
	default:
		logLevel = logger.Info
	}

	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数表名
		},
	}

	var err error
	db, err = gorm.Open(mysql.Open(dsn), gormConfig)
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("获取数据库实例失败: %v", err)
	}

	// 设置连接池
	sqlDB.SetMaxIdleConns(mysqlConfig.MaxIdleConns)
	sqlDB.SetMaxOpenConns(mysqlConfig.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Hour)
}

// GetDB 获取数据库连接实例
func GetDB() *gorm.DB {
	once.Do(func() {
		initDB()
	})
	return db
}

// CloseDB 关闭数据库连接
func CloseDB() {
	if db != nil {
		sqlDB, err := db.DB()
		if err != nil {
			log.Printf("获取数据库实例失败: %v", err)
			return
		}
		err = sqlDB.Close()
		if err != nil {
			log.Printf("关闭数据库连接失败: %v", err)
		}
	}
}
