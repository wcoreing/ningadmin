package distribution

import (
	"time"
	"wnadmin/app/model/do/distribution"
	"wnadmin/app/model/do/shop"
	"wnadmin/app/provider/db"

	"gorm.io/gorm"
)

type OrderListenerService struct {
	db            *gorm.DB
	incomeService *IncomeService
}

func NewOrderListenerService() *OrderListenerService {
	return &OrderListenerService{
		db:            db.DB.WnAdmin,
		incomeService: NewIncomeService(),
	}
}

// HandleOrderComplete 处理订单完成事件
func (s *OrderListenerService) HandleOrderComplete(orderID uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 1. 查询订单信息
		var order shop.ShopOrder
		if err := tx.Where("id = ?", orderID).First(&order).Error; err != nil {
			return err
		}

		// 2. 查询分销关系
		var relation distribution.DistributionRelation
		if err := tx.Where("user_id = ? AND status = 1", order.UserId).First(&relation).Error; err != nil {
			return nil // 无分销关系，直接返回
		}

		// 3. 查询分销设置
		var setting distribution.DistributionSetting
		if err := tx.First(&setting).Error; err != nil {
			return err
		}

		// 4. 计算一级分销收益
		level1Amount := order.PayAmount * setting.Level1Rate / 100
		level1Income := &distribution.DistributionIncome{
			UserID:           relation.RelatedUserID,
			OrderID:          orderID,
			OrderAmount:      order.PayAmount,
			CommissionRate:   setting.Level1Rate,
			CommissionAmount: level1Amount,
			Level:            1,
			CommissionStatus: 0, // 待结算
		}

		// 5. 如果有冻结期，设置解冻时间
		if setting.FreezeDays > 0 {
			unfreezeTime := time.Now().AddDate(0, 0, setting.FreezeDays)
			level1Income.UnfreezeTime = &unfreezeTime
		} else {
			// 无冻结期，直接入账
			if err := s.incomeService.AddIncome(level1Income); err != nil {
				return err
			}
		}

		// 6. 如果是二级分销，处理二级收益
		if setting.Level == 2 {
			var level2Relation distribution.DistributionRelation
			if err := tx.Where("user_id = ? AND level = 2 AND status = 1", order.UserId).First(&level2Relation).Error; err == nil {
				level2Amount := order.PayAmount * setting.Level2Rate / 100
				level2Income := &distribution.DistributionIncome{
					UserID:           level2Relation.RelatedUserID,
					OrderID:          orderID,
					OrderAmount:      order.PayAmount,
					CommissionRate:   setting.Level2Rate,
					CommissionAmount: level2Amount,
					Level:            2,
					CommissionStatus: 0, // 待结算
				}

				if setting.FreezeDays > 0 {
					unfreezeTime2 := time.Now().AddDate(0, 0, setting.FreezeDays)
					level2Income.UnfreezeTime = &unfreezeTime2
				} else {
					// 无冻结期，直接入账
					if err := s.incomeService.AddIncome(level2Income); err != nil {
						return err
					}
				}

				if err := tx.Create(level2Income).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})
}
