package distribution

import (
	"errors"
	"fmt"
	"time"
	"wnadmin/app/model/do/distribution"
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
	err := s.db.Where("status = 1").Find(&incomes).Error
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
