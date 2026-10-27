package share

import (
	"wnadmin/app/common/response"
	"wnadmin/app/model/dto"
	"wnadmin/app/service/share"

	"github.com/gin-gonic/gin"
)

type ApplyController struct {
	applyService *share.ApplyService
}

func NewApplyController() *ApplyController {
	return &ApplyController{
		applyService: share.NewApplyService(),
	}
}

// List 获取分享申请列表
func (c *ApplyController) List(ctx *gin.Context) {
	var req dto.ShareApplyListReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(400, err.Error(), ctx)
		return
	}

	list, total, err := c.applyService.GetApplyList(&req)
	if err != nil {
		response.FailWithMessage(500, err.Error(), ctx)
		return
	}

	response.OkWithDetailed(gin.H{"list": list, "total": total}, "获取成功", ctx)
}

// Approve 审核通过
func (c *ApplyController) Approve(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.applyService.ApproveApply(id); err != nil {
		response.FailWithMessage(500, err.Error(), ctx)
		return
	}
	response.OkWithMessage("审核通过", ctx)
}

// Reject 审核拒绝
func (c *ApplyController) Reject(ctx *gin.Context) {
	var req dto.ShareApplyRejectReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(400, err.Error(), ctx)
		return
	}

	if err := c.applyService.RejectApply(&req); err != nil {
		response.FailWithMessage(500, err.Error(), ctx)
		return
	}
	response.OkWithMessage("审核拒绝", ctx)
}

// Detail 获取申请详情
func (c *ApplyController) Detail(ctx *gin.Context) {
	id := ctx.Param("id")
	detail, err := c.applyService.GetApplyDetail(id)
	if err != nil {
		response.FailWithMessage(500, err.Error(), ctx)
		return
	}
	response.OkWithDetailed(detail, "获取成功", ctx)
}
