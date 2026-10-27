package shop

import "time"

// ShopOrderReceiver 订单地址表
type ShopOrderReceiver struct {
	Id               int64     `json:"id" gorm:"primaryKey"`
	OrderId          int64     `json:"orderId" gorm:"not null"`
	Receiver         string    `json:"receiver" gorm:"not null"`
	Phone            string    `json:"phone" gorm:"not null"`
	Province         string    `json:"province" gorm:"not null"`
	City             string    `json:"city" gorm:"not null"`
	District         string    `json:"district" gorm:"not null"`
	Detail           string    `json:"detail" gorm:"not null"`
	LogisticsCompany string    `json:"logisticsCompany"`
	LogisticsNo      string    `json:"logisticsNo"`
	IsDelete         int       `json:"isDelete" gorm:"not null;default:0"`
	CreateTime       time.Time `json:"createTime" gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdateTime       time.Time `json:"updateTime" gorm:"not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
}
