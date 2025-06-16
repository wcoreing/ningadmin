package task

import (
	"fmt"
	"time"
	"wnadmin/app/model/do/shop"
	"wnadmin/app/provider/db"
	service "wnadmin/app/service/shop"
)

// AutoCompleteOrder 自动完成订单任务
func AutoCompleteOrder() {
	fmt.Println("自动完成订单任务")
	orderService := &service.OrderService{}

	// 查找已发货超过N天的订单
	var orders []shop.ShopOrder
	autoReceiveDays := 7 // 可以从配置中读取
	deadline := time.Now().AddDate(0, 0, -autoReceiveDays)

	err := db.WnAdmin().Where("status = ? AND ship_time <= ?", "delivered", deadline).Find(&orders).Error
	if err != nil {
		return
	}

	for _, order := range orders {
		orderService.CompleteOrder(uint(order.Id))
	}
}
