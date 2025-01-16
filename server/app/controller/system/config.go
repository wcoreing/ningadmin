package system

import (
	"strconv"
	"wnadmin/app/common/response"
	"wnadmin/app/model/do/SystemDB"

	"github.com/gin-gonic/gin"
)

type ConfigController struct{}

// List 获取配置列表
func (c *ConfigController) List(ctx *gin.Context) {
	var query struct {
		PageNum   int    `form:"pageNum"`
		PageSize  int    `form:"pageSize"`
		Keyword   string `form:"keyword"`
		ModuleKey string `form:"moduleKey"`
	}

	if err := ctx.ShouldBindQuery(&query); err != nil {
		response.FailWithMessage(201, err.Error(), ctx)
		return
	}

	// 设置默认值
	if query.PageNum <= 0 {
		query.PageNum = 1
	}
	if query.PageSize <= 0 {
		query.PageSize = 10
	}

	var configs []SystemDB.ConfigDO
	configDO := &SystemDB.ConfigDO{}
	db := configDO.Query()

	if query.Keyword != "" {
		db = db.Where("config_name LIKE ? OR config_key LIKE ?", "%"+query.Keyword+"%", "%"+query.Keyword+"%")
	}
	if query.ModuleKey != "" {
		db = db.Where("module_key = ?", query.ModuleKey)
	}

	var total int64
	db.Count(&total)

	offset := (query.PageNum - 1) * query.PageSize
	err := db.Offset(offset).Limit(query.PageSize).Order("id desc").Find(&configs).Error
	if err != nil {
		response.FailWithMessage(201, err.Error(), ctx)
		return
	}

	response.OkWithDetailed(gin.H{
		"total": total,
		"list":  configs,
	}, "获取成功", ctx)
}

// Add 新增配置
func (c *ConfigController) Add(ctx *gin.Context) {
	var config SystemDB.ConfigDO
	if err := ctx.ShouldBindJSON(&config); err != nil {
		response.FailWithMessage(201, err.Error(), ctx)
		return
	}

	// 检查配置键是否已存在
	var count int64
	err := config.Query().Where("config_key = ?", config.ConfigKey).Count(&count).Error
	if err != nil {
		response.FailWithMessage(201, err.Error(), ctx)
		return
	}
	if count > 0 {
		response.FailWithMessage(201, "配置键已存在", ctx)
		return
	}

	if err := config.Create(); err != nil {
		response.FailWithMessage(201, err.Error(), ctx)
		return
	}

	response.OkWithMessage("创建成功", ctx)
}

// Update 更新配置
func (c *ConfigController) Update(ctx *gin.Context) {
	var config SystemDB.ConfigDO
	if err := ctx.ShouldBindJSON(&config); err != nil {
		response.FailWithMessage(201, err.Error(), ctx)
		return
	}

	if err := config.Update(); err != nil {
		response.FailWithMessage(201, err.Error(), ctx)
		return
	}

	response.OkWithMessage("更新成功", ctx)
}

// Delete 删除配置
func (c *ConfigController) Delete(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.FailWithMessage(201, "无效的ID", ctx)
		return
	}

	config := &SystemDB.ConfigDO{ID: id}
	if err := config.Delete(); err != nil {
		response.FailWithMessage(201, err.Error(), ctx)
		return
	}

	response.OkWithMessage("删除成功", ctx)
}

// GetInfo 获取配置详情
func (c *ConfigController) GetInfo(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.FailWithMessage(201, "无效的ID", ctx)
		return
	}

	config := &SystemDB.ConfigDO{ID: id}
	if err := config.First(); err != nil {
		response.FailWithMessage(100, err.Error(), ctx)
		return
	}

	response.OkWithDetailed(config, "获取成功", ctx)
}
