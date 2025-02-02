package shop

import (
	"wnadmin/app/common/response"
	"wnadmin/app/model/do/shop"

	"github.com/gin-gonic/gin"
)

// CategoryController 商品分类控制器
type CategoryController struct{}

// NewCategoryController 创建商品分类控制器
func NewCategoryController() *CategoryController {
	return &CategoryController{}
}

// List 获取分类列表（树形结构）
// @Summary 获取分类列表（树形结构）
// @Description 获取分类列表（树形结构）
// @Tags 商品分类
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Router /shop/category/list [get]
func (c *CategoryController) List(ctx *gin.Context) {
	var categories []shop.ShopCategory
	if err := (&shop.ShopCategory{}).Query().
		Where("status = ?", 1).
		Order("sort DESC, id ASC").
		Find(&categories).Error; err != nil {
		response.FailWithMessage(201, "获取分类列表失败", ctx)
		return
	}

	// 构建树形结构
	tree := shop.BuildCategoryTree(categories)
	response.OkWithData(tree, ctx)
}
