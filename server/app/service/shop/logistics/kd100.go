package logistics

import (
	"wnadmin/app/model/dto"
)

type Kd100Provider struct {
	Key    string
	Secret string
	ApiUrl string
}

func NewKd100Provider(key, secret string) *Kd100Provider {
	return &Kd100Provider{
		Key:    key,
		Secret: secret,
		ApiUrl: "https://poll.kuaidi100.com/poll/query.do",
	}
}

func (p *Kd100Provider) GetTrack(company, number string) (*LogisticsTrackInfo, error) {
	// TODO: 实现快递100的查询逻辑
	return &LogisticsTrackInfo{
		Status: "在途",
		Traces: []dto.LogisticsTrace{
			// ... 轨迹信息
		},
	}, nil
}
