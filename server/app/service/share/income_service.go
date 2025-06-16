package share

import (
	"errors"
	"fmt"
	"time"
	"wnadmin/app/model/do/share"
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

// AddIncome 添加分享收益
func (s *IncomeService) AddIncome(income *share.ShareIncome) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 1. 检查是否已处理过
		var existIncome share.ShareIncome
		err := tx.Where("id = ?", income.ID).First(&existIncome).Error
		if err == nil && existIncome.CommissionStatus >= 1 {
			return fmt.Errorf("该收益记录已处理，状态：%d", existIncome.CommissionStatus)
		}

		// 2. 查询用户当前余额
		var user share.ShareUser
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
			"commission_status": 1, // 已结算
			"process_time":      time.Now(),
			"remark":            "收益入账成功",
		}).Error; err != nil {
			return err
		}

		// 6. 记录资金变动日志
		balanceLog := &share.ShareBalanceLog{
			UserID:        income.UserID,
			Type:          1,
			Amount:        income.CommissionAmount,
			BeforeBalance: user.AvailableBalance,
			AfterBalance:  newBalance,
			SourceType:    1,
			SourceID:      &income.ID,
			Remark:        fmt.Sprintf("订单[%s]分享收益", income.SourceOrderNo),
		}
		if err := tx.Create(balanceLog).Error; err != nil {
			return err
		}

		return nil
	})
}

// Unfreeze 解冻分享收益
func (s *IncomeService) Unfreeze(req *dto.ShareIncomeUnfreezeReq) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		return s.UnfreezeIncome(uint(req.Id))
	})
}

// UnfreezeIncome 解冻收益
func (s *IncomeService) UnfreezeIncome(incomeID uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 1. 查询收益记录
		var income share.ShareIncome
		if err := tx.Where("id = ? AND commission_status = 0", incomeID).First(&income).Error; err != nil {
			return err
		}

		// 2. 查询用户信息
		var user share.ShareUser
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
			"commission_status": 1, // 已结算
			"unfreeze_time":     time.Now(),
		}).Error; err != nil {
			return err
		}

		// 6. 记录资金变动日志
		balanceLog := &share.ShareBalanceLog{
			UserID:        income.UserID,
			Type:          1, // 收入
			Amount:        income.CommissionAmount,
			BeforeBalance: user.AvailableBalance,
			AfterBalance:  newBalance,
			SourceType:    1, // 分享收益
			SourceID:      &income.ID,
			Remark:        "分享订单收益解冻",
		}
		if err := tx.Create(balanceLog).Error; err != nil {
			return err
		}

		return nil
	})
}

// CheckIncomeProcess 检查收益处理状态
func (s *IncomeService) CheckIncomeProcess() error {
	var incomes []share.ShareIncome

	// 查找所有已结算的记录
	err := s.db.Where("commission_status = 1").Find(&incomes).Error
	if err != nil {
		return err
	}

	for _, income := range incomes {
		// 检查是否有对应的资金变动记录
		var count int64
		err := s.db.Model(&share.ShareBalanceLog{}).
			Where("source_type = 1 AND source_id = ?", income.ID).
			Count(&count).Error
		if err != nil {
			continue
		}

		if count > 0 {
			// 已有资金变动记录，标记为已结算
			s.db.Model(&income).Updates(map[string]interface{}{
				"commission_status": 1,
				"process_time":      time.Now(),
				"remark":            "系统自动标记-已结算",
			})
		}
	}

	return nil
}

// GetPage 分页获取分享收益列表
func (s *IncomeService) GetIncomList(req *dto.ShareIncomeListReq) (list []distributionvo.IncomeVO, total int64, err error) {
	var incomes []share.ShareIncome
	db := s.db.Model(&share.ShareIncome{})

	// 关联用户表(推荐官和被邀请人)
	db = db.Select(`
		share_income.*,
		referrer.nickname as user_name,
		invitee.nickname as invitee_name,
		share_income.source_order_no as order_no
	`).
		Joins("LEFT JOIN user as referrer ON referrer.id = share_income.user_id").
		Joins("LEFT JOIN user as invitee ON invitee.id = share_income.invitee_id")

	// 构建查询条件
	if req.UserId != "" {
		db = db.Where("share_income.user_id = ?", req.UserId)
	}
	if req.Nickname != "" {
		db = db.Where("referrer.nickname LIKE ? OR invitee.nickname LIKE ?", "%"+req.Nickname+"%", "%"+req.Nickname+"%")
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
			UserName:         income.UserName,        // 推荐官名称
			InviteeId:        income.InviteeID,       // 被邀请人ID
			InviteeName:      income.InviteeNickname, // 被邀请人昵称(从user表读取)
			OrderId:          0,                      // share_income表中没有order_id字段
			OrderNo:          income.SourceOrderNo,
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
