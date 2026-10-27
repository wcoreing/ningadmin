package shop

import (
	"strconv"
	"wnadmin/app/common/response"
	"wnadmin/app/model/dto/shop"
	shop_service "wnadmin/app/service/shop"

	"github.com/gin-gonic/gin"
)

// CategoryController 商品分类控制器
type CategoryController struct {
	categoryService *shop_service.CategoryService
}

// NewCategoryController 创建商品分类控制器
func NewCategoryController() *CategoryController {
	return &CategoryController{
		categoryService: shop_service.NewCategoryService(),
	}
}

// List 获取分类列表
// @Summary 获取分类列表
// @Tags 商品分类
// @Accept json
// @Produce json
// @Param name query string false "分类名称"
// @Param status query int false "状态"
// @Param parent_id query int false "父级ID"
// @Success 200 {object} response.Response{data=[]shop.CategoryVO}
// @Router /shop/categories [get]
func (c *CategoryController) List(ctx *gin.Context) {
	var req shop.CategoryListReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(response.ERROR, err.Error(), ctx)
		return
	}

	list, err := c.categoryService.GetList(&req)
	if err != nil {
		response.FailWithMessage(response.ERROR, err.Error(), ctx)
		return
	}

	response.OkWithData(list, ctx)
}

// Create 创建分类
// @Summary 创建分类
// @Tags 商品分类
// @Accept json
// @Produce json
// @Param data body shop.CategoryCreateReq true "创建分类请求参数"
// @Success 200 {object} response.Response
// @Router /shop/categories [post]
func (c *CategoryController) Create(ctx *gin.Context) {
	var req shop.CategoryCreateReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(response.ERROR, err.Error(), ctx)
		return
	}

	err := c.categoryService.Create(&req)
	if err != nil {
		response.FailWithMessage(response.ERROR, err.Error(), ctx)
		return
	}

	response.Ok(ctx)
}

// Update 更新分类
// @Summary 更新分类
// @Tags 商品分类
// @Accept json
// @Produce json
// @Param id path int true "分类ID"
// @Param data body shop.CategoryUpdateReq true "更新分类请求参数"
// @Success 200 {object} response.Response
// @Router /shop/categories/{id} [put]
func (c *CategoryController) Update(ctx *gin.Context) {
	var req shop.CategoryUpdateReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(response.ERROR, err.Error(), ctx)
		return
	}

	err := c.categoryService.Update(&req)
	if err != nil {
		response.FailWithMessage(response.ERROR, err.Error(), ctx)
		return
	}

	response.Ok(ctx)
}

// Delete 删除分类
// @Summary 删除分类
// @Tags 商品分类
// @Accept json
// @Produce json
// @Param id path int true "分类ID"
// @Success 200 {object} response.Response
// @Router /shop/categories/{id} [delete]
func (c *CategoryController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		response.FailWithMessage(response.ERROR, "ID不能为空", ctx)
		return
	}

	// 转换ID为int64
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		response.FailWithMessage(response.ERROR, "ID格式错误", ctx)
		return
	}

	err = c.categoryService.Delete(idInt)
	if err != nil {
		response.FailWithMessage(response.ERROR, err.Error(), ctx)
		return
	}

	response.Ok(ctx)
}
