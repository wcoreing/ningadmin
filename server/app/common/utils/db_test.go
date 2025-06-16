package utils

import (
	"database/sql"
	"fmt"
	"time"
	"wnadmin/app/provider/db"
	"wnadmin/app/provider/mylog"

	_ "github.com/go-sql-driver/mysql"
)

// TestDBConnection 测试数据库连接
func TestDBConnection() error {
	if db.DB.WnAdmin == nil {
		return fmt.Errorf("数据库实例未初始化")
	}

	sqlDB, err := db.WnAdmin().DB()
	if err != nil {
		return fmt.Errorf("获取数据库连接失败: %v", err)
	}

	// 测试连接
	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("数据库连接测试失败: %v", err)
	}

	// 获取连接池状态
	stats := sqlDB.Stats()
	mylog.Log.WithField("max_open_connections", stats.MaxOpenConnections).
		WithField("open_connections", stats.OpenConnections).
		WithField("in_use", stats.InUse).
		WithField("idle", stats.Idle).
		Info("数据库连接池状态")

	return nil
}

// MonitorDBConnection 监控数据库连接状态
func MonitorDBConnection(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		if err := TestDBConnection(); err != nil {
			mylog.Log.WithField("error", err.Error()).Error("数据库连接监控检测到问题")
		} else {
			mylog.Log.Debug("数据库连接监控正常")
		}
	}
}

// CheckMySQLVersion 检查MySQL版本
func CheckMySQLVersion() (string, error) {
	if db.DB.WnAdmin == nil {
		return "", fmt.Errorf("数据库实例未初始化")
	}

	var version string
	err := db.WnAdmin().Raw("SELECT VERSION()").Scan(&version).Error
	if err != nil {
		return "", fmt.Errorf("获取MySQL版本失败: %v", err)
	}

	return version, nil
}

// TestRawConnection 测试原生SQL连接
func TestRawConnection(dsn string) error {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("打开数据库连接失败: %v", err)
	}
	defer db.Close()

	// 设置连接参数
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Hour)

	// 测试连接
	if err := db.Ping(); err != nil {
		return fmt.Errorf("数据库连接测试失败: %v", err)
	}

	// 执行简单查询
	var result int
	err = db.QueryRow("SELECT 1").Scan(&result)
	if err != nil {
		return fmt.Errorf("执行测试查询失败: %v", err)
	}

	if result != 1 {
		return fmt.Errorf("测试查询结果异常: 期望1，实际%d", result)
	}

	mylog.Log.Info("原生数据库连接测试成功")
	return nil
}
