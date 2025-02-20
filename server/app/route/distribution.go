package route

import (
	"wnadmin/app/controller/distribution"

	"github.com/gin-gonic/gin"
)

func RegisterDistributionRoutes(r *gin.RouterGroup) {
	distributionGroup := r.Group("/distribution")

	{
		// 分销设置
		settingGroup := distributionGroup.Group("/setting")
		{
			settingController := distribution.NewSettingController()
			settingGroup.GET("", settingController.GetSetting)    // 获取分销设置
			settingGroup.PUT("", settingController.UpdateSetting) // 更新分销设置
		}

		// 分销申请
		applyController := distribution.NewApplyController()
		distributionGroup.GET("/apply/list", applyController.List)
		distributionGroup.POST("/apply/approve/:id", applyController.Approve)
		distributionGroup.POST("/apply/reject", applyController.Reject)
		distributionGroup.GET("/apply/detail/:id", applyController.Detail)

		// 提现管理
		withdrawController := distribution.NewWithdrawController()
		distributionGroup.GET("/withdraw/list", withdrawController.List)
		distributionGroup.GET("/withdraw/detail/:id", withdrawController.Detail)
		distributionGroup.POST("/withdraw/approve/:id", withdrawController.Approve)
		distributionGroup.POST("/withdraw/reject", withdrawController.Reject)
		distributionGroup.POST("/withdraw/pay", withdrawController.Pay)
		// 分销收益相关路由

		incomeController := distribution.NewIncomeController()
		incomeRouter := distributionGroup.Group("/income")
		{
			incomeRouter.GET("/list", incomeController.GetIncomList)  // 分页查询
			incomeRouter.POST("/unfreeze", incomeController.Unfreeze) // 解冻收益
		}

	}
}
