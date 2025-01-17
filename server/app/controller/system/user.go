package system

import (
	"strconv"
	"wnadmin/app/common/response"
	"wnadmin/app/common/utils"
	"wnadmin/app/model/do/SystemDB"

	"github.com/gin-gonic/gin"
)

type UserController struct{}

// List 获取用户列表
func (c *UserController) List(ctx *gin.Context) {
	var query struct {
		PageNum  int    `form:"pageNum"`
		PageSize int    `form:"pageSize"`
		UserName string `form:"userName"`
		Status   string `form:"status"`
	}

	if err := ctx.ShouldBindQuery(&query); err != nil {
		response.FailWithMessage(201, err.Error(), ctx)
		return
	}

	if query.PageNum < 1 {
		query.PageNum = 1
	}
	if query.PageSize < 1 {
		query.PageSize = 10
	}

	userDO := &SystemDB.UserDO{}
	db := userDO.Query().Where("del_flag = ?", "0")

	if query.UserName != "" {
		db = db.Where("user_name LIKE ?", "%"+query.UserName+"%")
	}
	if query.Status != "" {
		db = db.Where("status = ?", query.Status)
	}

	var total int64
	err := db.Count(&total).Error
	if err != nil {
		response.FailWithMessage(201, err.Error(), ctx)
		return
	}

	var users []SystemDB.UserDO
	err = db.Offset((query.PageNum - 1) * query.PageSize).
		Limit(query.PageSize).
		Order("user_id desc").
		Find(&users).Error
	if err != nil {
		response.FailWithMessage(201, err.Error(), ctx)
		return
	}

	response.OkWithDetailed(gin.H{
		"total": total,
		"list":  users,
	}, "获取成功", ctx)
}

// GetInfo 获取用户详情
func (c *UserController) GetInfo(ctx *gin.Context) {
	idStr := ctx.Param("userId")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.FailWithMessage(201, "无效的用户ID", ctx)
		return
	}

	user := &SystemDB.UserDO{ID: id}
	if err := user.First(); err != nil {
		response.FailWithMessage(201, err.Error(), ctx)
		return
	}

	response.OkWithDetailed(user, "获取成功", ctx)
}

// Add 新增用户
func (c *UserController) Add(ctx *gin.Context) {
	var user SystemDB.UserDO
	if err := ctx.ShouldBindJSON(&user); err != nil {
		response.FailWithMessage(201, err.Error(), ctx)
		return
	}

	// 检查用户名是否存在
	exists, err := user.CheckUserNameExists(user.UserName, 0)
	if err != nil {
		response.FailWithMessage(201, err.Error(), ctx)
		return
	}
	if exists {
		response.FailWithMessage(201, "用户名已存在", ctx)
		return
	}

	// 加密密码
	encryptedPassword, err := utils.EncryptPassword(user.Password)
	if err != nil {
		response.FailWithMessage(201, "密码加密失败", ctx)
		return
	}
	user.Password = encryptedPassword

	user.DelFlag = "0"
	if err := user.Create(); err != nil {
		response.FailWithMessage(201, err.Error(), ctx)
		return
	}

	response.OkWithMessage("创建成功", ctx)
}

// Update 更新用户
func (c *UserController) Update(ctx *gin.Context) {
	var user SystemDB.UserDO
	if err := ctx.ShouldBindJSON(&user); err != nil {
		response.FailWithMessage(201, err.Error(), ctx)
		return
	}

	// 检查用户名是否存在
	exists, err := user.CheckUserNameExists(user.UserName, user.ID)
	if err != nil {
		response.FailWithMessage(201, err.Error(), ctx)
		return
	}
	if exists {
		response.FailWithMessage(201, "用户名已存在", ctx)
		return
	}

	if err := user.Update(); err != nil {
		response.FailWithMessage(201, err.Error(), ctx)
		return
	}

	response.OkWithMessage("更新成功", ctx)
}

// Delete 删除用户
func (c *UserController) Delete(ctx *gin.Context) {
	idStr := ctx.Param("userId")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.FailWithMessage(201, "无效的用户ID", ctx)
		return
	}

	user := &SystemDB.UserDO{ID: id}
	if err := user.Delete(); err != nil {
		response.FailWithMessage(201, err.Error(), ctx)
		return
	}

	response.OkWithMessage("删除成功", ctx)
}

// UpdateStatus 更新用户状态
func (c *UserController) UpdateStatus(ctx *gin.Context) {
	var req struct {
		UserID int64  `json:"userId"`
		Status string `json:"status"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(201, err.Error(), ctx)
		return
	}

	user := &SystemDB.UserDO{
		ID:     req.UserID,
		Status: req.Status,
	}

	if err := user.UpdateStatus(); err != nil {
		response.FailWithMessage(201, err.Error(), ctx)
		return
	}

	response.OkWithMessage("状态更新成功", ctx)
}

// ResetPassword 重置密码
func (c *UserController) ResetPassword(ctx *gin.Context) {
	var req struct {
		UserID   int64  `json:"userId"`
		Password string `json:"password"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(201, err.Error(), ctx)
		return
	}

	// 加密新密码
	encryptedPassword, err := utils.EncryptPassword(req.Password)
	if err != nil {
		response.FailWithMessage(201, "密码加密失败", ctx)
		return
	}

	user := &SystemDB.UserDO{
		ID:       req.UserID,
		Password: encryptedPassword,
	}

	if err := user.UpdatePassword(); err != nil {
		response.FailWithMessage(201, err.Error(), ctx)
		return
	}

	response.OkWithMessage("密码重置成功", ctx)
}
