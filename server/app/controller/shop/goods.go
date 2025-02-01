package shop

import (
	"fmt"
	"wnadmin/app/common/response"
	"wnadmin/app/common/utils"
	"wnadmin/app/model/dto"
	"wnadmin/app/service/shop"

	"github.com/gin-gonic/gin"
)

// GoodsController 商品控制器
type GoodsController struct {
	goodsService *shop.GoodsService
}

// NewGoodsController 创建商品控制器
func NewGoodsController() *GoodsController {
	return &GoodsController{
		goodsService: &shop.GoodsService{},
	}
}

// List 获取商品列表
// @Summary 获取商品列表
// @Description 获取商品列表，支持分页和条件查询
// @Tags 商品管理
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Param pageNum query int false "页码"
// @Param pageSize query int false "每页数量"
// @Param goodsName query string false "商品名称"
// @Param status query string false "状态"
// @Success 200 {object} response.Response "成功"
// @Router /shop/goods/list [get]
func (c *GoodsController) List(ctx *gin.Context) {
	var query dto.GoodsListReq
	// 打印原始请求参数
	fmt.Printf("原始请求参数: %+v\n", ctx.Request.URL.Query())

	// 使用 ShouldBindQuery 替代 ShouldBind
	if err := ctx.ShouldBindQuery(&query); err != nil {
		response.FailWithMessage(201, "参数错误: "+err.Error(), ctx)
		return
	}

	// 打印绑定后的参数
	fmt.Printf("绑定后的参数: %+v\n", query)

	// 设置默认值
	if query.PageNum <= 0 {
		query.PageNum = 1
	}
	if query.PageSize <= 0 {
		query.PageSize = 10
	}

	result, err := c.goodsService.GetGoodsList(&query)
	if err != nil {
		response.FailWithMessage(201, "获取商品列表失败: "+err.Error(), ctx)
		return
	}
	response.OkWithData(result, ctx)
}

// Delete 删除商品
// @Summary 删除商品
// @Description 删除商品
// @Tags 商品管理
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Param id path int true "商品ID"
// @Success 200 {object} response.Response "成功"
// @Router /shop/goods/delete/{id} [delete]
func (c *GoodsController) Delete(ctx *gin.Context) {
	id := utils.ParseParamID(ctx)
	if id <= 0 {
		response.FailWithMessage(201, "无效的商品ID", ctx)
		return
	}

	if err := c.goodsService.DeleteGoods(id); err != nil {
		response.FailWithMessage(201, "删除商品失败: "+err.Error(), ctx)
		return
	}

	response.OkWithMessage("删除成功", ctx)
}

// UpdateStatus 更新商品状态
// @Summary 更新商品状态
// @Description 更新商品状态
// @Tags 商品管理
// @Accept json
// @Produce json
// @Param data body dto.GoodsStatusReq true "状态信息"
// @Success 200 {object} response.Response "成功"
// @Router /shop/goods/status [put]
func (c *GoodsController) UpdateStatus(ctx *gin.Context) {
	var req dto.GoodsStatusReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(201, "参数错误: "+err.Error(), ctx)
		return
	}

	// 转换状态值为int8
	status := int8(0)
	if req.Status == "1" {
		status = 1
	}

	if err := c.goodsService.UpdateGoodsStatus(req.Id, status); err != nil {
		response.FailWithMessage(201, "更新商品状态失败: "+err.Error(), ctx)
		return
	}

	response.OkWithMessage("更新成功", ctx)
}
