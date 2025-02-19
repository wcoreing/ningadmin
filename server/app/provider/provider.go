package provider

import (
	_ "wnadmin/app/provider/conf"
	_ "wnadmin/app/provider/db"
	"wnadmin/app/provider/mylog"
	"wnadmin/app/task"
)

func Register() {
	mylog.Log.Info("start")
	task.InitTask()
}
