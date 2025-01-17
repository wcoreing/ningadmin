package system

import (
	"time"
	"wnadmin/app/model/do/SystemDB"
	"wnadmin/app/model/dto"
)

// ConfigService 配置服务
type ConfigService struct{}

// GetConfigList 获取配置列表
func (s *ConfigService) GetConfigList(query *dto.ConfigQuery) ([]SystemDB.ConfigDO, int64, error) {
	var configs []SystemDB.ConfigDO
	configDO := &SystemDB.ConfigDO{}
	db := configDO.Query()

	if query.Keyword != "" {
		db = db.Where("config_name LIKE ? OR config_key LIKE ?", "%"+query.Keyword+"%", "%"+query.Keyword+"%")
	}
	if query.ModuleKey != "" {
		db = db.Where("module_key = ?", query.ModuleKey)
	}
	if query.Status != "" {
		db = db.Where("status = ?", query.Status)
	}

	var total int64
	db.Count(&total)

	offset := (query.PageNum - 1) * query.PageSize
	err := db.Offset(offset).Limit(query.PageSize).Order("id desc").Find(&configs).Error

	return configs, total, err
}

// CreateConfig 创建配置
func (s *ConfigService) CreateConfig(req *dto.CreateConfigReq, createBy string) error {
	now := time.Now()
	config := &SystemDB.ConfigDO{
		ModuleKey:   req.ModuleKey,
		ConfigName:  req.ConfigName,
		ConfigKey:   req.ConfigKey,
		ConfigType:  req.ConfigType,
		ConfigValue: req.ConfigValue,
		Status:      req.Status,
		Remark:      req.Remark,
		CreateTime:  now,
		UpdateTime:  now,
		DelFlag:     "0",
	}

	// 设置默认状态
	if config.Status == "" {
		config.Status = "0"
	}

	// 检查配置键是否已存在
	var count int64
	err := config.Query().Where("config_key = ?", config.ConfigKey).Count(&count).Error
	if err != nil {
		return err
	}
	if count > 0 {
		return ErrConfigKeyExists
	}

	return config.Create()
}

// UpdateConfig 更新配置
func (s *ConfigService) UpdateConfig(req *dto.UpdateConfigReq, updateBy string) error {
	config := &SystemDB.ConfigDO{}
	err := config.Query().Where("id = ?", req.ID).First(config).Error
	if err != nil {
		return err
	}

	config.ModuleKey = req.ModuleKey
	config.ConfigName = req.ConfigName
	config.ConfigKey = req.ConfigKey
	config.ConfigType = req.ConfigType
	config.ConfigValue = req.ConfigValue
	config.Status = req.Status
	config.Remark = req.Remark
	config.UpdateTime = time.Now()

	// 检查配置键是否已存在（排除自身）
	var count int64
	err = config.Query().Where("config_key = ? AND id != ?", config.ConfigKey, config.ID).Count(&count).Error
	if err != nil {
		return err
	}
	if count > 0 {
		return ErrConfigKeyExists
	}

	return config.Update()
}

// DeleteConfig 删除配置
func (s *ConfigService) DeleteConfig(id int64) error {
	config := &SystemDB.ConfigDO{ID: id}
	return config.Delete()
}

// GetConfigInfo 获取配置详情
func (s *ConfigService) GetConfigInfo(id int64) (*SystemDB.ConfigDO, error) {
	config := &SystemDB.ConfigDO{ID: id}
	err := config.First()
	return config, err
}
