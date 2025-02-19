package distribution

import (
	"time"
	"wnadmin/app/model/do/distribution"
	"wnadmin/app/model/dto"
	"wnadmin/app/provider/db"

	"gorm.io/gorm"
)

type ApplyService struct {
	db *gorm.DB
}

func NewApplyService() *ApplyService {
	return &ApplyService{
		db: db.DB.WnAdmin,
	}
}

// GetApplyList 获取分销员申请列表
func (s *ApplyService) GetApplyList(req *dto.DistributionApplyListReq) (voList []dto.DistributionApplyVO, total int64, err error) {
	db := s.db.Model(&distribution.DistributionUser{})

	if req.RealName != "" {
		db = db.Where("real_name LIKE ?", "%"+req.RealName+"%")
	}
	if req.Phone != "" {
		db = db.Where("phone = ?", req.Phone)
	}
	if req.Status != nil {
		db = db.Where("status = ?", *req.Status)
	}

	err = db.Count(&total).Error
	if err != nil {
		return
	}

	var list []distribution.DistributionUser
	err = db.Offset((req.PageNum - 1) * req.PageSize).
		Limit(req.PageSize).
		Order("create_time DESC").
		Find(&list).Error
	if err != nil {
		return
	}

	// 转换为VO对象
	for _, item := range list {
		vo := dto.DistributionApplyVO{
			ID:           item.ID,
			UserID:       item.UserID,
			RealName:     item.RealName,
			Phone:        item.Phone,
			IDCard:       item.IDCard,
			Status:       item.Status,
			AuditUser:    item.AuditUser,
			RejectReason: item.RejectReason,
		}
		// 日期格式化
		if !item.ApplyTime.IsZero() {
			vo.ApplyTime = item.ApplyTime.Format("2006-01-02 15:04:05")
		}
		if !item.AuditTime.IsZero() {
			vo.AuditTime = item.AuditTime.Format("2006-01-02 15:04:05")
		}
		voList = append(voList, vo)
	}

	return
}

// ApproveApply 审核通过
func (s *ApplyService) ApproveApply(id string) error {
	return s.db.Model(&distribution.DistributionUser{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":     1,
			"audit_time": time.Now(),
		}).Error
}

// RejectApply 审核拒绝
func (s *ApplyService) RejectApply(req *dto.DistributionApplyRejectReq) error {
	return s.db.Model(&distribution.DistributionUser{}).
		Where("id = ?", req.ID).
		Updates(map[string]interface{}{
			"status":        2,
			"audit_time":    time.Now(),
			"reject_reason": req.Reason,
		}).Error
}

// GetApplyDetail 获取申请详情
func (s *ApplyService) GetApplyDetail(id string) (*dto.DistributionApplyVO, error) {
	var user distribution.DistributionUser
	err := s.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}

	vo := &dto.DistributionApplyVO{
		ID:           user.ID,
		UserID:       user.UserID,
		RealName:     user.RealName,
		Phone:        user.Phone,
		IDCard:       user.IDCard,
		Status:       user.Status,
		AuditUser:    user.AuditUser,
		RejectReason: user.RejectReason,
	}

	if !user.ApplyTime.IsZero() {
		vo.ApplyTime = user.ApplyTime.Format("2006-01-02 15:04:05")
	}
	if !user.AuditTime.IsZero() {
		vo.AuditTime = user.AuditTime.Format("2006-01-02 15:04:05")
	}

	return vo, nil
}
