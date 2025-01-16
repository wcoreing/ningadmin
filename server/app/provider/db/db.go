package db

import (
	"fmt"
	"wnadmin/app/config"
	"wnadmin/app/provider/db/driver"

	"gorm.io/gorm"
)

type DBModel struct {
	WnAdmin *gorm.DB
}

var DB DBModel

func init() {
	fmt.Println(config.DB.Mysql)
	DB.WnAdmin = driver.GormMysqlByConfig(config.DB.Mysql["default"])

}
