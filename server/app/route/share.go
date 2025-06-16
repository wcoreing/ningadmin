package route

import (
	"wnadmin/app/controller/share"

	"github.com/gin-gonic/gin"
)

func RegisterShareRoutes(r *gin.RouterGroup) {
	shareGroup := r.Group("/share")

	{
		// 分享申请
		applyController := share.NewApplyController()
		shareGroup.GET("/apply/list", applyController.List)
		shareGroup.POST("/apply/approve/:id", applyController.Approve)
		shareGroup.POST("/apply/reject", applyController.Reject)
		shareGroup.GET("/apply/detail/:id", applyController.Detail)

		// 提现管理
		withdrawController := share.NewWithdrawController()
		shareGroup.GET("/withdraw/list", withdrawController.List)
		shareGroup.GET("/withdraw/detail/:id", withdrawController.Detail)
		shareGroup.POST("/withdraw/approve/:id", withdrawController.Approve)
		shareGroup.POST("/withdraw/reject", withdrawController.Reject)
		shareGroup.POST("/withdraw/pay", withdrawController.Pay)

		// 分享收益相关路由
		incomeController := share.NewIncomeController()
		incomeRouter := shareGroup.Group("/income")
		{
			incomeRouter.GET("/list", incomeController.GetIncomList)  // 分页查询
			incomeRouter.POST("/unfreeze", incomeController.Unfreeze) // 解冻收益
		}

	}
}
