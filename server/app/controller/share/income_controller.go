package share

import (
	"wnadmin/app/common/response"
	"wnadmin/app/model/dto"
	"wnadmin/app/service/share"

	"github.com/gin-gonic/gin"
)

type IncomeController struct {
	incomeService *share.IncomeService
}

func NewIncomeController() *IncomeController {
	return &IncomeController{
		incomeService: share.NewIncomeService(),
	}
}

// GetPage 分页获取分享收益列表
func (c *IncomeController) GetIncomList(ctx *gin.Context) {
	var req dto.ShareIncomeListReq
	if err := ctx.ShouldBind(&req); err != nil {
		response.FailWithMessage(200, err.Error(), ctx)
		return
	}

	list, total, err := c.incomeService.GetIncomList(&req)
	if err != nil {
		response.FailWithMessage(200, err.Error(), ctx)
		return
	}

	response.OkWithData(map[string]interface{}{
		"list":  list,
		"total": total,
	}, ctx)
}

// Unfreeze 解冻分享收益
func (c *IncomeController) Unfreeze(ctx *gin.Context) {
	var req dto.ShareIncomeUnfreezeReq
	if err := ctx.ShouldBind(&req); err != nil {
		response.FailWithMessage(200, err.Error(), ctx)
		return
	}

	err := c.incomeService.Unfreeze(&req)
	if err != nil {
		response.FailWithMessage(200, err.Error(), ctx)
		return
	}

	response.Ok(ctx)
}

// Settle 结算单个收益
func (c *IncomeController) Settle(ctx *gin.Context) {
	var req struct {
		ID uint `json:"id" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(201, "参数错误: "+err.Error(), ctx)
		return
	}

	// 获取操作人
	operator := ctx.GetString("userName")
	if operator == "" {
		operator = "admin"
	}

	err := c.incomeService.SettleIncome(req.ID, operator)
	if err != nil {
		response.FailWithMessage(201, "结算失败: "+err.Error(), ctx)
		return
	}

	response.OkWithMessage("结算成功", ctx)
}

// BatchSettle 批量结算收益
func (c *IncomeController) BatchSettle(ctx *gin.Context) {
	var req struct {
		IDs []uint `json:"ids" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(201, "参数错误: "+err.Error(), ctx)
		return
	}

	if len(req.IDs) == 0 {
		response.FailWithMessage(201, "请选择要结算的收益记录", ctx)
		return
	}

	// 获取操作人
	operator := ctx.GetString("userName")
	if operator == "" {
		operator = "admin"
	}

	err := c.incomeService.BatchSettleIncome(req.IDs, operator)
	if err != nil {
		response.FailWithMessage(201, "批量结算失败: "+err.Error(), ctx)
		return
	}

	response.OkWithMessage("批量结算成功", ctx)
}
