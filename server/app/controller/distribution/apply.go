package distribution

import (
	"wnadmin/app/model/dto"
	"wnadmin/app/service/distribution"

	"github.com/gin-gonic/gin"
)

type ApplyController struct {
	applyService *distribution.ApplyService
}

func NewApplyController() *ApplyController {
	return &ApplyController{
		applyService: distribution.NewApplyService(),
	}
}

// List 获取分销员申请列表
func (c *ApplyController) List(ctx *gin.Context) {
	var req dto.DistributionApplyListReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(200, gin.H{"code": 400, "msg": err.Error()})
		return
	}

	list, total, err := c.applyService.GetApplyList(&req)
	if err != nil {
		ctx.JSON(200, gin.H{"code": 500, "msg": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"code": 200, "msg": "success", "data": gin.H{"list": list, "total": total}})
}

// Approve 审核通过
func (c *ApplyController) Approve(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.applyService.ApproveApply(id); err != nil {
		ctx.JSON(200, gin.H{"code": 500, "msg": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{"code": 200, "msg": "success"})
}

// Reject 审核拒绝
func (c *ApplyController) Reject(ctx *gin.Context) {
	var req dto.DistributionApplyRejectReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(200, gin.H{"code": 400, "msg": err.Error()})
		return
	}

	if err := c.applyService.RejectApply(&req); err != nil {
		ctx.JSON(200, gin.H{"code": 500, "msg": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{"code": 200, "msg": "success"})
}

// Detail 获取申请详情
func (c *ApplyController) Detail(ctx *gin.Context) {
	id := ctx.Param("id")
	detail, err := c.applyService.GetApplyDetail(id)
	if err != nil {
		ctx.JSON(200, gin.H{"code": 500, "msg": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{"code": 200, "msg": "success", "data": detail})
}
