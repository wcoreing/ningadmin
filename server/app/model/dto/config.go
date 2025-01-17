package dto

// ConfigQuery 配置查询参数
type ConfigQuery struct {
	PageNum   int    `form:"pageNum"`
	PageSize  int    `form:"pageSize"`
	Keyword   string `form:"keyword"`
	ModuleKey string `form:"moduleKey"`
	Status    string `form:"status"`
}

// CreateConfigReq 创建配置请求
type CreateConfigReq struct {
	ModuleKey   string `json:"moduleKey" binding:"required"`
	ConfigName  string `json:"configName" binding:"required"`
	ConfigKey   string `json:"configKey" binding:"required"`
	ConfigType  string `json:"configType" binding:"required"`
	ConfigValue string `json:"configValue"`
	Status      string `json:"status"`
	Remark      string `json:"remark"`
}

// UpdateConfigReq 更新配置请求
type UpdateConfigReq struct {
	ID          int64  `json:"id" binding:"required"`
	ModuleKey   string `json:"moduleKey" binding:"required"`
	ConfigName  string `json:"configName" binding:"required"`
	ConfigKey   string `json:"configKey" binding:"required"`
	ConfigType  string `json:"configType" binding:"required"`
	ConfigValue string `json:"configValue"`
	Status      string `json:"status"`
	Remark      string `json:"remark"`
}
