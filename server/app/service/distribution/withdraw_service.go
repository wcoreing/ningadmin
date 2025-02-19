package distribution

import (
	"time"
	"wnadmin/app/model/do/distribution"
	"wnadmin/app/model/dto"
	"wnadmin/app/provider/db"

	"errors"

	"gorm.io/gorm"
)

type WithdrawService struct {
	db *gorm.DB
}

func NewWithdrawService() *WithdrawService {
	return &WithdrawService{
		db: db.DB.WnAdmin,
	}
}

// GetWithdrawList 获取提现列表
func (s *WithdrawService) GetWithdrawList(req *dto.WithdrawListReq) (voList []dto.WithdrawVO, total int64, err error) {
	db := s.db.Model(&distribution.DistributionWithdraw{})

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

	var list []distribution.DistributionWithdraw
	err = db.Offset((req.PageNum - 1) * req.PageSize).
		Limit(req.PageSize).
		Order("id DESC").
		Find(&list).Error
	if err != nil {
		return
	}

	for _, item := range list {
		vo := dto.WithdrawVO{
			ID:          item.ID,
			UserID:      item.UserID,
			Amount:      item.Amount,
			Method:      item.Method,
			Account:     item.Account,
			RealName:    item.RealName,
			Status:      item.Status,
			AuditRemark: item.AuditRemark,
			PayRemark:   item.PayRemark,
		}

		if !item.CreateTime.IsZero() {
			vo.CreateTime = item.CreateTime.Format("2006-01-02 15:04:05")
		}
		if !item.AuditTime.IsZero() {
			vo.AuditTime = item.AuditTime.Format("2006-01-02 15:04:05")
		}
		if !item.PayTime.IsZero() {
			vo.PayTime = item.PayTime.Format("2006-01-02 15:04:05")
		}

		voList = append(voList, vo)
	}

	return
}

// GetWithdrawDetail 获取提现详情
func (s *WithdrawService) GetWithdrawDetail(id string) (*dto.WithdrawVO, error) {
	var withdraw distribution.DistributionWithdraw
	err := s.db.Where("id = ?", id).First(&withdraw).Error
	if err != nil {
		return nil, err
	}

	vo := &dto.WithdrawVO{
		ID:          withdraw.ID,
		UserID:      withdraw.UserID,
		Amount:      withdraw.Amount,
		Method:      withdraw.Method,
		Account:     withdraw.Account,
		RealName:    withdraw.RealName,
		Status:      withdraw.Status,
		AuditRemark: withdraw.AuditRemark,
		PayRemark:   withdraw.PayRemark,
	}

	if !withdraw.CreateTime.IsZero() {
		vo.CreateTime = withdraw.CreateTime.Format("2006-01-02 15:04:05")
	}
	if !withdraw.AuditTime.IsZero() {
		vo.AuditTime = withdraw.AuditTime.Format("2006-01-02 15:04:05")
	}
	if !withdraw.PayTime.IsZero() {
		vo.PayTime = withdraw.PayTime.Format("2006-01-02 15:04:05")
	}

	return vo, nil
}

// ApproveWithdraw 审核通过
func (s *WithdrawService) ApproveWithdraw(id string) error {
	return s.db.Model(&distribution.DistributionWithdraw{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":     1,
			"audit_time": time.Now(),
		}).Error
}

// RejectWithdraw 审核拒绝
func (s *WithdrawService) RejectWithdraw(req *dto.WithdrawRejectReq) error {
	return s.db.Model(&distribution.DistributionWithdraw{}).
		Where("id = ?", req.ID).
		Updates(map[string]interface{}{
			"status":       2,
			"audit_time":   time.Now(),
			"audit_remark": req.Remark,
		}).Error
}

// ConfirmPay 确认打款
func (s *WithdrawService) ConfirmPay(req *dto.WithdrawPayReq) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 1. 查询提现记录
		var withdraw distribution.DistributionWithdraw
		if err := tx.Where("id = ?", req.ID).First(&withdraw).Error; err != nil {
			return err
		}

		// 2. 使用 FOR UPDATE 锁定用户记录
		var user distribution.DistributionUser
		if err := tx.Set("gorm:query_option", "FOR UPDATE").
			Where("user_id = ?", withdraw.UserID).
			First(&user).Error; err != nil {
			return err
		}

		// 3. 检查余额
		if user.AvailableBalance < withdraw.Amount {
			return errors.New("用户可用余额不足")
		}

		// 4. 计算新余额
		newBalance := user.AvailableBalance - withdraw.Amount

		// 5. 使用乐观锁更新用户余额
		result := tx.Model(&user).
			Where("version = ?", user.Version).
			Updates(map[string]interface{}{
				"available_balance": newBalance,
				"total_withdraw":    gorm.Expr("total_withdraw + ?", withdraw.Amount),
				"version":           gorm.Expr("version + 1"),
			})

		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return errors.New("余额更新失败，请重试")
		}

		// 6. 记录资金变动
		balanceLog := &distribution.DistributionBalanceLog{
			UserID:        withdraw.UserID,
			Type:          2, // 支出
			Amount:        withdraw.Amount,
			BeforeBalance: user.AvailableBalance,
			AfterBalance:  newBalance,
			SourceType:    2, // 提现
			SourceID:      &withdraw.ID,
			Remark:        "提现支出",
		}
		if err := tx.Create(balanceLog).Error; err != nil {
			return err
		}

		// 7. 更新提现状态
		if err := tx.Model(&withdraw).Updates(map[string]interface{}{
			"status":     3,
			"pay_time":   time.Now(),
			"pay_remark": req.Remark,
		}).Error; err != nil {
			return err
		}

		return nil
	})
}
