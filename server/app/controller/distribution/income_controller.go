package distribution

import (
	"wnadmin/app/common/response"
	"wnadmin/app/model/dto"
	"wnadmin/app/service/distribution"

	"github.com/gin-gonic/gin"
)

type IncomeController struct {
	incomeService *distribution.IncomeService
}

func NewIncomeController() *IncomeController {
	return &IncomeController{
		incomeService: distribution.NewIncomeService(),
	}
}

// GetPage 分页获取分销收益列表
func (c *IncomeController) GetIncomList(ctx *gin.Context) {
	var req dto.DistributionIncomeListReq
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

// Unfreeze 解冻分销收益
func (c *IncomeController) Unfreeze(ctx *gin.Context) {
	var req dto.DistributionIncomeUnfreezeReq
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
