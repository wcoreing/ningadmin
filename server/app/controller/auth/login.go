package auth

import (
	"github.com/gin-gonic/gin"
	"wnadmin/app/common/response"
	"wnadmin/app/model/dto"
	"wnadmin/app/service/auth"
)

type LoginController struct {
	authService *auth.AuthService
}

// NewLoginController 创建登录控制器
func NewLoginController() *LoginController {
	return &LoginController{
		authService: &auth.AuthService{},
	}
}

// Login 用户登录
func (c *LoginController) Login(ctx *gin.Context) {
	var req dto.LoginReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(201, err.Error(), ctx)
		return
	}

	// 调用服务层处理登录逻辑
	resp, err := c.authService.Login(&req)
	if err != nil {
		response.FailWithMessage(201, err.Error(), ctx)
		return
	}

	response.OkWithDetailed(resp, "登录成功", ctx)
}

// GetUserInfo 获取当前登录用户信息
func (c *LoginController) GetUserInfo(ctx *gin.Context) {
	// 从上下文中获取用户ID（需要配合JWT中间件使用）
	userID, exists := ctx.Get("userId")
	if !exists {
		response.FailWithMessage(201, "获取用户信息失败", ctx)
		return
	}

	// 调用服务层获取用户信息
	userInfo, err := c.authService.GetUserInfo(userID.(int64))
	if err != nil {
		response.FailWithMessage(201, err.Error(), ctx)
		return
	}

	response.OkWithDetailed(userInfo, "获取成功", ctx)
}

// Logout 用户登出
func (c *LoginController) Logout(ctx *gin.Context) {
	// 调用服务层处理登出逻辑
	if err := c.authService.Logout(); err != nil {
		response.FailWithMessage(201, err.Error(), ctx)
		return
	}

	response.OkWithMessage("退出成功", ctx)
}
