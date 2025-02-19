package distribution

import (
	"errors"
	"time"
	"wnadmin/app/model/do/shop"
	"wnadmin/app/model/dto"
	"wnadmin/app/provider/db"
)

type DistributionService struct{}

func NewDistributionService() *DistributionService {
	return &DistributionService{}
}

// GetConfig 获取分销配置
func (s *DistributionService) GetConfig() (*dto.DistributionConfigDTO, error) {
	var config shop.DistributionConfig
	if err := db.DB.WnAdmin.First(&config).Error; err != nil {
		return nil, err
	}
	return &dto.DistributionConfigDTO{
		ID:           config.ID,
		Enabled:      config.Enabled == 1,
		Level:        config.Level,
		Mode:         config.Mode,
		ApplyEnabled: config.ApplyEnabled == 1,
		BindType:     config.BindType,
		BindMode:     config.BindMode,
		PosterBg:     config.PosterBg,
		CenterType:   config.CenterType,
	}, nil
}

// UpdateConfig 更新分销配置
func (s *DistributionService) UpdateConfig(req *dto.DistributionConfigDTO) error {
	config := &shop.DistributionConfig{
		ID:           req.ID,
		Enabled:      boolToInt8(req.Enabled),
		Level:        req.Level,
		Mode:         req.Mode,
		ApplyEnabled: boolToInt8(req.ApplyEnabled),
		BindType:     req.BindType,
		BindMode:     req.BindMode,
		PosterBg:     req.PosterBg,
		CenterType:   req.CenterType,
		UpdateTime:   time.Now(),
	}
	return db.DB.WnAdmin.Save(config).Error
}

// GetDistributorList 获取分销员列表
func (s *DistributionService) GetDistributorList(query *dto.DistributorQueryDTO) ([]*dto.DistributorDTO, int64, error) {
	var (
		list  []*shop.Distributor
		total int64
	)

	db := db.DB.WnAdmin.Model(&shop.Distributor{})
	if query.Status > 0 {
		db = db.Where("status = ?", query.Status)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := db.Offset((query.PageNum - 1) * query.PageSize).
		Limit(query.PageSize).
		Find(&list).Error; err != nil {
		return nil, 0, err
	}

	result := make([]*dto.DistributorDTO, 0, len(list))
	for _, item := range list {
		result = append(result, &dto.DistributorDTO{
			ID:          item.ID,
			UserID:      item.UserID,
			ParentID:    item.ParentID,
			Level:       item.Level,
			Status:      item.Status,
			TotalIncome: item.TotalIncome,
			TotalOrders: item.TotalOrders,
			BindTime:    item.BindTime,
			ExpireTime:  item.ExpireTime,
		})
	}

	return result, total, nil
}

// AuditDistributor 审核分销员
func (s *DistributionService) AuditDistributor(req *dto.DistributorAuditDTO) error {
	if req.Status != 1 && req.Status != 2 {
		return errors.New("无效的审核状态")
	}

	return db.DB.WnAdmin.Model(&shop.Distributor{}).
		Where("id = ?", req.ID).
		Updates(map[string]interface{}{
			"status":      req.Status,
			"update_time": time.Now(),
		}).Error
}

func boolToInt8(b bool) int8 {
	if b {
		return 1
	}
	return 0
}
