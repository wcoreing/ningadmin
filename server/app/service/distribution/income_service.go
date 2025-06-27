package distribution

import (
	"errors"
	"fmt"
	"time"
	"wnadmin/app/model/do/distribution"
	"wnadmin/app/model/dto"
	distributionvo "wnadmin/app/model/vo/distribution"

	"wnadmin/app/provider/db"

	"gorm.io/gorm"
)

type IncomeService struct {
	db *gorm.DB
}

func NewIncomeService() *IncomeService {
	return &IncomeService{
		db: db.DB.WnAdmin,
	}
}

// AddIncome 添加分销收益
func (s *IncomeService) AddIncome(income *distribution.DistributionIncome) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 1. 检查是否已处理过
		var existIncome distribution.DistributionIncome
		err := tx.Where("id = ?", income.ID).First(&existIncome).Error
		if err == nil && existIncome.Status >= 3 {
			return fmt.Errorf("该收益记录已入账，状态：%d", existIncome.Status)
		}

		// 2. 查询用户当前余额
		var user distribution.DistributionUser
		if err := tx.Set("gorm:query_option", "FOR UPDATE").
			Where("user_id = ?", income.UserID).First(&user).Error; err != nil {
			return err
		}

		// 3. 计算新余额
		newBalance := user.AvailableBalance + income.CommissionAmount

		// 4. 使用乐观锁更新用户余额
		result := tx.Model(&user).
			Where("version = ?", user.Version).
			Updates(map[string]interface{}{
				"available_balance": newBalance,
				"total_income":      gorm.Expr("total_income + ?", income.CommissionAmount),
				"version":           gorm.Expr("version + 1"),
			})

		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return errors.New("余额更新失败，请重试")
		}

		// 5. 更新收益状态
		if err := tx.Model(income).Updates(map[string]interface{}{
			"status":         3, // 已入账
			"process_time":   time.Now(),
			"process_remark": "收益入账成功",
		}).Error; err != nil {
			return err
		}

		// 6. 记录资金变动日志
		balanceLog := &distribution.DistributionBalanceLog{
			UserID:        income.UserID,
			Type:          1,
			Amount:        income.CommissionAmount,
			BeforeBalance: user.AvailableBalance,
			AfterBalance:  newBalance,
			SourceType:    1,
			SourceID:      &income.ID,
			Remark:        fmt.Sprintf("订单[%d]分销收益", income.OrderID),
		}
		if err := tx.Create(balanceLog).Error; err != nil {
			return err
		}

		return nil
	})
}

// Unfreeze 解冻分销收益
func (s *IncomeService) Unfreeze(req *dto.DistributionIncomeUnfreezeReq) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		return s.UnfreezeIncome(uint(req.Id))
	})
}

// UnfreezeIncome 解冻收益
func (s *IncomeService) UnfreezeIncome(incomeID uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 1. 查询收益记录
		var income distribution.DistributionIncome
		if err := tx.Where("id = ? AND status = 0", incomeID).First(&income).Error; err != nil {
			return err
		}

		// 2. 查询用户信息
		var user distribution.DistributionUser
		if err := tx.Where("user_id = ?", income.UserID).First(&user).Error; err != nil {
			return err
		}

		// 3. 计算新余额
		newBalance := user.AvailableBalance + income.CommissionAmount

		// 4. 使用乐观锁更新用户余额
		result := tx.Model(&user).
			Where("version = ?", user.Version).
			Updates(map[string]interface{}{
				"available_balance": newBalance,
				"frozen_balance":    gorm.Expr("frozen_balance - ?", income.CommissionAmount),
				"version":           gorm.Expr("version + 1"),
			})

		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return errors.New("余额更新失败，请重试")
		}

		// 5. 更新收益状态
		if err := tx.Model(&income).Updates(map[string]interface{}{
			"status":        1, // 已解冻
			"unfreeze_time": time.Now(),
		}).Error; err != nil {
			return err
		}

		// 6. 记录资金变动日志
		balanceLog := &distribution.DistributionBalanceLog{
			UserID:        income.UserID,
			Type:          1, // 收入
			Amount:        income.CommissionAmount,
			BeforeBalance: user.AvailableBalance,
			AfterBalance:  newBalance,
			SourceType:    1, // 分销收益
			SourceID:      &income.ID,
			Remark:        "分销订单收益解冻",
		}
		if err := tx.Create(balanceLog).Error; err != nil {
			return err
		}

		return nil
	})
}

// CheckIncomeProcess 检查收益处理状态
func (s *IncomeService) CheckIncomeProcess() error {
	var incomes []distribution.DistributionIncome

	// 查找所有已解冻但未入账的记录
	err := s.db.Where("commission_status = 1").Find(&incomes).Error
	if err != nil {
		return err
	}

	for _, income := range incomes {
		// 检查是否有对应的资金变动记录
		var count int64
		err := s.db.Model(&distribution.DistributionBalanceLog{}).
			Where("source_type = 1 AND source_id = ?", income.ID).
			Count(&count).Error
		if err != nil {
			continue
		}

		if count > 0 {
			// 已有资金变动记录，标记为已入账
			s.db.Model(&income).Updates(map[string]interface{}{
				"status":         3,
				"process_time":   time.Now(),
				"process_remark": "系统自动标记-已入账",
			})
		}
	}

	return nil
}

// GetPage 分页获取分销收益列表
func (s *IncomeService) GetIncomList(req *dto.DistributionIncomeListReq) (list []distributionvo.IncomeVO, total int64, err error) {
	var incomes []distribution.DistributionIncome
	db := s.db.Model(&distribution.DistributionIncome{})

	// 关联用户表
	db = db.Select(`
		share_income.*,
		shop_user.nickname as user_name
	`).
		Joins("LEFT JOIN shop_user ON shop_user.id = share_income.user_id")

	// 构建查询条件
	if req.UserId != "" {
		db = db.Where("share_income.user_id = ?", req.UserId)
	}
	if req.Nickname != "" {
		db = db.Where("shop_user.nickname LIKE ?", "%"+req.Nickname+"%")
	}
	if req.Status != "" {
		db = db.Where("share_income.commission_status = ?", req.Status)
	}
	if req.StartTime != "" {
		db = db.Where("share_income.create_time >= ?", req.StartTime)
	}
	if req.EndTime != "" {
		db = db.Where("share_income.create_time <= ?", req.EndTime)
	}

	// 获取总数
	err = db.Count(&total).Error
	if err != nil {
		return
	}

	// 获取分页数据
	err = db.Offset((req.PageNum - 1) * req.PageSize).
		Limit(req.PageSize).
		Order("share_income.id DESC").
		Find(&incomes).Error
	// 转换为 VO
	for _, income := range incomes {
		vo := distributionvo.IncomeVO{
			Id:               income.ID,
			UserId:           income.UserID,
			UserName:         income.UserName,
			OrderId:          income.OrderID,
			OrderNo:          income.OrderNo,
			OrderAmount:      income.OrderAmount,
			CommissionRate:   income.CommissionRate,
			CommissionAmount: income.CommissionAmount,
			Level:            income.Level,
			Status:           income.CommissionStatus,
			CreateTime:       income.CreateTime.Format("2006-01-02 15:04:05"),
		}

		// 处理解冻时间
		if income.UnfreezeTime != nil {
			vo.UnfreezeTime = *income.UnfreezeTime
		}

		// 处理状态显示
		switch income.CommissionStatus {
		case 0:
			vo.StatusText = "待结算"
		case 1:
			vo.StatusText = "已结算"
		case 2:
			vo.StatusText = "已取消"
		}

		list = append(list, vo)
	}

	return
}

// SettleIncome 结算单个收益
func (s *IncomeService) SettleIncome(incomeID uint, operator string) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 1. 查询收益记录
		var income distribution.DistributionIncome
		if err := tx.Where("id = ? AND commission_status IN (0, 1)", incomeID).First(&income).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return errors.New("收益记录不存在或状态不允许结算")
			}
			return err
		}

		// 2. 查询用户信息
		var user distribution.DistributionUser
		if err := tx.Set("gorm:query_option", "FOR UPDATE").
			Where("user_id = ?", income.UserID).First(&user).Error; err != nil {
			return err
		}

		// 3. 根据当前状态进行不同处理
		var newAvailableBalance, newFrozenBalance float64
		var remark string

		if income.CommissionStatus == 0 {
			// 冻结中 -> 直接结算到可用余额
			newAvailableBalance = user.AvailableBalance + income.CommissionAmount
			newFrozenBalance = user.FrozenBalance
			remark = "冻结收益直接结算"
		} else {
			// 已解冻 -> 从冻结余额转到可用余额（实际上已解冻的应该已经在可用余额中）
			newAvailableBalance = user.AvailableBalance
			newFrozenBalance = user.FrozenBalance
			remark = "已解冻收益结算"
		}

		// 4. 使用乐观锁更新用户余额
		result := tx.Model(&user).
			Where("version = ?", user.Version).
			Updates(map[string]interface{}{
				"available_balance": newAvailableBalance,
				"frozen_balance":    newFrozenBalance,
				"total_income":      gorm.Expr("total_income + ?", income.CommissionAmount),
				"version":           gorm.Expr("version + 1"),
			})

		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return errors.New("余额更新失败，请重试")
		}

		// 5. 更新收益状态为已入账
		if err := tx.Model(&income).Updates(map[string]interface{}{
			"commission_status": 1, // 已结算
			"process_time":      time.Now(),
			"remark":            fmt.Sprintf("手动结算 - 操作人: %s", operator),
		}).Error; err != nil {
			return err
		}

		// 6. 记录资金变动日志
		balanceLog := &distribution.DistributionBalanceLog{
			UserID:        income.UserID,
			Type:          1, // 收入
			Amount:        income.CommissionAmount,
			BeforeBalance: user.AvailableBalance,
			AfterBalance:  newAvailableBalance,
			SourceType:    1, // 分销收益
			SourceID:      &income.ID,
			Remark:        remark,
		}
		if err := tx.Create(balanceLog).Error; err != nil {
			return err
		}

		return nil
	})
}

// BatchSettleIncome 批量结算收益
func (s *IncomeService) BatchSettleIncome(incomeIDs []uint, operator string) error {
	if len(incomeIDs) == 0 {
		return errors.New("请选择要结算的收益记录")
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		for _, incomeID := range incomeIDs {
			if err := s.settleIncomeInTx(tx, incomeID, operator); err != nil {
				return fmt.Errorf("结算收益ID %d 失败: %v", incomeID, err)
			}
		}
		return nil
	})
}

// settleIncomeInTx 在事务中结算单个收益
func (s *IncomeService) settleIncomeInTx(tx *gorm.DB, incomeID uint, operator string) error {
	// 1. 查询收益记录
	var income distribution.DistributionIncome
	if err := tx.Where("id = ? AND commission_status IN (0)", incomeID).First(&income).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("收益记录不存在或状态不允许结算")
		}
		return err
	}

	// 2. 查询用户信息
	var user distribution.DistributionUser
	if err := tx.Set("gorm:query_option", "FOR UPDATE").
		Where("user_id = ?", income.UserID).First(&user).Error; err != nil {
		return err
	}

	// 3. 计算新余额 - 直接结算到可用余额
	newAvailableBalance := user.AvailableBalance + income.CommissionAmount
	remark := "批量结算-待结算收益"

	// 4. 使用乐观锁更新用户余额
	result := tx.Model(&user).
		Where("version = ?", user.Version).
		Updates(map[string]interface{}{
			"available_balance": newAvailableBalance,
			"total_income":      gorm.Expr("total_income + ?", income.CommissionAmount),
			"version":           gorm.Expr("version + 1"),
		})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("余额更新失败，请重试")
	}

	// 5. 更新收益状态为已结算
	if err := tx.Model(&income).Updates(map[string]interface{}{
		"commission_status": 1, // 已结算
		"process_time":      time.Now(),
		"remark":            fmt.Sprintf("批量结算 - 操作人: %s", operator),
	}).Error; err != nil {
		return err
	}

	// 6. 记录资金变动日志
	balanceLog := &distribution.DistributionBalanceLog{
		UserID:        income.UserID,
		Type:          1, // 收入
		Amount:        income.CommissionAmount,
		BeforeBalance: user.AvailableBalance,
		AfterBalance:  newAvailableBalance,
		SourceType:    1, // 分销收益
		SourceID:      &income.ID,
		Remark:        remark,
	}
	if err := tx.Create(balanceLog).Error; err != nil {
		return err
	}

	return nil
}
