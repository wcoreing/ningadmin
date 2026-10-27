package dto

// ConfigQuery 配置查询参数
type ConfigQuery struct {
	PageNum   int    `form:"pageNum"`
	PageSize  int    `form:"pageSize"`
	Keyword   string `form:"keyword"`
	ModuleKey string `form:"moduleKey"`
	Status    string `form:"status"`
	ParentId  *int64 `form:"parentId"`
}

// CreateConfigReq 创建配置请求
type CreateConfigReq struct {
	ParentId    int64  `json:"parentId"`
	ModuleKey   string `json:"moduleKey" binding:"required"`
	ModuleName  string `json:"moduleName" binding:"required"`
	ConfigName  string `json:"configName" binding:"required"`
	ConfigKey   string `json:"configKey"`   // 父节点可不填
	ConfigType  string `json:"configType"`  // 父节点可不填
	ConfigValue string `json:"configValue"` // 父节点可不填
	Sort        int    `json:"sort"`
	Status      int    `json:"status"`
	Remark      string `json:"remark"`
}

// UpdateConfigReq 更新配置请求
type UpdateConfigReq struct {
	ID          int64  `json:"id" binding:"required"`
	ParentId    int64  `json:"parentId"`
	ModuleKey   string `json:"moduleKey" binding:"required"`
	ModuleName  string `json:"moduleName" binding:"required"`
	ConfigName  string `json:"configName" binding:"required"`
	ConfigKey   string `json:"configKey"`   // 父节点可不填
	ConfigType  string `json:"configType"`  // 父节点可不填
	ConfigValue string `json:"configValue"` // 父节点可不填
	Sort        int    `json:"sort"`
	Status      int    `json:"status"`
	Remark      string `json:"remark"`
}
