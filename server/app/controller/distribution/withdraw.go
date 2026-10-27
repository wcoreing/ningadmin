package distribution

import (
	"wnadmin/app/common/response"
	"wnadmin/app/model/dto"
	"wnadmin/app/service/distribution"

	"github.com/gin-gonic/gin"
)

type WithdrawController struct {
	withdrawService *distribution.WithdrawService
}

func NewWithdrawController() *WithdrawController {
	return &WithdrawController{
		withdrawService: distribution.NewWithdrawService(),
	}
}

// List 获取提现列表
func (c *WithdrawController) List(ctx *gin.Context) {
	var req dto.WithdrawListReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(400, err.Error(), ctx)
		return
	}

	list, total, err := c.withdrawService.GetWithdrawList(&req)
	if err != nil {
		response.Fail(ctx)
		return
	}

	response.OkWithData(gin.H{
		"list":  list,
		"total": total,
	}, ctx)
}

// Detail 获取提现详情
func (c *WithdrawController) Detail(ctx *gin.Context) {
	id := ctx.Param("id")
	detail, err := c.withdrawService.GetWithdrawDetail(id)
	if err != nil {
		response.FailWithMessage(400, "获取提现详情失败", ctx)
		return
	}

	response.OkWithData(detail, ctx)
}

// Approve 审核通过
func (c *WithdrawController) Approve(ctx *gin.Context) {
	id := ctx.Param("id")
	err := c.withdrawService.ApproveWithdraw(id)
	if err != nil {
		response.FailWithMessage(400, "审核失败", ctx)
		return
	}

	response.Ok(ctx)
}

// Reject 审核拒绝
func (c *WithdrawController) Reject(ctx *gin.Context) {
	var req dto.WithdrawRejectReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(400, err.Error(), ctx)
		return
	}

	err := c.withdrawService.RejectWithdraw(&req)
	if err != nil {
		response.FailWithMessage(400, "审核失败", ctx)
		return
	}

	response.Ok(ctx)
}

// Pay 确认打款
func (c *WithdrawController) Pay(ctx *gin.Context) {
	var req dto.WithdrawPayReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(400, err.Error(), ctx)
		return
	}

	err := c.withdrawService.ConfirmPay(&req)
	if err != nil {
		if err.Error() == "用户可用余额不足" {
			response.FailWithMessage(400, err.Error(), ctx)
			return
		}
		response.FailWithMessage(400, "操作失败", ctx)
		return
	}

	response.Ok(ctx)
}
