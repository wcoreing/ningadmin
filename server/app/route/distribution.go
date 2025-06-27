package route

import (
	"wnadmin/app/controller/distribution"
	"wnadmin/app/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterDistributionRoutes(r *gin.RouterGroup) {
	distributionGroup := r.Group("/distribution")

	{
		// 分销设置
		settingGroup := distributionGroup.Group("/setting")
		{
			settingController := distribution.NewSettingController()
			settingGroup.GET("", middleware.WithSelect("分销设置", "查询分销设置"), settingController.GetSetting)    // 获取分销设置
			settingGroup.PUT("", middleware.WithUpdate("分销设置", "修改分销设置"), settingController.UpdateSetting) // 更新分销设置
		}

		// 分销申请
		applyController := distribution.NewApplyController()
		distributionGroup.GET("/apply/list", middleware.WithSelect("分销申请", "查询申请列表"), applyController.List)
		distributionGroup.POST("/apply/approve/:id", middleware.WithUpdate("分销申请", "审批申请"), applyController.Approve)
		distributionGroup.POST("/apply/reject", middleware.WithUpdate("分销申请", "拒绝申请"), applyController.Reject)
		distributionGroup.GET("/apply/detail/:id", middleware.WithSelect("分销申请", "查询申请详情"), applyController.Detail)

		// 提现管理
		withdrawController := distribution.NewWithdrawController()
		distributionGroup.GET("/withdraw/list", middleware.WithSelect("提现管理", "查询提现列表"), withdrawController.List)
		distributionGroup.GET("/withdraw/detail/:id", middleware.WithSelect("提现管理", "查询提现详情"), withdrawController.Detail)
		distributionGroup.POST("/withdraw/approve/:id", middleware.WithUpdate("提现管理", "审批提现"), withdrawController.Approve)
		distributionGroup.POST("/withdraw/reject", middleware.WithUpdate("提现管理", "拒绝提现"), withdrawController.Reject)
		distributionGroup.POST("/withdraw/pay", middleware.WithUpdate("提现管理", "提现打款"), withdrawController.Pay)
		// 分销收益相关路由

		incomeController := distribution.NewIncomeController()
		incomeRouter := distributionGroup.Group("/income")
		{
			incomeRouter.GET("/list", middleware.WithSelect("分销收益", "查询收益列表"), incomeController.GetIncomList)         // 分页查询
			incomeRouter.POST("/unfreeze", middleware.WithUpdate("分销收益", "解冻收益"), incomeController.Unfreeze)          // 解冻收益
			incomeRouter.POST("/settle", middleware.WithUpdate("分销收益", "结算收益"), incomeController.Settle)              // 单个结算
			incomeRouter.POST("/batch-settle", middleware.WithUpdate("分销收益", "批量结算收益"), incomeController.BatchSettle) // 批量结算
		}

	}
}
