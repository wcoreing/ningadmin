package system

import (
	"wnadmin/app/common/response"
	"wnadmin/app/common/utils"
	"wnadmin/app/model/do/SystemDB"
	"wnadmin/app/model/dto"

	"github.com/gin-gonic/gin"
)

type LoginController struct{}

// Login 用户登录
func (c *LoginController) Login(ctx *gin.Context) {
	var req dto.LoginReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(201, err.Error(), ctx)
		return
	}

	// 1. 根据用户名查询用户
	user := &SystemDB.UserDO{}
	err := user.Query().Where("user_name = ? AND del_flag = ?", req.UserName, "0").First(user).Error
	if err != nil {
		response.FailWithMessage(201, "用户名或密码错误", ctx)
		return
	}

	// 2. 检查用户状态
	if user.Status == "1" {
		response.FailWithMessage(201, "账号已被停用", ctx)
		return
	}

	// 3. 验证密码
	if !utils.ComparePassword(user.Password, req.Password) {
		response.FailWithMessage(201, "用户名或密码错误", ctx)
		return
	}

	// 4. 生成令牌
	token, err := utils.GenerateToken(user.ID, user.UserName)
	if err != nil {
		response.FailWithMessage(201, "生成令牌失败", ctx)
		return
	}

	// 5. 构建响应数据
	userInfo := dto.UserInfo{
		UserID:      user.ID,
		UserName:    user.UserName,
		NickName:    user.NickName,
		Email:       user.Email,
		Phonenumber: user.Phonenumber,
		Sex:         user.Sex,
		Avatar:      user.Avatar,
		Status:      user.Status,
		CreateTime:  user.CreateTime.Format("2006-01-02 15:04:05"),
	}

	resp := dto.LoginResp{
		Token:    token,
		UserInfo: userInfo,
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

	// 查询用户信息
	user := &SystemDB.UserDO{ID: userID.(int64)}
	if err := user.First(); err != nil {
		response.FailWithMessage(201, err.Error(), ctx)
		return
	}

	// 构建用户信息
	userInfo := dto.UserInfo{
		UserID:      user.ID,
		UserName:    user.UserName,
		NickName:    user.NickName,
		Email:       user.Email,
		Phonenumber: user.Phonenumber,
		Sex:         user.Sex,
		Avatar:      user.Avatar,
		Status:      user.Status,
		CreateTime:  user.CreateTime.Format("2006-01-02 15:04:05"),
	}

	response.OkWithDetailed(userInfo, "获取成功", ctx)
}

// Logout 用户登出
func (c *LoginController) Logout(ctx *gin.Context) {
	// 这里可以添加一些清理工作，比如清除缓存中的用户信息等
	response.OkWithMessage("退出成功", ctx)
}
