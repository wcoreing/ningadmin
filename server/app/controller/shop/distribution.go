package shop

import (
	"wnadmin/app/service/shop"

	"github.com/gin-gonic/gin"
)

// DistributionController 分销控制器
type DistributionController struct {
	distributionService *shop.DistributionService
}

// NewDistributionController 创建分销控制器
func NewDistributionController() *DistributionController {
	return &DistributionController{
		distributionService: shop.NewDistributionService(),
	}
}

// 获取分销配置
func (c *DistributionController) GetConfig(ctx *gin.Context) {
	// ...
}

// 更新分销配置
func (c *DistributionController) UpdateConfig(ctx *gin.Context) {
	// ...
}

// 分销员列表
func (c *DistributionController) ListDistributors(ctx *gin.Context) {
	// ...
}

// 审核分销员
func (c *DistributionController) AuditDistributor(ctx *gin.Context) {
	// ...
}
