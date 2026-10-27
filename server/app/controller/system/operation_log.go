package system

import (
	"strconv"
	"wnadmin/app/common/response"
	"wnadmin/app/model/dto"
	"wnadmin/app/service/system"

	"github.com/gin-gonic/gin"
)

type OperationLogController struct {
	operationLogService *system.OperationLogService
}

func NewOperationLogController() *OperationLogController {
	return &OperationLogController{
		operationLogService: &system.OperationLogService{},
	}
}

func (c *OperationLogController) List(ctx *gin.Context) {
	var req dto.OperationLogListReq

	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(201, "参数错误", ctx)
		return
	}

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	result, err := c.operationLogService.GetLogList(&req)
	if err != nil {
		response.FailWithMessage(201, "查询失败", ctx)
		return
	}

	response.OkWithData(result, ctx)
}

func (c *OperationLogController) Delete(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.FailWithMessage(201, "参数错误", ctx)
		return
	}

	if err := c.operationLogService.DeleteLog(id); err != nil {
		response.FailWithMessage(201, "删除失败", ctx)
		return
	}

	response.Ok(ctx)
}

func (c *OperationLogController) BatchDelete(ctx *gin.Context) {
	var req dto.OperationLogBatchDeleteReq

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(201, "参数错误", ctx)
		return
	}

	if len(req.Ids) == 0 {
		response.FailWithMessage(201, "请选择要删除的记录", ctx)
		return
	}

	if err := c.operationLogService.BatchDeleteLogs(req.Ids); err != nil {
		response.FailWithMessage(201, "批量删除失败", ctx)
		return
	}

	response.Ok(ctx)
}

func (c *OperationLogController) Clear(ctx *gin.Context) {
	var req dto.OperationLogClearReq

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(201, "参数错误", ctx)
		return
	}

	if req.Days <= 0 {
		req.Days = 30
	}

	count, err := c.operationLogService.ClearLogs(req.Days)
	if err != nil {
		response.FailWithMessage(201, "清空日志失败", ctx)
		return
	}

	response.OkWithData(map[string]interface{}{
		"count": count,
	}, ctx)
}

func (c *OperationLogController) Export(ctx *gin.Context) {
	var req dto.OperationLogListReq

	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(201, "参数错误", ctx)
		return
	}

	// 导出功能暂未实现
	response.FailWithMessage(201, "导出功能开发中", ctx)
}
