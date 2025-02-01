package shop

import (
	"wnadmin/app/common/response"
	"wnadmin/app/service/shop"

	"github.com/gin-gonic/gin"
)

// BrandController 品牌控制器
type BrandController struct {
	brandService *shop.BrandService
}

// NewBrandController 创建品牌控制器
func NewBrandController() *BrandController {
	return &BrandController{
		brandService: &shop.BrandService{},
	}
}

// List 获取品牌列表
// @Summary 获取品牌列表
// @Description 获取品牌列表
// @Tags 商品品牌
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Router /shop/brand/list [get]
func (c *BrandController) List(ctx *gin.Context) {
	list, err := c.brandService.GetBrandList()
	if err != nil {
		response.FailWithMessage(201, "获取品牌列表失败: "+err.Error(), ctx)
		return
	}
	response.OkWithData(list, ctx)
}
