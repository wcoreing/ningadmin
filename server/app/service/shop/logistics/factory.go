package logistics

import "sync"

type ProviderType string

const (
	ProviderKd100  ProviderType = "kd100"
	ProviderAliyun ProviderType = "aliyun"
)

var (
	providerFactory *Factory
	once            sync.Once
)

type Factory struct {
	providers map[ProviderType]Provider
}

// GetFactory 获取工厂实例
func GetFactory() *Factory {
	once.Do(func() {
		providerFactory = &Factory{
			providers: make(map[ProviderType]Provider),
		}
		// 注册默认提供者
		providerFactory.Register(ProviderKd100, NewKd100Provider("your_key", "your_secret"))
		providerFactory.Register(ProviderAliyun, NewAliyunProvider("your_key", "your_secret"))
	})
	return providerFactory
}

// Register 注册物流查询提供者
func (f *Factory) Register(name ProviderType, provider Provider) {
	f.providers[name] = provider
}

// Get 获取物流查询提供者
func (f *Factory) Get(name ProviderType) Provider {
	return f.providers[name]
}
