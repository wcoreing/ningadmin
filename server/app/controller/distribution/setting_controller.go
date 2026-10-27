package distribution

import (
	"wnadmin/app/common/response"
	"wnadmin/app/model/dto"
	"wnadmin/app/service/distribution"

	"github.com/gin-gonic/gin"
)

type SettingController struct {
	settingService *distribution.SettingService
}

// NewSettingController 创建分销设置控制器
func NewSettingController() *SettingController {
	return &SettingController{
		settingService: distribution.NewSettingService(),
	}
}

// GetSetting 获取分销设置
func (c *SettingController) GetSetting(ctx *gin.Context) {
	setting, err := c.settingService.GetSetting(ctx)
	if err != nil {
		response.FailWithMessage(500, err.Error(), ctx)
		return
	}
	response.OkWithData(setting, ctx)
}

// UpdateSetting 更新分销设置
func (c *SettingController) UpdateSetting(ctx *gin.Context) {
	var req dto.DistributionSettingDTO
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(400, err.Error(), ctx)
		return
	}

	if err := c.settingService.UpdateSetting(ctx, &req); err != nil {
		response.FailWithMessage(500, err.Error(), ctx)
		return
	}

	response.Ok(ctx)
}
