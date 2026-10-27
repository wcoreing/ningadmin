package shop

import (
	"wnadmin/app/common/response"
	"wnadmin/app/service/shop"

	"github.com/gin-gonic/gin"
)

// TagController 商品标签控制器
type TagController struct {
	tagService *shop.TagService
}

// NewTagController 创建商品标签控制器
func NewTagController() *TagController {
	return &TagController{
		tagService: &shop.TagService{},
	}
}

// List 获取商品标签列表
// @Summary 获取商品标签列表
// @Description 获取商品标签列表
// @Tags 商品标签
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Router /shop/tag/list [get]
func (c *TagController) List(ctx *gin.Context) {
	list, err := c.tagService.GetTagList()
	if err != nil {
		response.FailWithMessage(201, "获取商品标签列表失败: "+err.Error(), ctx)
		return
	}
	response.OkWithData(list, ctx)
}
