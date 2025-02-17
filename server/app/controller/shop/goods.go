package shop

import (
	"fmt"
	"strconv"
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

	if err := c.goodsService.UpdateGoodsStatus(req.GoodsId, int8(req.Status)); err != nil {
		response.FailWithMessage(201, "更新商品状态失败: "+err.Error(), ctx)
		return
	}

	response.OkWithMessage("更新成功", ctx)
}

// GetInfo 获取商品详情
// @Summary 获取商品详情
// @Description 获取商品详情信息
// @Tags 商品管理
// @Accept json
// @Produce json
// @Param id path int true "商品ID"
// @Success 200 {object} response.Response{data=dto.GoodsInfoResp} "成功"
// @Router /shop/goods/info/{id} [get]
func (c *GoodsController) GetInfo(ctx *gin.Context) {
	// 获取商品ID
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.FailWithMessage(201, "无效的商品ID", ctx)
		return
	}

	// 获取商品详情
	info, err := c.goodsService.GetGoodsInfo(id)
	if err != nil {
		response.FailWithMessage(201, err.Error(), ctx)
		return
	}

	response.OkWithData(info, ctx)
}

// Update 更新商品
// @Summary 更新商品
// @Description 更新商品信息
// @Tags 商品管理
// @Accept json
// @Produce json
// @Param data body dto.GoodsUpdateReq true "商品信息"
// @Success 200 {object} response.Response "成功"
// @Router /shop/goods/update [put]
func (c *GoodsController) Update(ctx *gin.Context) {
	var req dto.GoodsUpdateReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(201, "参数错误: "+err.Error(), ctx)
		return
	}
	fmt.Printf("更新请求: %+v\n", req)
	// 更新商品
	if err := c.goodsService.UpdateGoods(&req); err != nil {
		response.FailWithMessage(201, err.Error(), ctx)
		return
	}

	response.OkWithMessage("更新成功", ctx)
}

// Create 创建商品
// @Summary 创建商品
// @Description 创建商品信息
// @Tags 商品管理
// @Accept json
// @Produce json
// @Param data body dto.GoodsUpdateReq true "商品信息"
// @Success 200 {object} response.Response "成功"
// @Router /shop/goods/create [post]
func (c *GoodsController) Create(ctx *gin.Context) {
	var req dto.GoodsUpdateReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(201, "参数错误", ctx)
		return
	}

	if err := c.goodsService.CreateGoods(&req); err != nil {
		response.FailWithMessage(201, err.Error(), ctx)
		return
	}

	response.OkWithMessage("创建成功", ctx)
}

// CopyGoods 复制商品
// @Summary 复制商品
// @Description 复制商品
// @Tags 商品管理
// @Accept json
// @Produce json
// @Param id path int true "商品ID"
// @Success 200 {object} response.Response "成功"
// @Router /shop/goods/copy/{id} [post]
func (c *GoodsController) CopyGoods(ctx *gin.Context) {
	id := ctx.Param("id")
	goodsId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		response.FailWithMessage(400, "商品ID格式错误", ctx)
		return
	}

	// 获取复制数量
	countStr := ctx.Query("count")
	count, err := strconv.Atoi(countStr)
	if err != nil || count < 1 {
		count = 1 // 默认复制1份
	}

	if err := c.goodsService.CopyGoods(goodsId, count); err != nil {
		response.FailWithMessage(500, err.Error(), ctx)
		return
	}

	response.Ok(ctx)
}

// BatchMove 批量移动商品
// @Summary 批量移动商品
// @Description 批量移动商品到指定分类
// @Tags 商品管理
// @Accept json
// @Produce json
// @Param data body dto.GoodsBatchMoveReq true "移动信息"
// @Success 200 {object} response.Response "成功"
// @Router /shop/goods/move [put]
func (c *GoodsController) BatchMove(ctx *gin.Context) {
	var req dto.GoodsBatchMoveReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(400, "参数错误", ctx)
		return
	}

	if err := c.goodsService.BatchMoveGoods(&req); err != nil {
		response.FailWithMessage(500, err.Error(), ctx)
		return
	}

	response.Ok(ctx)
}

// BatchDelete 批量删除商品
// @Summary 批量删除商品
// @Description 批量删除商品
// @Tags 商品管理
// @Accept json
// @Produce json
// @Param data body dto.GoodsBatchDeleteReq true "商品ID列表"
// @Success 200 {object} response.Response "成功"
// @Router /shop/goods/batch [delete]
func (c *GoodsController) BatchDelete(ctx *gin.Context) {
	var req dto.GoodsBatchDeleteReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(400, "参数错误", ctx)
		return
	}

	if err := c.goodsService.BatchDeleteGoods(&req); err != nil {
		response.FailWithMessage(500, err.Error(), ctx)
		return
	}

	response.Ok(ctx)
}
