package driver

import (
	"fmt"
	"time"
	"wnadmin/app/provider/db/model"

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
	}
	fmt.Println(m.Dsn())
	if db, err := gorm.Open(mysql.New(mysqlConfig)); err != nil {
		panic(err)
	} else {
		db.InstanceSet("gorm:table_options", "ENGINE=InnoDB")
		sqlDB, _ := db.DB()
		// 设置连接池
		sqlDB.SetMaxIdleConns(m.MaxIdleConns)   // 设置空闲连接池中连接的最大数量
		sqlDB.SetMaxOpenConns(m.MaxOpenConns)   // 设置打开数据库连接的最大数量
		sqlDB.SetConnMaxLifetime(time.Hour)     // 设置了连接可复用的最大时间
		sqlDB.SetConnMaxIdleTime(time.Hour * 2) // 设置空闲连接最大存活时间
		return db
	}
}
