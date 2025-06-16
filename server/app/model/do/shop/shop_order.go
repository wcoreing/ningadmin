package shop

import (
	"time"
	"wnadmin/app/provider/db"

	"gorm.io/gorm"
)

// ShopOrder 订单表
type ShopOrder struct {
	Id              int64              `json:"id" gorm:"primaryKey"`
	OrderNo         string             `json:"orderNo" gorm:"not null"`
	UserId          int64              `json:"userId" gorm:"not null"`
	TotalAmount     float64            `json:"totalAmount" gorm:"not null"`
	PayAmount       float64            `json:"payAmount" gorm:"not null"`
	FreightAmount   float64            `json:"freightAmount" gorm:"not null"`
	DiscountAmount  float64            `json:"discountAmount" gorm:"not null"`
	Status          string             `json:"status" gorm:"not null;default:'unpaid'"`
	DiningType      string             `json:"diningType" gorm:"default:'takeout'"`
	TableId         *int64             `json:"tableId"`
	TableNo         string             `json:"tableNo"`
	PayType         int                `json:"payType"`
	PayTime         *time.Time         `json:"payTime"`
	ReceiverName    string             `json:"receiverName" gorm:"not null"`
	ReceiverPhone   string             `json:"receiverPhone" gorm:"not null"`
	ReceiverAddress string             `json:"receiverAddress" gorm:"not null"`
	ExpressCompany  string             `json:"expressCompany"`
	ExpressNo       string             `json:"expressNo"`
	ShipTime        *time.Time         `json:"shipTime"`
	ReceiveTime     *time.Time         `json:"receiveTime"`
	FinishTime      *time.Time         `json:"finishTime"`
	CancelTime      *time.Time         `json:"cancelTime"`
	CancelReason    string             `json:"cancelReason"`
	Remark          string             `json:"remark"`
	IsDelete        int                `json:"isDelete" gorm:"not null;default:0"`
	CreateTime      time.Time          `json:"createTime" gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdateTime      time.Time          `json:"updateTime" gorm:"not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
	Items           []ShopOrderGoods   `json:"items" gorm:"foreignKey:OrderId"`
	Receiver        *ShopOrderReceiver `json:"receiver" gorm:"foreignKey:OrderId"`
}

// Query 获取查询构建器
func (m *ShopOrder) Query() *gorm.DB {
	return db.WnAdmin().Model(m).Where("is_delete = ?", 0)
}
