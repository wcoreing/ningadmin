package logistics

import "wnadmin/app/model/dto"

// Provider 物流查询服务提供者接口
type Provider interface {
	// GetTrack 获取物流轨迹
	GetTrack(company, number string) (*LogisticsTrackInfo, error)
}

// LogisticsTrackInfo 物流轨迹信息
type LogisticsTrackInfo struct {
	Status string               // 物流状态
	Traces []dto.LogisticsTrace // 物流轨迹
}
