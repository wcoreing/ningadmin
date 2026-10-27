package utils

import (
	"fmt"
	"mime/multipart"
	"path"
	"strings"
	"time"
	"wnadmin/app/config"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

// OSSUploader OSS上传器
type OSSUploader struct {
	client *oss.Client
	bucket *oss.Bucket
}

// NewOSSUploader 创建OSS上传器
func NewOSSUploader() (*OSSUploader, error) {
	client, err := oss.New(config.App.Aliyun.OSS.Endpoint, config.App.Aliyun.OSS.AccessKeyId, config.App.Aliyun.OSS.AccessKeySecret)
	if err != nil {
		return nil, fmt.Errorf("创建OSS客户端失败: %v", err)
	}

	bucket, err := client.Bucket(config.App.Aliyun.OSS.BucketName)
	if err != nil {
		return nil, fmt.Errorf("获取OSS存储桶失败: %v", err)
	}

	return &OSSUploader{
		client: client,
		bucket: bucket,
	}, nil
}

// UploadFile 上传文件到OSS，返回相对路径
func (u *OSSUploader) UploadFile(file *multipart.FileHeader) (string, error) {
	// 生成文件名
	ext := path.Ext(file.Filename)
	fileName := GenerateUUID() + ext

	// 构建OSS对象键
	objectKey := config.App.Aliyun.OSS.BasePath + "/" + time.Now().Format("2006/01/02") + "/" + fileName

	// 打开文件
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("打开文件失败: %v", err)
	}
	defer src.Close()

	// 上传到OSS
	err = u.bucket.PutObject(objectKey, src)
	if err != nil {
		return "", fmt.Errorf("上传到OSS失败: %v", err)
	}

	// 返回相对路径（OSS对象键）
	return "/" + objectKey, nil
}

// DeleteFile 从OSS删除文件
func (u *OSSUploader) DeleteFile(relativePath string) error {
	// 去掉开头的斜杠，转换为OSS对象键
	objectKey := strings.TrimPrefix(relativePath, "/")

	err := u.bucket.DeleteObject(objectKey)
	if err != nil {
		return fmt.Errorf("从OSS删除文件失败: %v", err)
	}
	return nil
}

// IsValidFileType 检查文件类型是否允许
func IsValidFileType(filename string) bool {
	ext := strings.ToLower(path.Ext(filename))
	if ext != "" && ext[0] == '.' {
		ext = ext[1:] // 去掉点号
	}

	for _, allowedType := range config.App.Upload.AllowedTypes {
		if ext == strings.ToLower(allowedType) {
			return true
		}
	}
	return false
}

// IsValidFileSize 检查文件大小是否允许
func IsValidFileSize(size int64) bool {
	maxSize := int64(config.App.Upload.MaxSize) * 1024 * 1024 // 转换为字节
	return size <= maxSize
}
