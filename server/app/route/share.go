package route

import (
	"wnadmin/app/controller/share"
	"wnadmin/app/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterShareRoutes(r *gin.RouterGroup) {
	shareGroup := r.Group("/share")

	{
		// 分享申请
		applyController := share.NewApplyController()
		shareGroup.GET("/apply/list", middleware.WithSelect("分享申请", "查询申请列表"), applyController.List)
		shareGroup.POST("/apply/approve/:id", middleware.WithUpdate("分享申请", "审批申请"), applyController.Approve)
		shareGroup.POST("/apply/reject", middleware.WithUpdate("分享申请", "拒绝申请"), applyController.Reject)
		shareGroup.GET("/apply/detail/:id", middleware.WithSelect("分享申请", "查询申请详情"), applyController.Detail)

		// 提现管理
		withdrawController := share.NewWithdrawController()
		shareGroup.GET("/withdraw/list", middleware.WithSelect("分享提现", "查询提现列表"), withdrawController.List)
		shareGroup.GET("/withdraw/detail/:id", middleware.WithSelect("分享提现", "查询提现详情"), withdrawController.Detail)
		shareGroup.POST("/withdraw/approve/:id", middleware.WithUpdate("分享提现", "审批提现"), withdrawController.Approve)
		shareGroup.POST("/withdraw/reject", middleware.WithUpdate("分享提现", "拒绝提现"), withdrawController.Reject)
		shareGroup.POST("/withdraw/pay", middleware.WithUpdate("分享提现", "提现打款"), withdrawController.Pay)

		// 分享收益相关路由
		incomeController := share.NewIncomeController()
		incomeRouter := shareGroup.Group("/income")
		{
			incomeRouter.GET("/list", middleware.WithSelect("分享收益", "查询收益列表"), incomeController.GetIncomList)         // 分页查询
			incomeRouter.POST("/unfreeze", middleware.WithUpdate("分享收益", "解冻收益"), incomeController.Unfreeze)          // 解冻收益
			incomeRouter.POST("/settle", middleware.WithUpdate("分享收益", "结算收益"), incomeController.Settle)              // 单个结算
			incomeRouter.POST("/batch-settle", middleware.WithUpdate("分享收益", "批量结算收益"), incomeController.BatchSettle) // 批量结算
		}

	}
}
