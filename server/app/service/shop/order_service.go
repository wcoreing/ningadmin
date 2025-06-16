package shop

import (
	"errors"
	"time"
	"wnadmin/app/model/do/shop"
	"wnadmin/app/model/dto"
	"wnadmin/app/provider/db"
	"wnadmin/app/service/distribution"

	"gorm.io/gorm"
)

type OrderService struct{}

const timeFormat = "2006-01-02 15:04:05"

// formatTime 格式化时间
func formatTime(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format(timeFormat)
}

// GetOrderList 获取订单列表
func (s *OrderService) GetOrderList(req *dto.OrderListReq) (*dto.PageResult, error) {
	var (
		list  []shop.ShopOrder
		total int64
	)

	// 构建查询
	db := (&shop.ShopOrder{}).Query()

	// 添加查询条件
	if req.OrderNo != "" {
		db = db.Where("order_no LIKE ?", "%"+req.OrderNo+"%")
	}
	if req.Status != "" {
		db = db.Where("status = ?", req.Status)
	}
	if req.DiningType != "" {
		db = db.Where("dining_type = ?", req.DiningType)
	}
	if !req.StartTime.IsZero() {
		db = db.Where("create_time >= ?", req.StartTime)
	}
	if !req.EndTime.IsZero() {
		db = db.Where("create_time <= ?", req.EndTime)
	}

	// 查询总数
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}

	// 查询列表
	if err := db.Preload("Items").Preload("Receiver").
		Offset((req.PageNum - 1) * req.PageSize).
		Limit(req.PageSize).
		Order("create_time DESC").
		Find(&list).Error; err != nil {
		return nil, err
	}

	// 转换响应
	respList := make([]dto.OrderListResp, 0)
	for _, order := range list {
		items := make([]dto.OrderItemResp, 0)
		for _, item := range order.Items {
			items = append(items, dto.OrderItemResp{
				GoodsId:    item.GoodsId,
				GoodsName:  item.GoodsName,
				GoodsImage: item.GoodsImage,
				SpecText:   item.Specs,
				Price:      item.Price,
				Quantity:   item.Quantity,
			})
		}

		var address string
		if order.Receiver != nil {
			address = order.Receiver.Province + order.Receiver.City + order.Receiver.District + order.Receiver.Detail
		}

		respList = append(respList, dto.OrderListResp{
			Id:          order.Id,
			OrderNo:     order.OrderNo,
			TotalAmount: order.TotalAmount,
			PayAmount:   order.PayAmount,
			Status:      order.Status,
			DiningType:  order.DiningType,
			TableId:     order.TableId,
			TableNo:     order.TableNo,
			PayType:     order.PayType,
			PayTime:     formatTime(order.PayTime),
			CreateTime:  order.CreateTime.Format(timeFormat),
			ShipTime:    formatTime(order.ShipTime),
			FinishTime:  formatTime(order.FinishTime),
			Receiver: dto.ReceiverInfo{
				Receiver:         order.Receiver.Receiver,
				Phone:            order.Receiver.Phone,
				Address:          address,
				LogisticsCompany: order.Receiver.LogisticsCompany,
				LogisticsNo:      order.Receiver.LogisticsNo,
			},
			Items: items,
		})
	}

	return &dto.PageResult{
		List:     respList,
		Total:    total,
		PageNum:  req.PageNum,
		PageSize: req.PageSize,
	}, nil
}

// ShipOrder 订单发货
func (s *OrderService) ShipOrder(req *dto.OrderShipReq) error {
	// 开启事务
	return db.WnAdmin().Transaction(func(tx *gorm.DB) error {
		// 查询订单
		order := &shop.ShopOrder{}
		if err := tx.Where("id = ? AND is_delete = 0", req.OrderId).First(order).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return errors.New("订单不存在")
			}
			return err
		}

		// 检查订单状态
		if order.Status != "unshipped" {
			return errors.New("订单状态不正确，无法发货")
		}

		// 更新订单状态和发货信息
		updates := map[string]interface{}{
			"status":    "delivered",
			"ship_time": time.Now(),
		}
		if err := tx.Model(order).Updates(updates).Error; err != nil {
			return err
		}

		// 更新收货地址表的物流信息
		receiver := &shop.ShopOrderReceiver{}
		if err := tx.Where("order_id = ? AND is_delete = 0", req.OrderId).First(receiver).Error; err != nil {
			return err
		}

		if err := tx.Model(receiver).Updates(map[string]interface{}{
			"logistics_company": req.LogisticsCompany,
			"logistics_no":      req.LogisticsNo,
		}).Error; err != nil {
			return err
		}

		return nil
	})
}

// CompleteOrder 完成订单
func (s *OrderService) CompleteOrder(orderID uint) error {
	return db.WnAdmin().Transaction(func(tx *gorm.DB) error {
		// 1. 更新订单状态
		if err := tx.Model(&shop.ShopOrder{}).
			Where("id = ?", orderID).
			Updates(map[string]interface{}{
				"status":      "completed",
				"finish_time": time.Now(),
			}).Error; err != nil {
			return err
		}

		// 2. 处理分销收益
		orderListener := distribution.NewOrderListenerService()
		if err := orderListener.HandleOrderComplete(orderID); err != nil {
			return err
		}

		return nil
	})
}
