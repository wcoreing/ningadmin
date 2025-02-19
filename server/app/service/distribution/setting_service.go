package distribution

import (
	"errors"
	"wnadmin/app/model/do/distribution"
	"wnadmin/app/model/dto"
	"wnadmin/app/provider/db"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SettingService 分销设置服务
type SettingService struct {
	db *gorm.DB
}

// NewSettingService 创建分销设置服务实例
func NewSettingService() *SettingService {
	return &SettingService{
		db: db.DB.WnAdmin,
	}
}

// GetSetting 获取分销设置
func (s *SettingService) GetSetting(c *gin.Context) (*distribution.DistributionSetting, error) {
	var setting distribution.DistributionSetting
	err := s.db.Model(&distribution.DistributionSetting{}).First(&setting).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("分销设置不存在")
		}
		return nil, err
	}
	return &setting, nil
}

// UpdateSetting 更新分销设置
func (s *SettingService) UpdateSetting(c *gin.Context, req *dto.DistributionSettingDTO) error {
	setting := &distribution.DistributionSetting{
		Enabled:                req.Enabled,
		Level:                  req.Level,
		Mode:                   req.Mode,
		ApplyEnabled:           req.ApplyEnabled,
		UserType:               req.UserType,
		BindMode:               req.BindMode,
		CenterType:             req.CenterType,
		PosterImages:           req.PosterImages,
		CalcType:               req.CalcType,
		Level1Rate:             req.Level1Rate,
		Level2Rate:             req.Level2Rate,
		FreezeDays:             req.FreezeDays,
		SelfBuyEnabled:         req.SelfBuyEnabled,
		PromotionEnabled:       req.PromotionEnabled,
		MaxWithdraw:            req.MaxWithdraw,
		MinWithdraw:            req.MinWithdraw,
		WithdrawFeeRate:        req.WithdrawFeeRate,
		BankWithdrawEnabled:    req.BankWithdrawEnabled,
		WechatWithdrawEnabled:  req.WechatWithdrawEnabled,
		WechatPayMode:          req.WechatPayMode,
		AlipayWithdrawEnabled:  req.AlipayWithdrawEnabled,
		BalanceWithdrawEnabled: req.BalanceWithdrawEnabled,
	}

	// 检查是否存在记录
	var count int64
	err := s.db.Model(&distribution.DistributionSetting{}).Count(&count).Error
	if err != nil {
		return errors.New("查询分销设置失败: " + err.Error())
	}

	if count == 0 {
		// 不存在则创建
		if err := s.db.Create(setting).Error; err != nil {
			return errors.New("创建分销设置失败: " + err.Error())
		}
	} else {
		// 存在则更新
		if err := s.db.Model(&distribution.DistributionSetting{}).Where("id = ?", 1).Updates(setting).Error; err != nil {
			return errors.New("更新分销设置失败: " + err.Error())
		}
	}

	return nil
}
