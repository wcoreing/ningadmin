package task

import (
	"fmt"

	"github.com/robfig/cron"
)

var cronTab *cron.Cron

func InitTask() {
	// 每天凌晨1点执行自动完成订单任务
	cronTab.AddFunc("0 1 * * *", AutoCompleteOrder)

	// 每小时执行一次解冻收益任务
	cronTab.AddFunc("0 * * * *", UnfreezeIncome)

	// 每小时执行一次收益处理状态检查
	cronTab.AddFunc("0 * * * *", CheckIncomeProcess)

	// 每小时执行一次订单自动收货任务
	cronTab.AddFunc("* * * * *", func() {
		fmt.Println("每分钟执行一次任务")
	})
}
