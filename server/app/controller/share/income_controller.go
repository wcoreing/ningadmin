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
