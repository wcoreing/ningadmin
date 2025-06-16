package config

type AppConfig struct {
	ENV        string       `yaml:"ENV"`
	Host       string       `yaml:"Host"`
	Port       int          `yaml:"Port"`
	Name       string       `yaml:"Name"`
	FileDomain string       `yaml:"FileDomain"` // 文件访问域名
	Upload     UploadConfig `yaml:"Upload"`
	Aliyun     AliyunConfig `yaml:"aliyun"` // 阿里云配置
}

// UploadConfig 文件上传配置
type UploadConfig struct {
	Type         string   `yaml:"Type"`         // 上传方式: local 或 oss
	MaxSize      int      `yaml:"MaxSize"`      // 最大文件大小(MB)
	AllowedTypes []string `yaml:"AllowedTypes"` // 允许的文件类型
}

// AliyunConfig 阿里云配置
type AliyunConfig struct {
	OSS OSSConfig `yaml:"oss"` // OSS配置
	SMS SMSConfig `yaml:"sms"` // 短信配置
}

// OSSConfig 阿里云OSS配置
type OSSConfig struct {
	Endpoint        string `yaml:"endpoint"`        // OSS访问域名
	AccessKeyId     string `yaml:"accessKeyId"`     // AccessKey ID
	AccessKeySecret string `yaml:"accessKeySecret"` // AccessKey Secret
	BucketName      string `yaml:"bucketName"`      // 存储空间名称
	Domain          string `yaml:"domain"`          // OSS访问域名（CDN或默认域名）
	BasePath        string `yaml:"basePath"`        // 文件存储基础路径
}

// SMSConfig 阿里云短信配置
type SMSConfig struct {
	AccessKeyId     string                `yaml:"accessKeyId"`     // AccessKey ID
	AccessKeySecret string                `yaml:"accessKeySecret"` // AccessKey Secret
	SignName        string                `yaml:"signName"`        // 签名名称
	TemplateCode    SMSTemplateCodeConfig `yaml:"templateCode"`    // 模板配置
}

// SMSTemplateCodeConfig 短信模板配置
type SMSTemplateCodeConfig struct {
	Register      string `yaml:"register"`      // 注册模板
	Login         string `yaml:"login"`         // 登录模板
	ResetPassword string `yaml:"resetPassword"` // 重置密码模板
}

var App AppConfig
