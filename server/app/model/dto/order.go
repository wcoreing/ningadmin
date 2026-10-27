package dto

import "time"

// OrderListReq 订单列表请求参数
type OrderListReq struct {
	OrderNo    string    `form:"orderNo"`    // 订单编号
	Status     string    `form:"status"`     // 订单状态
	DiningType string    `form:"diningType"` // 用餐类型
	StartTime  time.Time `form:"startTime"`  // 开始时间
	EndTime    time.Time `form:"endTime"`    // 结束时间
	PageNum    int       `form:"pageNum"`    // 页码
	PageSize   int       `form:"pageSize"`   // 每页数量
}

// OrderListResp 订单列表响应
type OrderListResp struct {
	Id          int64           `json:"id"`          // 订单ID
	OrderNo     string          `json:"orderNo"`     // 订单编号
	TotalAmount float64         `json:"totalAmount"` // 订单总金额
	PayAmount   float64         `json:"payAmount"`   // 实付金额
	Status      string          `json:"status"`      // 订单状态
	DiningType  string          `json:"diningType"`  // 用餐类型
	TableId     *int64          `json:"tableId"`     // 餐桌ID
	TableNo     string          `json:"tableNo"`     // 餐桌编号
	PayType     int             `json:"payType"`     // 支付方式
	PayTime     string          `json:"payTime"`     // 支付时间
	CreateTime  string          `json:"createTime"`  // 创建时间
	ShipTime    string          `json:"shipTime"`    // 发货时间
	FinishTime  string          `json:"finishTime"`  // 完成时间
	Receiver    ReceiverInfo    `json:"receiver"`    // 收货人信息
	Items       []OrderItemResp `json:"items"`       // 订单商品
}

// OrderItemResp 订单商品响应
type OrderItemResp struct {
	GoodsId    int64   `json:"goodsId"`    // 商品ID
	GoodsName  string  `json:"goodsName"`  // 商品名称
	GoodsImage string  `json:"goodsImage"` // 商品图片
	SpecText   string  `json:"specText"`   // 规格文本
	Price      float64 `json:"price"`      // 商品单价
	Quantity   int     `json:"quantity"`   // 购买数量
}

// ReceiverInfo 收货人信息
type ReceiverInfo struct {
	Receiver         string `json:"receiver"`         // 收货人
	Phone            string `json:"phone"`            // 联系电话
	Address          string `json:"address"`          // 完整地址
	LogisticsCompany string `json:"logisticsCompany"` // 物流公司
	LogisticsNo      string `json:"logisticsNo"`      // 物流单号
}

// OrderShipReq 订单发货请求
type OrderShipReq struct {
	OrderId          int64  `json:"orderId" binding:"required"`          // 订单ID
	LogisticsCompany string `json:"logisticsCompany" binding:"required"` // 物流公司
	LogisticsNo      string `json:"logisticsNo" binding:"required"`      // 物流单号
	Remark           string `json:"remark"`                              // 发货备注
}

// LogisticsTrackReq 物流轨迹查询请求
type LogisticsTrackReq struct {
	OrderId int64 `json:"orderId" binding:"required"` // 订单ID
}

// LogisticsTrackResp 物流轨迹响应
type LogisticsTrackResp struct {
	OrderNo          string           `json:"orderNo"`          // 订单编号
	LogisticsCompany string           `json:"logisticsCompany"` // 物流公司
	LogisticsNo      string           `json:"logisticsNo"`      // 物流单号
	Status           string           `json:"status"`           // 物流状态
	Traces           []LogisticsTrace `json:"traces"`           // 物流轨迹
}

// LogisticsTrace 物流轨迹
type LogisticsTrace struct {
	AcceptTime    string `json:"acceptTime"`    // 接收时间
	AcceptStation string `json:"acceptStation"` // 当前位置
	Remark        string `json:"remark"`        // 备注信息
}
