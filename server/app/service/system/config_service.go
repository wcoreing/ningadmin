package system

import (
	"errors"
	"time"
	"wnadmin/app/model/do/SystemDB"
	"wnadmin/app/model/dto"
)

// ConfigService 配置服务
type ConfigService struct{}

// GetConfigList 获取配置列表(树形结构)
func (s *ConfigService) GetConfigList(query *dto.ConfigQuery) ([]SystemDB.ConfigDO, int64, error) {
	configDO := &SystemDB.ConfigDO{}

	// 如果是树形查询，返回树形结构
	if query.PageNum == 0 && query.PageSize == 0 {
		configs, err := configDO.BuildTreeWithFilter(query.Keyword, query.ModuleKey, query.Status)
		return configs, int64(len(configs)), err
	}

	// 否则返回分页列表
	var configs []SystemDB.ConfigDO
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
	if query.ParentId != nil {
		db = db.Where("parent_id = ?", *query.ParentId)
	}

	var total int64
	db.Count(&total)

	offset := (query.PageNum - 1) * query.PageSize
	err := db.Offset(offset).Limit(query.PageSize).Order("sort desc, id desc").Find(&configs).Error

	return configs, total, err
}

// CreateConfig 创建配置
func (s *ConfigService) CreateConfig(req *dto.CreateConfigReq, createBy string) error {
	now := time.Now()

	// 检查父级配置
	level := 1
	if req.ParentId > 0 {
		var parent SystemDB.ConfigDO
		if err := parent.Query().Where("id = ?", req.ParentId).First(&parent).Error; err != nil {
			return errors.New("父级配置不存在")
		}
		if parent.Level >= 3 {
			return errors.New("最多支持三级配置")
		}
		level = parent.Level + 1
	}

	config := &SystemDB.ConfigDO{
		ParentId:    req.ParentId,
		ModuleKey:   req.ModuleKey,
		ModuleName:  req.ModuleName,
		ConfigName:  req.ConfigName,
		ConfigKey:   req.ConfigKey,
		ConfigType:  req.ConfigType,
		ConfigValue: req.ConfigValue,
		Level:       level,
		Sort:        req.Sort,
		Status:      req.Status,
		Remark:      req.Remark,
		CreateTime:  now,
		UpdateTime:  now,
		DelFlag:     "0",
	}

	// 设置默认状态
	if config.Status == 0 {
		config.Status = 1
	}

	// 只有当配置键不为空时才检查是否已存在（父节点可以没有配置键）
	if config.ConfigKey != "" {
		var count int64
		err := config.Query().Where("config_key = ?", config.ConfigKey).Count(&count).Error
		if err != nil {
			return err
		}
		if count > 0 {
			return ErrConfigKeyExists
		}
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

	// 检查父级配置
	if req.ParentId > 0 && req.ParentId != config.ParentId {
		var parent SystemDB.ConfigDO
		if err := parent.Query().Where("id = ?", req.ParentId).First(&parent).Error; err != nil {
			return errors.New("父级配置不存在")
		}
		if parent.Level >= 3 {
			return errors.New("最多支持三级配置")
		}
		config.Level = parent.Level + 1
	}

	config.ParentId = req.ParentId
	config.ModuleKey = req.ModuleKey
	config.ModuleName = req.ModuleName
	config.ConfigName = req.ConfigName
	config.ConfigKey = req.ConfigKey
	config.ConfigType = req.ConfigType
	config.ConfigValue = req.ConfigValue
	config.Sort = req.Sort
	config.Status = req.Status
	config.Remark = req.Remark
	config.UpdateTime = time.Now()

	// 只有当配置键不为空时才检查是否已存在（父节点可以没有配置键）
	if config.ConfigKey != "" {
		var count int64
		err = config.Query().Where("config_key = ? AND id != ?", config.ConfigKey, config.ID).Count(&count).Error
		if err != nil {
			return err
		}
		if count > 0 {
			return ErrConfigKeyExists
		}
	}

	return config.Update()
}

// DeleteConfig 删除配置
func (s *ConfigService) DeleteConfig(id int64) error {
	// 检查是否有子配置
	var count int64
	configDO := &SystemDB.ConfigDO{}
	err := configDO.Query().Where("parent_id = ?", id).Count(&count).Error
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("请先删除子配置")
	}

	config := &SystemDB.ConfigDO{ID: id}
	return config.Delete()
}

// GetConfigInfo 获取配置详情
func (s *ConfigService) GetConfigInfo(id int64) (*SystemDB.ConfigDO, error) {
	config := &SystemDB.ConfigDO{ID: id}
	err := config.First()
	return config, err
}
