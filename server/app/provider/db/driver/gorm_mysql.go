package driver

import (
	"fmt"
	"time"
	"wnadmin/app/config"
	"wnadmin/app/provider/db/model"

	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// GormMysqlByConfig 初始化Mysql数据库用过传入配置
func GormMysqlByConfig(m model.Mysql) *gorm.DB {
	if m.Dbname == "" {
		return nil
	}
	mysqlConfig := mysql.Config{
		DSN:                       m.Dsn(), // DSN data source name
		DefaultStringSize:         191,     // string 类型字段的默认长度
		SkipInitializeWithVersion: false,   // 根据版本自动配置
		DisableDatetimePrecision:  true,    // 禁用 datetime 精度
		DontSupportRenameIndex:    true,    // 重命名索引优化
		DontSupportRenameColumn:   true,    // 重命名列优化
	}
	fmt.Println(m.Dsn())
	mysqlConfig2 := config.DB.Mysql["default"]
	var logLevel logger.LogLevel
	switch mysqlConfig2.LogMode {
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
		SkipDefaultTransaction: true, // 禁用默认事务
		PrepareStmt:            true, // 预编译缓存
	}
	if db, err := gorm.Open(mysql.New(mysqlConfig), gormConfig); err != nil {
		panic(err)
	} else {
		db.InstanceSet("gorm:table_options", "ENGINE=InnoDB")
		sqlDB, _ := db.DB()
		// 设置连接池
		sqlDB.SetMaxIdleConns(m.MaxIdleConns)     // 空闲连接池大小
		sqlDB.SetMaxOpenConns(m.MaxOpenConns)     // 最大连接数
		sqlDB.SetConnMaxLifetime(time.Hour * 4)   // 连接最大生命周期
		sqlDB.SetConnMaxIdleTime(time.Minute * 5) // 空闲连接超时时间

		go func() {
			for {
				time.Sleep(time.Second * 30) // 每30秒检查一次
				err := sqlDB.Ping()
				if err != nil {
					// 发现连接异常，尝试重连
					newDB, err := gorm.Open(mysql.New(mysqlConfig), gormConfig)
					if err == nil {
						// 重连成功，更新连接
						db = newDB
					}
				}
			}
		}()

		return db
	}
}
