package provider

import (
	_ "wnadmin/app/provider/conf"
	_ "wnadmin/app/provider/db"
	"wnadmin/app/provider/mylog"
	_ "wnadmin/app/provider/mylog"
)

func Register() {
	mylog.Log.Info("start")
}
