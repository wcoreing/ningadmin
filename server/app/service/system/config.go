package system

import (
	"errors"
	"wnadmin/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ConfigService struct {
	db *gorm.DB
}

func NewConfigService(db *gorm.DB) *ConfigService {
	return &ConfigService{db: db}
}

// List 获取配置列表
func (s *ConfigService) List(ctx *gin.Context, query model.SystemConfigQuery) ([]model.SystemConfig, int64, error) {
	var configs []model.SystemConfig
	var total int64

	db := s.db.Model(&model.SystemConfig{})

	if query.Keyword != "" {
		db = db.Where("config_name LIKE ? OR config_key LIKE ?", "%"+query.Keyword+"%", "%"+query.Keyword+"%")
	}
	if query.ModuleKey != "" {
		db = db.Where("module_key = ?", query.ModuleKey)
	}

	err := db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (query.PageNum - 1) * query.PageSize
	err = db.Offset(offset).Limit(query.PageSize).Order("id desc").Find(&configs).Error
	if err != nil {
		return nil, 0, err
	}

	return configs, total, nil
}

// Add 新增配置
func (s *ConfigService) Add(ctx *gin.Context, config model.SystemConfigAdd) error {
	// 检查配置键是否已存在
	var count int64
	err := s.db.Model(&model.SystemConfig{}).Where("config_key = ?", config.ConfigKey).Count(&count).Error
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("配置键已存在")
	}

	// 创建新配置
	newConfig := model.SystemConfig{
		ConfigName:  config.ConfigName,
		ConfigKey:   config.ConfigKey,
		ConfigValue: config.ConfigValue,
		ConfigType:  config.ConfigType,
		ModuleKey:   config.ModuleKey,
		Remark:      config.Remark,
	}

	return s.db.Create(&newConfig).Error
}

// Update 更新配置
func (s *ConfigService) Update(ctx *gin.Context, config model.SystemConfigUpdate) error {
	// 检查配置是否存在
	var existConfig model.SystemConfig
	err := s.db.First(&existConfig, config.ID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("配置不存在")
		}
		return err
	}

	// 更新配置
	updates := map[string]interface{}{
		"config_name":  config.ConfigName,
		"config_value": config.ConfigValue,
		"config_type":  config.ConfigType,
		"module_key":   config.ModuleKey,
		"remark":       config.Remark,
	}

	return s.db.Model(&existConfig).Updates(updates).Error
}

// Delete 删除配置
func (s *ConfigService) Delete(ctx *gin.Context, id int64) error {
	result := s.db.Delete(&model.SystemConfig{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("配置不存在")
	}
	return nil
}

// GetInfo 获取配置详情
func (s *ConfigService) GetInfo(ctx *gin.Context, id int64) (*model.SystemConfig, error) {
	var config model.SystemConfig
	err := s.db.First(&config, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("配置不存在")
		}
		return nil, err
	}
	return &config, nil
}
