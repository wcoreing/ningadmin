package common

import (
	"wnadmin/app/common/response"
	"wnadmin/app/model/dto/common"
	"wnadmin/app/service/common/file"

	"github.com/gin-gonic/gin"
)

const (
	ERROR = response.ERROR // 使用 response 包中定义的错误码
)

type FileController struct {
	fileService *file.Service
}

func NewFileController() *FileController {
	return &FileController{
		fileService: file.NewService(),
	}
}

// Upload 上传文件
// @Summary 上传文件
// @Tags 文件管理
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "文件"
// @Param category_id formData int false "分类ID"
// @Success 200 {object} model.Response{data=model.FileVO}
// @Router /api/v1/files/upload [post]
func (c *FileController) Upload(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		response.FailWithMessage(ERROR, err.Error(), ctx)
		return
	}

	categoryId := ctx.PostForm("category_id")
	result, err := c.fileService.Upload(file, categoryId)
	if err != nil {
		response.FailWithMessage(ERROR, err.Error(), ctx)
		return
	}

	response.OkWithData(result, ctx)
}

// List 获取文件列表
// @Summary 获取文件列表
// @Tags 文件管理
// @Produce json
// @Param category_id query int false "分类ID"
// @Param keyword query string false "关键字"
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} model.Response{data=model.PageResult{list=[]model.FileVO}}
// @Router /api/v1/files [get]
func (c *FileController) List(ctx *gin.Context) {
	var query common.FileQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		response.FailWithMessage(ERROR, err.Error(), ctx)
		return
	}

	list, total, err := c.fileService.GetList(&query)
	if err != nil {
		response.FailWithMessage(ERROR, err.Error(), ctx)
		return
	}

	response.OkWithData(gin.H{
		"list":     list,
		"total":    total,
		"page":     query.Page,
		"pageSize": query.PageSize,
	}, ctx)
}

// Delete 删除文件
// @Summary 删除文件
// @Tags 文件管理
// @Produce json
// @Param id path int true "文件ID"
// @Success 200 {object} model.Response
// @Router /api/v1/files/{id} [delete]
func (c *FileController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	err := c.fileService.Delete(id)
	if err != nil {
		response.FailWithMessage(ERROR, err.Error(), ctx)
		return
	}

	response.Ok(ctx)
}

// CategoryTree 获取文件分类树
// @Summary 获取文件分类树
// @Tags 文件管理
// @Produce json
// @Success 200 {object} model.Response{data=[]model.FileCategoryTreeVO}
// @Router /api/v1/files/categories [get]
func (c *FileController) CategoryTree(ctx *gin.Context) {
	result, err := c.fileService.GetCategoryTree()
	if err != nil {
		response.FailWithMessage(ERROR, err.Error(), ctx)
		return
	}

	response.OkWithData(result, ctx)
}
