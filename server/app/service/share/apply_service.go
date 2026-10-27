package share

import (
	"time"
	"wnadmin/app/model/do/share"
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

// GetApplyList 获取分享申请列表
func (s *ApplyService) GetApplyList(req *dto.ShareApplyListReq) (voList []dto.ShareApplyVO, total int64, err error) {
	db := s.db.Model(&share.ShareUser{})

	// 关联user表获取用户基本信息
	db = db.Select(`
		share_user.*,
		user.username,
		user.nickname,
		user.avatar,
		user.phone
	`).
		Joins("LEFT JOIN user ON user.id = share_user.user_id")

	if req.RealName != "" {
		db = db.Where("share_user.real_name LIKE ?", "%"+req.RealName+"%")
	}
	if req.Phone != "" {
		db = db.Where("share_user.phone = ?", req.Phone)
	}
	if req.Status != nil {
		db = db.Where("share_user.status = ?", *req.Status)
	}

	err = db.Count(&total).Error
	if err != nil {
		return
	}

	var list []share.ShareUser
	err = db.Offset((req.PageNum - 1) * req.PageSize).
		Limit(req.PageSize).
		Order("share_user.create_time DESC").
		Find(&list).Error
	if err != nil {
		return
	}

	// 转换为VO对象
	for _, item := range list {
		vo := dto.ShareApplyVO{
			ID:           item.ID,
			UserID:       item.UserID,
			Username:     item.Username, // 从user表读取
			Nickname:     item.Nickname, // 从user表读取
			Avatar:       item.Avatar,   // 从user表读取
			RealName:     item.RealName,
			Phone:        item.Phone,
			IDCard:       item.IDCard,
			Status:       item.Status,
			AuditUser:    item.AuditUser,
			RejectReason: item.RejectReason,
		}
		// 日期格式化
		if item.ApplyTime != nil {
			vo.ApplyTime = item.ApplyTime.Format("2006-01-02 15:04:05")
		}
		if item.AuditTime != nil {
			vo.AuditTime = item.AuditTime.Format("2006-01-02 15:04:05")
		}
		voList = append(voList, vo)
	}

	return
}

// ApproveApply 审核通过
func (s *ApplyService) ApproveApply(id string) error {
	return s.db.Model(&share.ShareUser{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":     1,
			"audit_time": time.Now(),
		}).Error
}

// RejectApply 审核拒绝
func (s *ApplyService) RejectApply(req *dto.ShareApplyRejectReq) error {
	return s.db.Model(&share.ShareUser{}).
		Where("id = ?", req.ID).
		Updates(map[string]interface{}{
			"status":        2,
			"audit_time":    time.Now(),
			"reject_reason": req.Reason,
		}).Error
}

// GetApplyDetail 获取申请详情
func (s *ApplyService) GetApplyDetail(id string) (*dto.ShareApplyVO, error) {
	var user share.ShareUser
	err := s.db.Select(`
		share_user.*,
		user.username,
		user.nickname,
		user.avatar
	`).
		Joins("LEFT JOIN user ON user.id = share_user.user_id").
		Where("share_user.id = ?", id).
		First(&user).Error
	if err != nil {
		return nil, err
	}

	vo := &dto.ShareApplyVO{
		ID:           user.ID,
		UserID:       user.UserID,
		Username:     user.Username, // 从user表读取
		Nickname:     user.Nickname, // 从user表读取
		Avatar:       user.Avatar,   // 从user表读取
		RealName:     user.RealName,
		Phone:        user.Phone,
		IDCard:       user.IDCard,
		Status:       user.Status,
		AuditUser:    user.AuditUser,
		RejectReason: user.RejectReason,
	}

	if user.ApplyTime != nil {
		vo.ApplyTime = user.ApplyTime.Format("2006-01-02 15:04:05")
	}
	if user.AuditTime != nil {
		vo.AuditTime = user.AuditTime.Format("2006-01-02 15:04:05")
	}

	return vo, nil
}
