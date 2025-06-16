package common

import (
	"strconv"
	"wnadmin/app/common/response"
	"wnadmin/app/model/dto/common"
	commonService "wnadmin/app/service/common"

	"github.com/gin-gonic/gin"
)

// CategoryController 通用分类控制器
type CategoryController struct {
	categoryService *commonService.CategoryService
}

// NewCategoryController 创建分类控制器
func NewCategoryController() *CategoryController {
	return &CategoryController{
		categoryService: &commonService.CategoryService{},
	}
}

// GetCategoryList 获取分类列表
func (c *CategoryController) GetCategoryList(ctx *gin.Context) {
	module := ctx.Query("module")
	if module == "" {
		module = "goods" // 默认为商品分类
	}

	categories, err := c.categoryService.GetCategoryOptions(module)
	if err != nil {
		response.FailWithMessage(response.ERROR, "获取分类列表失败: "+err.Error(), ctx)
		return
	}

	response.OkWithData(categories, ctx)
}

// GetCategoryTree 获取分类树
func (c *CategoryController) GetCategoryTree(ctx *gin.Context) {
	module := ctx.Query("module")
	if module == "" {
		module = "goods" // 默认为商品分类
	}

	categories, err := c.categoryService.GetCategoryTree(module)
	if err != nil {
		response.FailWithMessage(response.ERROR, "获取分类树失败: "+err.Error(), ctx)
		return
	}

	response.OkWithData(categories, ctx)
}

// CreateCategory 创建分类
func (c *CategoryController) CreateCategory(ctx *gin.Context) {
	var req common.CategoryCreateReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(response.ERROR, "参数错误: "+err.Error(), ctx)
		return
	}

	err := c.categoryService.CreateCategory(&req)
	if err != nil {
		response.FailWithMessage(response.ERROR, "创建分类失败: "+err.Error(), ctx)
		return
	}

	response.OkWithMessage("创建成功", ctx)
}

// UpdateCategory 更新分类
func (c *CategoryController) UpdateCategory(ctx *gin.Context) {
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

	var req common.CategoryUpdateReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(response.ERROR, "参数错误: "+err.Error(), ctx)
		return
	}

	req.Id = idInt
	err = c.categoryService.UpdateCategory(&req)
	if err != nil {
		response.FailWithMessage(response.ERROR, "更新分类失败: "+err.Error(), ctx)
		return
	}

	response.OkWithMessage("更新成功", ctx)
}

// DeleteCategory 删除分类
func (c *CategoryController) DeleteCategory(ctx *gin.Context) {
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

	err = c.categoryService.DeleteCategory(idInt)
	if err != nil {
		response.FailWithMessage(response.ERROR, "删除分类失败: "+err.Error(), ctx)
		return
	}

	response.OkWithMessage("删除成功", ctx)
}
