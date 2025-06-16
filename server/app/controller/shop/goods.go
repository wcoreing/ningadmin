package shop

import (
	"fmt"
	"strconv"
	"wnadmin/app/common/response"
	"wnadmin/app/common/utils"
	"wnadmin/app/model/do/shop"
	"wnadmin/app/model/dto"
	"wnadmin/app/service/shop/goods"

	"github.com/gin-gonic/gin"
)

// GoodsController 商品控制器
type GoodsController struct{}

// NewGoodsController 创建商品控制器
func NewGoodsController() *GoodsController {
	return &GoodsController{}
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
// @Param type query int false "商品类型：1-普通商品 2-虚拟商品 3-套餐商品"
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

	// 直接调用商品列表查询服务
	service := &goods.ListService{}
	result, err := service.Execute(&query)
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

	// 直接调用商品删除服务
	service := &goods.DeleteService{}
	if err := service.Execute(id); err != nil {
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

	// 直接调用商品状态更新服务
	service := &goods.UpdateStatusService{}
	if err := service.Execute(req.GoodsId, int8(req.Status)); err != nil {
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

	// 直接调用商品详情查询服务
	service := &goods.GetInfoService{}
	info, err := service.Execute(id)
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

	// 直接调用商品更新服务
	service := &goods.UpdateService{}
	if err := service.Execute(&req); err != nil {
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

	// 直接调用商品创建服务
	service := &goods.CreateService{}
	if err := service.Execute(&req); err != nil {
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

	// 直接调用商品复制服务
	service := &goods.CopyService{}
	if err := service.Execute(goodsId, count); err != nil {
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

	// 直接调用批量移动服务
	service := &goods.BatchMoveService{}
	if err := service.Execute(&req); err != nil {
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

	// 直接调用批量删除服务
	service := &goods.BatchDeleteService{}
	if err := service.Execute(&req); err != nil {
		response.FailWithMessage(500, err.Error(), ctx)
		return
	}

	response.Ok(ctx)
}

// GetByID 根据ID获取商品（兼容方法，供其他服务调用）
func (c *GoodsController) GetByID(id int64) (*shop.GoodsDO, error) {
	service := &goods.GetInfoService{}
	return service.GetGoodsById(id)
}

// ========== 套餐相关接口 ==========

// GetComboConfig 获取套餐配置
// @Summary 获取套餐配置
// @Description 获取套餐的组和选项配置信息
// @Tags 套餐管理
// @Accept json
// @Produce json
// @Param id path int true "套餐ID"
// @Success 200 {object} response.Response "成功"
// @Router /shop/goods/combo/config/{id} [get]
func (c *GoodsController) GetComboConfig(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.FailWithMessage(201, "无效的套餐ID", ctx)
		return
	}

	// 调用套餐配置服务
	service := &goods.GetComboConfigService{}
	config, err := service.Execute(id)
	if err != nil {
		response.FailWithMessage(201, err.Error(), ctx)
		return
	}

	response.OkWithData(config, ctx)
}

// CalculateComboPrice 计算套餐价格
// @Summary 计算套餐价格
// @Description 根据用户选择计算套餐总价格
// @Tags 套餐管理
// @Accept json
// @Produce json
// @Param data body ComboCalculateReq true "套餐选择信息"
// @Success 200 {object} response.Response "成功"
// @Router /shop/goods/combo/calculate [post]
func (c *GoodsController) CalculateComboPrice(ctx *gin.Context) {
	var req ComboCalculateReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(201, "参数错误: "+err.Error(), ctx)
		return
	}

	// 调用套餐价格计算服务
	service := &goods.CalculateComboService{}
	totalPrice, err := service.Execute(req.ComboId, req.Selections)
	if err != nil {
		response.FailWithMessage(201, err.Error(), ctx)
		return
	}

	response.OkWithData(gin.H{
		"totalPrice": totalPrice,
	}, ctx)
}

// ValidateComboSelection 验证套餐选择
// @Summary 验证套餐选择
// @Description 验证用户的套餐选择是否符合规则
// @Tags 套餐管理
// @Accept json
// @Produce json
// @Param data body ComboValidateReq true "套餐选择信息"
// @Success 200 {object} response.Response "成功"
// @Router /shop/goods/combo/validate [post]
func (c *GoodsController) ValidateComboSelection(ctx *gin.Context) {
	var req ComboValidateReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(201, "参数错误: "+err.Error(), ctx)
		return
	}

	// 调用套餐选择验证服务
	service := &goods.ValidateComboService{}
	err := service.Execute(req.ComboId, req.Selections)
	if err != nil {
		response.FailWithMessage(201, err.Error(), ctx)
		return
	}

	response.OkWithMessage("验证通过", ctx)
}

// GetDefaultComboSelection 获取默认套餐选择
// @Summary 获取默认套餐选择
// @Description 获取套餐的默认选择配置
// @Tags 套餐管理
// @Accept json
// @Produce json
// @Param id path int true "套餐ID"
// @Success 200 {object} response.Response "成功"
// @Router /shop/goods/combo/default/{id} [get]
func (c *GoodsController) GetDefaultComboSelection(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.FailWithMessage(201, "无效的套餐ID", ctx)
		return
	}

	// 调用默认套餐选择服务
	service := &goods.GetDefaultComboSelectionService{}
	selections, err := service.Execute(id)
	if err != nil {
		response.FailWithMessage(201, err.Error(), ctx)
		return
	}

	response.OkWithData(selections, ctx)
}

// ComboCalculateReq 套餐价格计算请求
type ComboCalculateReq struct {
	ComboId    int64                   `json:"comboId" binding:"required"`
	Selections []*goods.ComboSelection `json:"selections" binding:"required"`
}

// ComboValidateReq 套餐选择验证请求
type ComboValidateReq struct {
	ComboId    int64                   `json:"comboId" binding:"required"`
	Selections []*goods.ComboSelection `json:"selections" binding:"required"`
}
