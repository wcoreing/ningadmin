package user

import (
	"strconv"
	"wnadmin/app/common/response"
	"wnadmin/app/model/dto"
	"wnadmin/app/service/user"

	"github.com/gin-gonic/gin"
)

// UserController 用户控制器
type UserController struct {
	userService *user.Service
}

// NewUserController 创建用户控制器
func NewUserController() *UserController {
	return &UserController{
		userService: user.NewService(),
	}
}

// List 获取用户列表
// @Summary 获取用户列表
// @Description 获取用户列表，支持分页和条件查询
// @Tags 用户管理
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Param pageNum query int false "页码"
// @Param pageSize query int false "每页数量"
// @Param keyword query string false "关键字"
// @Param status query int false "状态"
// @Param level query int false "会员等级"
// @Param gender query int false "性别"
// @Success 200 {object} response.Response "成功"
// @Router /user/list [get]
func (c *UserController) List(ctx *gin.Context) {
	var query dto.UserQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		response.FailWithMessage(201, "参数错误: "+err.Error(), ctx)
		return
	}

	result, err := c.userService.GetUserList(&query)
	if err != nil {
		response.FailWithMessage(201, "获取用户列表失败: "+err.Error(), ctx)
		return
	}

	response.OkWithDetailed(gin.H{
		"list":  result.List,
		"total": result.Total,
	}, "获取成功", ctx)
}

// Add 新增用户
// @Summary 新增用户
// @Description 创建用户
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param data body dto.CreateUserReq true "用户信息"
// @Success 200 {object} response.Response "成功"
// @Router /user [post]
func (c *UserController) Add(ctx *gin.Context) {
	var req dto.CreateUserReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(201, "参数错误: "+err.Error(), ctx)
		return
	}

	if err := c.userService.CreateUser(&req); err != nil {
		response.FailWithMessage(201, "创建用户失败: "+err.Error(), ctx)
		return
	}

	response.OkWithMessage("创建成功", ctx)
}

// Update 更新用户
// @Summary 更新用户
// @Description 更新用户信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param data body dto.UpdateUserReq true "用户信息"
// @Success 200 {object} response.Response "成功"
// @Router /user [put]
func (c *UserController) Update(ctx *gin.Context) {
	var req dto.UpdateUserReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(201, "参数错误: "+err.Error(), ctx)
		return
	}

	if err := c.userService.UpdateUser(&req); err != nil {
		response.FailWithMessage(201, "更新用户失败: "+err.Error(), ctx)
		return
	}

	response.OkWithMessage("更新成功", ctx)
}

// Delete 删除用户
// @Summary 删除用户
// @Description 删除用户
// @Tags 用户管理
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Param id path int true "用户ID"
// @Success 200 {object} response.Response "成功"
// @Router /user/{id} [delete]
func (c *UserController) Delete(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.FailWithMessage(201, "无效的ID", ctx)
		return
	}

	if err := c.userService.DeleteUser(id); err != nil {
		response.FailWithMessage(201, "删除用户失败: "+err.Error(), ctx)
		return
	}

	response.OkWithMessage("删除成功", ctx)
}

// GetInfo 获取用户详情
// @Summary 获取用户详情
// @Description 获取用户详细信息
// @Tags 用户管理
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Param id path int true "用户ID"
// @Success 200 {object} response.Response "成功"
// @Router /user/{id} [get]
func (c *UserController) GetInfo(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.FailWithMessage(201, "无效的ID", ctx)
		return
	}

	userInfo, err := c.userService.GetUserInfo(id)
	if err != nil {
		response.FailWithMessage(201, "获取用户详情失败: "+err.Error(), ctx)
		return
	}

	response.OkWithDetailed(userInfo, "获取成功", ctx)
}

// UpdateStatus 更新用户状态
// @Summary 更新用户状态
// @Description 更新用户状态
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param data body dto.UpdateUserStatusReq true "状态信息"
// @Success 200 {object} response.Response "成功"
// @Router /user/status [put]
func (c *UserController) UpdateStatus(ctx *gin.Context) {
	var req dto.UpdateUserStatusReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(201, "参数错误: "+err.Error(), ctx)
		return
	}

	// 验证状态字段是否传递
	if req.Status == nil {
		response.FailWithMessage(201, "状态字段不能为空", ctx)
		return
	}

	if err := c.userService.UpdateUserStatus(req.ID, *req.Status); err != nil {
		response.FailWithMessage(201, "更新用户状态失败: "+err.Error(), ctx)
		return
	}

	response.OkWithMessage("更新成功", ctx)
}

// ResetPassword 重置用户密码
// @Summary 重置用户密码
// @Description 重置用户密码
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param data body dto.UpdateUserPasswordReq true "密码信息"
// @Success 200 {object} response.Response "成功"
// @Router /user/password [put]
func (c *UserController) ResetPassword(ctx *gin.Context) {
	var req dto.UpdateUserPasswordReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(201, "参数错误: "+err.Error(), ctx)
		return
	}

	if err := c.userService.UpdateUserPassword(req.ID, req.Password); err != nil {
		response.FailWithMessage(201, "重置密码失败: "+err.Error(), ctx)
		return
	}

	response.OkWithMessage("重置成功", ctx)
}

// UpdateBalance 更新用户余额
// @Summary 更新用户余额
// @Description 更新用户余额
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param data body dto.UpdateUserBalanceReq true "余额信息"
// @Success 200 {object} response.Response "成功"
// @Router /user/balance [put]
func (c *UserController) UpdateBalance(ctx *gin.Context) {
	var req dto.UpdateUserBalanceReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(201, "参数错误: "+err.Error(), ctx)
		return
	}

	// 验证余额字段是否传递
	if req.Balance == nil {
		response.FailWithMessage(201, "余额字段不能为空", ctx)
		return
	}

	if err := c.userService.UpdateUserBalance(&req); err != nil {
		response.FailWithMessage(201, "更新余额失败: "+err.Error(), ctx)
		return
	}

	response.OkWithMessage("更新成功", ctx)
}

// UpdatePoints 更新用户积分
// @Summary 更新用户积分
// @Description 更新用户积分
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param data body dto.UpdateUserPointsReq true "积分信息"
// @Success 200 {object} response.Response "成功"
// @Router /user/points [put]
func (c *UserController) UpdatePoints(ctx *gin.Context) {
	var req dto.UpdateUserPointsReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(201, "参数错误: "+err.Error(), ctx)
		return
	}

	// 验证积分字段是否传递
	if req.Points == nil {
		response.FailWithMessage(201, "积分字段不能为空", ctx)
		return
	}

	if err := c.userService.UpdateUserPoints(&req); err != nil {
		response.FailWithMessage(201, "更新积分失败: "+err.Error(), ctx)
		return
	}

	response.OkWithMessage("更新成功", ctx)
}

// GetChangeLogList 获取用户变动记录列表
// @Summary 获取用户变动记录列表
// @Description 获取用户余额积分变动记录列表
// @Tags 用户管理
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Param pageNum query int false "页码"
// @Param pageSize query int false "每页数量"
// @Param userId query int false "用户ID"
// @Param changeType query int false "变动数据类型(1-余额 2-积分)"
// @Param actionType query int false "变动行为(1-收入/增加 2-支出/减少)"
// @Param sourceType query int false "来源类型"
// @Param startTime query string false "开始时间"
// @Param endTime query string false "结束时间"
// @Success 200 {object} response.Response "成功"
// @Router /user/change-log [get]
func (c *UserController) GetChangeLogList(ctx *gin.Context) {
	var query dto.UserChangeLogQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		response.FailWithMessage(201, "参数错误: "+err.Error(), ctx)
		return
	}

	// 设置默认分页参数
	if query.PageNum <= 0 {
		query.PageNum = 1
	}
	if query.PageSize <= 0 {
		query.PageSize = 20
	}

	result, err := c.userService.GetUserChangeLogList(&query)
	if err != nil {
		response.FailWithMessage(201, "获取记录列表失败: "+err.Error(), ctx)
		return
	}

	response.OkWithDetailed(gin.H{
		"list":  result.List,
		"total": result.Total,
	}, "获取成功", ctx)
}

// GetUserChangeLog 获取指定用户的变动记录
// @Summary 获取指定用户的变动记录
// @Description 获取指定用户的余额积分变动记录
// @Tags 用户管理
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Param id path int true "用户ID"
// @Param pageNum query int false "页码"
// @Param pageSize query int false "每页数量"
// @Param changeType query int false "变动数据类型(1-余额 2-积分)"
// @Success 200 {object} response.Response "成功"
// @Router /user/{id}/change-log [get]
func (c *UserController) GetUserChangeLog(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.FailWithMessage(201, "无效的ID", ctx)
		return
	}

	var query dto.UserChangeLogQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		response.FailWithMessage(201, "参数错误: "+err.Error(), ctx)
		return
	}

	// 设置用户ID
	query.UserID = &id

	// 设置默认分页参数
	if query.PageNum <= 0 {
		query.PageNum = 1
	}
	if query.PageSize <= 0 {
		query.PageSize = 20
	}

	result, err := c.userService.GetUserChangeLogList(&query)
	if err != nil {
		response.FailWithMessage(201, "获取记录列表失败: "+err.Error(), ctx)
		return
	}

	response.OkWithDetailed(gin.H{
		"list":  result.List,
		"total": result.Total,
	}, "获取成功", ctx)
}
