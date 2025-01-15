package db

import (
	"fmt"
	"gorm.io/gorm"
	"wnadmin/app/config"
	"wnadmin/app/provider/db/driver"
)

type DBModel struct {
	WnAdmin *gorm.DB
}

var DB DBModel

func init() {
	fmt.Println(config.DB.Mysql)
	DB.WnAdmin = driver.GormMysqlByConfig(config.DB.Mysql["studentactive"])

}
