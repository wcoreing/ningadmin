package utils

import (
	"fmt"
	"strings"
	"wnadmin/app/config"
)

// RemoveDomainFromURL 从URL中移除域名部分，只保留相对路径
func RemoveDomainFromURL(url string) string {
	if url == "" {
		return ""
	}

	// 如果已经是相对路径，直接返回
	if strings.HasPrefix(url, "/") {
		return url
	}

	// 移除各种可能的域名前缀
	prefixes := []string{
		"http://localhost:8080",
		"https://localhost:8080",
		"http://127.0.0.1:8080",
		"https://127.0.0.1:8080",
	}

	// 如果配置了文件域名，也要移除
	if config.App.FileDomain != "" {
		prefixes = append(prefixes, config.App.FileDomain)
	}

	for _, prefix := range prefixes {
		if strings.HasPrefix(url, prefix) {
			return strings.TrimPrefix(url, prefix)
		}
	}

	// 如果没有匹配的前缀，检查是否是完整URL
	if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
		// 提取路径部分
		parts := strings.Split(url, "/")
		if len(parts) >= 4 {
			// 重新组合路径部分 (跳过 http:, "", domain:port)
			return "/" + strings.Join(parts[3:], "/")
		}
	}

	return url
}

// AddDomainToURL 为相对路径添加域名前缀，支持智能域名选择
func AddDomainToURL(url string) string {
	return BuildFileURL(url, "")
}

// BuildFileURL 根据配置文件的上传方式构建文件访问URL
func BuildFileURL(relativePath, _ string) string {
	if relativePath == "" {
		return ""
	}

	// 如果是完整URL（历史数据），直接返回
	if strings.HasPrefix(relativePath, "http://") || strings.HasPrefix(relativePath, "https://") {
		return relativePath
	}

	// 确保相对路径以/开头
	if !strings.HasPrefix(relativePath, "/") {
		relativePath = "/" + relativePath
	}

	// 根据配置文件的上传方式决定使用哪个域名
	if strings.ToLower(config.App.Upload.Type) == "oss" {
		// OSS配置，使用OSS域名
		if config.App.Aliyun.OSS.Domain != "" {
			return config.App.Aliyun.OSS.Domain + relativePath
		}
		// 如果没有配置OSS域名，使用默认OSS域名
		objectKey := strings.TrimPrefix(relativePath, "/")
		return fmt.Sprintf("https://%s.%s/%s", config.App.Aliyun.OSS.BucketName, config.App.Aliyun.OSS.Endpoint, objectKey)
	} else {
		// 本地配置，使用本地域名
		if config.App.FileDomain != "" {
			return config.App.FileDomain + relativePath
		}
		return relativePath
	}
}

// ProcessImageURLsForSave 处理图片URL列表，用于保存时去除域名
func ProcessImageURLsForSave(urls []string) []string {
	result := make([]string, len(urls))
	for i, url := range urls {
		result[i] = RemoveDomainFromURL(url)
	}
	return result
}

// ProcessImageURLsForResponse 处理图片URL列表，用于响应时添加域名
func ProcessImageURLsForResponse(urls []string) []string {
	result := make([]string, len(urls))
	for i, url := range urls {
		result[i] = AddDomainToURL(url)
	}
	return result
}
