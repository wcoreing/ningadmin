package shop

import (
	"wnadmin/app/common/response"
	"wnadmin/app/service/shop"

	"github.com/gin-gonic/gin"
)

// CategoryController 商品分类控制器
type CategoryController struct {
	categoryService *shop.CategoryService
}

// NewCategoryController 创建商品分类控制器
func NewCategoryController() *CategoryController {
	return &CategoryController{
		categoryService: &shop.CategoryService{},
	}
}

// List 获取商品分类列表
// @Summary 获取商品分类列表
// @Description 获取商品分类列表
// @Tags 商品分类
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Router /shop/category/list [get]
func (c *CategoryController) List(ctx *gin.Context) {
	list, err := c.categoryService.GetCategoryList()
	if err != nil {
		response.FailWithMessage(201, "获取商品分类列表失败: "+err.Error(), ctx)
		return
	}
	response.OkWithData(list, ctx)
}
