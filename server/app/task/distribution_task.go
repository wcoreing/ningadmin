package task

import (
	"fmt"
	"time"
	"wnadmin/app/model/do/distribution"
	"wnadmin/app/provider/db"
	distributionService "wnadmin/app/service/distribution"
)

// UnfreezeIncome 解冻收益任务
func UnfreezeIncome() {
	fmt.Println("解冻收益任务")
	incomeService := distributionService.NewIncomeService()

	var incomes []distribution.DistributionIncome
	err := db.WnAdmin().Where("status = 0 AND unfreeze_time <= ?", time.Now()).Find(&incomes).Error
	if err != nil {
		return
	}

	for _, income := range incomes {
		incomeService.UnfreezeIncome(income.ID)
	}
}

// CheckIncomeProcess 检查收益处理状态
func CheckIncomeProcess() {
	incomeService := distributionService.NewIncomeService()
	incomeService.CheckIncomeProcess()
}
