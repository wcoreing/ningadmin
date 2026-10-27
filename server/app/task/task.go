package task

import (
	"fmt"
	"time"

	croner "github.com/robfig/cron/v3"
)

var cronTab *croner.Cron

/*
Field name   | Mandatory? | Allowed values  | Allowed special characters
----------   | ---------- | --------------  | --------------------------
Seconds      | Yes        | 0-59            | * / , -
Minutes      | Yes        | 0-59            | * / , -
Hours        | Yes        | 0-23            | * / , -
Day of month | Yes        | 1-31            | * / , - ?
Month        | Yes        | 1-12 or JAN-DEC | * / , -
Day of week  | Yes        | 0-6 or SUN-SAT  | * / , - ?
*/
func InitTask() {
	// 初始化 cron 实例
	cronTab = croner.New()

	// 添加定时任务
	cronTab.AddFunc("0 0 * * *", func() {
		// ... task code ...
	})

	// 每天凌晨1点执行自动完成订单任务
	cronTab.AddFunc("0 1 * * *", AutoCompleteOrder)

	// 每小时执行一次解冻收益任务
	cronTab.AddFunc("0 * * * *", UnfreezeIncome)

	// 每小时执行一次收益处理状态检查
	cronTab.AddFunc("24 * * * *", CheckIncomeProcess)

	// 每分钟执行一次任务
	cronTab.AddFunc("*/2 * * * *", func() {
		fmt.Println("每分钟执行一次任务", time.Now().Format("2006-01-02 15:04:05"))
	})

	// 启动定时任务
	cronTab.Start()
}
