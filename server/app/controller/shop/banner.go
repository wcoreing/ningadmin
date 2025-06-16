package shop

import (
	"fmt"
	"strconv"
	"wnadmin/app/common/response"
	"wnadmin/app/model/dto"
	shop_service "wnadmin/app/service/shop"

	"github.com/gin-gonic/gin"
)

type BannerController struct {
	bannerService shop_service.BannerService
}

// NewBannerController 创建Banner控制器
func NewBannerController() *BannerController {
	return &BannerController{
		bannerService: shop_service.BannerService{},
	}
}

// 获取Banner列表
func (c *BannerController) GetBannerList(ctx *gin.Context) {
	var req dto.BannerListReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		fmt.Printf("Bind error: %v\n", err)
		response.FailWithMessage(400, err.Error(), ctx)
		return
	}

	fmt.Printf("Query params: %+v\n", req)

	result, err := c.bannerService.GetBannerList(&req)
	if err != nil {
		fmt.Printf("Query error: %v\n", err)
		response.FailWithMessage(500, err.Error(), ctx)
		return
	}

	fmt.Printf("Query result: %+v\n", result)

	response.OkWithData(gin.H{
		"list":  result.List,
		"total": result.Total,
	}, ctx)
}

// 获取指定位置的Banner列表
func (c *BannerController) GetBannersByPosition(ctx *gin.Context) {
	position := ctx.Query("position")
	if position == "" {
		response.FailWithMessage(400, "position不能为空", ctx)
		return
	}

	// 转换为int
	positionInt, err := strconv.Atoi(position)
	if err != nil {
		response.FailWithMessage(400, "position参数错误", ctx)
		return
	}

	list, err := c.bannerService.GetBannersByPosition(positionInt)
	if err != nil {
		response.FailWithMessage(500, err.Error(), ctx)
		return
	}

	response.OkWithData(list, ctx)
}

// 创建Banner
func (c *BannerController) CreateBanner(ctx *gin.Context) {
	var req dto.BannerCreateReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		fmt.Printf("Bind error: %v\n", err)
		response.FailWithMessage(400, "参数错误: "+err.Error(), ctx)
		return
	}

	fmt.Printf("Create banner request: %+v\n", req)
	err := c.bannerService.CreateBanner(&req)
	if err != nil {
		fmt.Printf("Create banner error: %v\n", err)
		response.FailWithMessage(500, "创建失败: "+err.Error(), ctx)
		return
	}

	response.Ok(ctx)
}

// 更新Banner
func (c *BannerController) UpdateBanner(ctx *gin.Context) {
	var req dto.BannerUpdateReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		fmt.Printf("Banner更新参数绑定失败: %v\n", err)
		response.FailWithMessage(400, "参数错误: "+err.Error(), ctx)
		return
	}

	fmt.Printf("Banner更新请求参数: %+v\n", req)
	err := c.bannerService.UpdateBanner(&req)
	if err != nil {
		fmt.Printf("Banner更新失败: %v\n", err)
		response.FailWithMessage(500, err.Error(), ctx)
		return
	}

	response.Ok(ctx)
}

// 删除Banner
func (c *BannerController) DeleteBanner(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		response.FailWithMessage(400, "id不能为空", ctx)
		return
	}

	// 转换id参数
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		response.FailWithMessage(400, "id参数错误", ctx)
		return
	}

	err = c.bannerService.DeleteBanner(idInt)
	if err != nil {
		response.FailWithMessage(500, err.Error(), ctx)
		return
	}

	response.Ok(ctx)
}

// 修改Banner状态
func (c *BannerController) ChangeBannerStatus(ctx *gin.Context) {
	id := ctx.Param("id")
	status := ctx.Query("status")
	if id == "" || status == "" {
		response.FailWithMessage(400, "参数错误", ctx)
		return
	}

	// 转换参数
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		response.FailWithMessage(400, "id参数错误", ctx)
		return
	}
	statusInt, err := strconv.Atoi(status)
	if err != nil {
		response.FailWithMessage(400, "status参数错误", ctx)
		return
	}

	err = c.bannerService.ChangeBannerStatus(idInt, statusInt)
	if err != nil {
		response.FailWithMessage(500, err.Error(), ctx)
		return
	}

	response.Ok(ctx)
}

// GetBannerInfo 获取Banner详情
// @Summary 获取Banner详情
// @Description 获取Banner详情
// @Tags Banner管理
// @Accept json
// @Produce json
// @Param id path int true "Banner ID"
// @Success 200 {object} response.Response{data=dto.BannerInfoResp} "成功"
// @Router /shop/banner/info/{id} [get]
func (c *BannerController) GetBannerInfo(ctx *gin.Context) {
	id := ctx.Param("id")
	bannerId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		response.FailWithMessage(400, "Banner ID格式错误", ctx)
		return
	}

	info, err := c.bannerService.GetBannerInfo(bannerId)
	if err != nil {
		response.FailWithMessage(500, err.Error(), ctx)
		return
	}

	response.OkWithData(info, ctx)
}
