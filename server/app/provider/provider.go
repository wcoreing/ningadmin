package provider

import (
	_ "wnadmin/app/provider/conf"
	_ "wnadmin/app/provider/db"
	"wnadmin/app/provider/mylog"
	"wnadmin/app/task"
)

func Register() {
	mylog.Log.Info("系统启动中...")
	// 初始化任务
	task.InitTask()

	mylog.Log.Info("系统启动完成")
}
