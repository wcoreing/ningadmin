package system

import (
	"strconv"
	"wnadmin/app/common/response"
	"wnadmin/app/model/dto"
	"wnadmin/app/service/system"

	"github.com/gin-gonic/gin"
)

// UserController 用户控制器
type UserController struct {
	userService *system.UserService
}

// NewUserController 创建用户控制器
func NewUserController() *UserController {
	return &UserController{
		userService: &system.UserService{},
	}
}

// List 获取用户列表
// @Summary 获取用户列表
// @Description 获取用户列表，支持分页和条件查询
// @Tags 系统用户
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Param pageNum query int false "页码"
// @Param pageSize query int false "每页数量"
// @Param keyword query string false "关键字"
// @Param status query string false "状态"
// @Success 200 {object} response.Response "成功"
// @Router /system/user/list [get]
func (c *UserController) List(ctx *gin.Context) {
	var query dto.SystemUserQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		response.FailWithMessage(201, "参数错误: "+err.Error(), ctx)
		return
	}

	// 设置默认值
	if query.PageNum <= 0 {
		query.PageNum = 1
	}
	if query.PageSize <= 0 {
		query.PageSize = 10
	}

	result, err := c.userService.GetUserList(&query)
	if err != nil {
		response.FailWithMessage(201, "获取用户列表失败: "+err.Error(), ctx)
		return
	}

	response.OkWithDetailed(gin.H{
		"total": result.Total,
		"list":  result.Rows,
	}, "获取成功", ctx)
}

// Add 新增用户
// @Summary 新增用户
// @Description 新增系统用户
// @Tags 系统用户
// @Accept json
// @Produce json
// @Param data body dto.CreateUserReq true "用户信息"
// @Success 200 {object} response.Response "成功"
// @Router /system/user [post]
func (c *UserController) Add(ctx *gin.Context) {
	var req dto.CreateSystemUserReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(201, "参数错误: "+err.Error(), ctx)
		return
	}

	// 从上下文获取创建人
	createBy := ctx.GetString("userName")
	if err := c.userService.CreateUser(&req, createBy); err != nil {
		response.FailWithMessage(201, "创建用户失败: "+err.Error(), ctx)
		return
	}

	response.OkWithMessage("创建成功", ctx)
}

// Update 更新用户
// @Summary 更新用户
// @Description 更新系统用户
// @Tags 系统用户
// @Accept json
// @Produce json
// @Param data body dto.UpdateUserReq true "用户信息"
// @Success 200 {object} response.Response "成功"
// @Router /system/user [put]
func (c *UserController) Update(ctx *gin.Context) {
	var req dto.UpdateSystemUserReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(201, "参数错误: "+err.Error(), ctx)
		return
	}

	// 从上下文获取更新人
	updateBy := ctx.GetString("userName")
	if err := c.userService.UpdateUser(&req, updateBy); err != nil {
		response.FailWithMessage(201, "更新用户失败: "+err.Error(), ctx)
		return
	}

	response.OkWithMessage("更新成功", ctx)
}

// Delete 删除用户
// @Summary 删除用户
// @Description 删除系统用户
// @Tags 系统用户
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Param userId path int true "用户ID"
// @Success 200 {object} response.Response "成功"
// @Router /system/user/{userId} [delete]
func (c *UserController) Delete(ctx *gin.Context) {
	idStr := ctx.Param("userId")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.FailWithMessage(201, "无效的用户ID", ctx)
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
// @Description 获取系统用户详情
// @Tags 系统用户
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Param userId path int true "用户ID"
// @Success 200 {object} response.Response "成功"
// @Router /system/user/{userId} [get]
func (c *UserController) GetInfo(ctx *gin.Context) {
	idStr := ctx.Param("userId")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.FailWithMessage(201, "无效的用户ID", ctx)
		return
	}

	user, err := c.userService.GetUserInfo(id)
	if err != nil {
		response.FailWithMessage(201, "获取用户详情失败: "+err.Error(), ctx)
		return
	}

	response.OkWithDetailed(user, "获取成功", ctx)
}

// UpdateStatus 更新用户状态
// @Summary 更新用户状态
// @Description 更新系统用户状态
// @Tags 系统用户
// @Accept json
// @Produce json
// @Param data body dto.UpdateUserStatusReq true "状态信息"
// @Success 200 {object} response.Response "成功"
// @Router /system/user/status [put]
func (c *UserController) UpdateStatus(ctx *gin.Context) {
	var req dto.UpdateUserStatusReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(201, "参数错误: "+err.Error(), ctx)
		return
	}

	// 转换状态类型：int -> string
	status := "0"
	if *req.Status == 1 {
		status = "1"
	}

	if err := c.userService.UpdateUserStatus(req.ID, status); err != nil {
		response.FailWithMessage(201, "更新用户状态失败: "+err.Error(), ctx)
		return
	}

	response.OkWithMessage("更新成功", ctx)
}

// ResetPassword 重置用户密码
// @Summary 重置用户密码
// @Description 重置系统用户密码
// @Tags 系统用户
// @Accept json
// @Produce json
// @Param data body dto.UpdateUserPasswordReq true "密码信息"
// @Success 200 {object} response.Response "成功"
// @Router /system/user/password [put]
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
