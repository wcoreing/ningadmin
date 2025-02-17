package shop

import (
	"wnadmin/app/common/response"
	"wnadmin/app/model/dto"
	"wnadmin/app/service/shop"
	"wnadmin/app/service/shop/logistics"

	"github.com/gin-gonic/gin"
)

type OrderController struct {
	orderService     shop.OrderService
	logisticsService *shop.LogisticsService
}

// NewOrderController 创建订单控制器
func NewOrderController() *OrderController {
	return &OrderController{
		orderService:     shop.OrderService{},
		logisticsService: shop.NewLogisticsService(logistics.ProviderAliyun),
	}
}

// GetOrderList 获取订单列表
func (c *OrderController) GetOrderList(ctx *gin.Context) {
	var req dto.OrderListReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(400, err.Error(), ctx)
		return
	}

	// 设置默认值
	if req.PageNum <= 0 {
		req.PageNum = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	result, err := c.orderService.GetOrderList(&req)
	if err != nil {
		response.FailWithMessage(500, err.Error(), ctx)
		return
	}

	response.OkWithData(result, ctx)
}

// ShipOrder 订单发货
func (c *OrderController) ShipOrder(ctx *gin.Context) {
	var req dto.OrderShipReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(400, err.Error(), ctx)
		return
	}

	if err := c.orderService.ShipOrder(&req); err != nil {
		response.FailWithMessage(500, err.Error(), ctx)
		return
	}

	response.Ok(ctx)
}

// GetLogisticsTrack 获取物流轨迹
func (c *OrderController) GetLogisticsTrack(ctx *gin.Context) {
	var req dto.LogisticsTrackReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(400, err.Error(), ctx)
		return
	}

	result, err := c.logisticsService.GetLogisticsTrack(&req)
	if err != nil {
		response.FailWithMessage(500, err.Error(), ctx)
		return
	}

	response.OkWithData(result, ctx)
}
