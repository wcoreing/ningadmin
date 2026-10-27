package logistics

import (
	"wnadmin/app/model/dto"
)

type AliyunProvider struct {
	AccessKeyId     string
	AccessKeySecret string
	ApiUrl          string
}

func NewAliyunProvider(accessKeyId, accessKeySecret string) *AliyunProvider {
	return &AliyunProvider{
		AccessKeyId:     accessKeyId,
		AccessKeySecret: accessKeySecret,
		ApiUrl:          "https://logistics.aliyuncs.com",
	}
}

func (p *AliyunProvider) GetTrack(company, number string) (*LogisticsTrackInfo, error) {
	// TODO: 实现阿里云的查询逻辑
	return &LogisticsTrackInfo{
		Status: "在途",
		Traces: []dto.LogisticsTrace{
			// ... 轨迹信息
		},
	}, nil
}
