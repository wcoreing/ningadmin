package driver

import (
	"fmt"
	"sync"
	"time"
	"wnadmin/app/config"
	"wnadmin/app/provider/db/model"

	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	dbInstance *gorm.DB
	dbMutex    sync.RWMutex
	dbConfig   model.Mysql
)

// GormMysqlByConfig 初始化Mysql数据库
func GormMysqlByConfig(m model.Mysql) *gorm.DB {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	// 保存配置
	dbConfig = m

	if dbInstance != nil {
		// 检查现有连接是否可用
		if sqlDB, err := dbInstance.DB(); err == nil {
			if err := sqlDB.Ping(); err == nil {
				return dbInstance
			}
		}
		// 关闭无效连接
		closeConnection()
	}

	// 创建新连接
	dbInstance = createNewConnection(m)
	return dbInstance
}

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	dbMutex.RLock()
	if dbInstance != nil {
		dbMutex.RUnlock()
		return dbInstance
	}
	dbMutex.RUnlock()

	// 如果没有实例，使用默认配置创建
	return GormMysqlByConfig(config.DB.Mysql["default"])
}

// createNewConnection 创建新的数据库连接
func createNewConnection(m model.Mysql) *gorm.DB {
	if m.Dbname == "" {
		panic("数据库名称不能为空")
	}

	// MySQL配置
	mysqlConfig := mysql.Config{
		DSN:                       m.Dsn(),
		DefaultStringSize:         191,
		SkipInitializeWithVersion: false,
		DisableDatetimePrecision:  true,
		DontSupportRenameIndex:    true,
		DontSupportRenameColumn:   true,
	}

	// 日志级别
	var logLevel logger.LogLevel
	switch m.LogMode {
	case "silent":
		logLevel = logger.Silent
	case "error":
		logLevel = logger.Error
	case "warn":
		logLevel = logger.Warn
	case "info":
		logLevel = logger.Info
	default:
		logLevel = logger.Error // 默认只显示错误
	}

	// GORM配置
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
	}

	// 连接数据库，带重试
	db, err := connectWithRetry(mysqlConfig, gormConfig, 3)
	if err != nil {
		panic(fmt.Sprintf("数据库连接失败: %v", err))
	}

	// 配置连接池
	setupConnectionPool(db, m)

	fmt.Printf("数据库连接成功: %s\n", m.Dbname)
	return db
}

// connectWithRetry 带重试的连接
func connectWithRetry(mysqlConfig mysql.Config, gormConfig *gorm.Config, maxRetries int) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	for i := 0; i < maxRetries; i++ {
		db, err = gorm.Open(mysql.New(mysqlConfig), gormConfig)
		if err == nil {
			// 测试连接
			if sqlDB, dbErr := db.DB(); dbErr == nil {
				if pingErr := sqlDB.Ping(); pingErr == nil {
					return db, nil
				} else {
					err = pingErr
				}
			} else {
				err = dbErr
			}
		}

		if i < maxRetries-1 {
			fmt.Printf("数据库连接失败，%d秒后重试: %v\n", (i+1)*2, err)
			time.Sleep(time.Duration((i+1)*2) * time.Second)
		}
	}

	return nil, err
}

// setupConnectionPool 设置连接池
func setupConnectionPool(db *gorm.DB, m model.Mysql) {
	sqlDB, err := db.DB()
	if err != nil {
		return
	}

	// 连接池配置
	maxOpen := m.MaxOpenConns
	maxIdle := m.MaxIdleConns

	if maxOpen <= 0 {
		maxOpen = 50 // 适中的默认值
	}
	if maxIdle <= 0 {
		maxIdle = 10
	}

	sqlDB.SetMaxOpenConns(maxOpen)
	sqlDB.SetMaxIdleConns(maxIdle)
	sqlDB.SetConnMaxLifetime(time.Hour * 4)    // 4小时生命周期
	sqlDB.SetConnMaxIdleTime(time.Minute * 30) // 30分钟空闲超时
}

// ReconnectDB 重新连接数据库
func ReconnectDB() error {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	// 关闭现有连接
	closeConnection()

	// 重新创建连接
	newDB := createNewConnection(dbConfig)
	if newDB == nil {
		return fmt.Errorf("重连失败")
	}

	dbInstance = newDB
	fmt.Println("数据库重连成功")
	return nil
}

// IsConnected 检查连接状态
func IsConnected() bool {
	dbMutex.RLock()
	defer dbMutex.RUnlock()

	if dbInstance == nil {
		return false
	}

	sqlDB, err := dbInstance.DB()
	if err != nil {
		return false
	}

	return sqlDB.Ping() == nil
}

// CloseDB 关闭数据库连接
func CloseDB() error {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	return closeConnection()
}

// closeConnection 内部关闭连接方法
func closeConnection() error {
	if dbInstance != nil {
		if sqlDB, err := dbInstance.DB(); err == nil {
			err = sqlDB.Close()
			dbInstance = nil
			return err
		}
		dbInstance = nil
	}
	return nil
}
