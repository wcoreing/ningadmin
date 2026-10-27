package shop

import (
	"wnadmin/app/model/do/shop"
	"wnadmin/app/model/dto"
	"wnadmin/app/provider/db"
	"wnadmin/app/service/shop/logistics"
)

type LogisticsService struct {
	providerType logistics.ProviderType
}

func NewLogisticsService(providerType logistics.ProviderType) *LogisticsService {
	return &LogisticsService{
		providerType: providerType,
	}
}

// GetLogisticsTrack 获取物流轨迹
func (s *LogisticsService) GetLogisticsTrack(req *dto.LogisticsTrackReq) (*dto.LogisticsTrackResp, error) {
	// 1. 查询订单信息
	order := &shop.ShopOrder{}
	if err := db.WnAdmin().Preload("Receiver").Where("id = ?", req.OrderId).First(order).Error; err != nil {
		return nil, err
	}

	// 2. 获取物流提供者并查询
	provider := logistics.GetFactory().Get(s.providerType)
	trackInfo, err := provider.GetTrack(
		order.Receiver.LogisticsCompany,
		order.Receiver.LogisticsNo,
	)
	if err != nil {
		return nil, err
	}

	// 3. 封装响应
	return &dto.LogisticsTrackResp{
		OrderNo:          order.OrderNo,
		LogisticsCompany: getExpressCompanyName(order.Receiver.LogisticsCompany),
		LogisticsNo:      order.Receiver.LogisticsNo,
		Status:           trackInfo.Status,
		Traces:           trackInfo.Traces,
	}, nil
}

// 物流公司映射
var expressCompanyMap = map[string]string{
	"SF":   "顺丰速运",
	"ZTO":  "中通快递",
	"YTO":  "圆通速递",
	"YD":   "韵达快递",
	"STO":  "申通快递",
	"EMS":  "EMS",
	"HTKY": "百世快递",
	"JD":   "京东物流",
}

func getExpressCompanyName(code string) string {
	if name, ok := expressCompanyMap[code]; ok {
		return name
	}
	return code
}
