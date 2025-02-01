package shop

import (
	"wnadmin/app/common/response"
	"wnadmin/app/service/shop"

	"github.com/gin-gonic/gin"
)

// SupplierController 供应商控制器
type SupplierController struct {
	supplierService *shop.SupplierService
}

// NewSupplierController 创建供应商控制器
func NewSupplierController() *SupplierController {
	return &SupplierController{
		supplierService: &shop.SupplierService{},
	}
}

// List 获取供应商列表
// @Summary 获取供应商列表
// @Description 获取供应商列表
// @Tags 商品供应商
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Router /shop/supplier/list [get]
func (c *SupplierController) List(ctx *gin.Context) {
	list, err := c.supplierService.GetSupplierList()
	if err != nil {
		response.FailWithMessage(201, "获取供应商列表失败: "+err.Error(), ctx)
		return
	}
	response.OkWithData(list, ctx)
}
