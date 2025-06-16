# 文件URL构建优化说明

## 🎯 优化目标

将文件URL的构建逻辑统一为**基于配置文件的上传方式**来决定域名选择，而不是依赖数据库中的存储类型字段。

## 📝 核心逻辑

### 新的判断逻辑

```yaml
# 配置文件 app.yaml
Upload:
  Type: oss  # 或 local
```

```go
// 根据配置文件决定域名
if config.App.Upload.Type == "oss" {
    // 使用OSS域名
    return config.App.Aliyun.OSS.Domain + relativePath
} else {
    // 使用本地域名  
    return config.App.FileDomain + relativePath
}
```

## 🔧 实现方式

### 1. 统一的工具方法

```go
// BuildFileURL 根据配置文件的上传方式构建文件访问URL
func BuildFileURL(relativePath, _ string) string {
    // 如果是完整URL（历史数据），直接返回
    if strings.HasPrefix(relativePath, "http://") || strings.HasPrefix(relativePath, "https://") {
        return relativePath
    }

    // 根据配置文件的Upload.Type决定域名
    if strings.ToLower(config.App.Upload.Type) == "oss" {
        return config.App.Aliyun.OSS.Domain + relativePath
    } else {
        return config.App.FileDomain + relativePath
    }
}
```

### 2. 简化的调用方式

```go
// 所有地方都可以简单调用
finalURL := utils.BuildFileURL(file.Url, "")
// 或者
finalURL := utils.AddDomainToURL(file.Url)
```

## 📊 使用场景

### 场景1：配置为本地存储
```yaml
Upload:
  Type: local
FileDomain: http://localhost:8080
```

**结果：**
- 数据库存储：`/uploads/file.jpg`
- 访问URL：`http://localhost:8080/uploads/file.jpg`

### 场景2：配置为OSS存储
```yaml
Upload:
  Type: oss
aliyun:
  oss:
    domain: "https://cdn.example.com"
```

**结果：**
- 数据库存储：`/uploads/2024/01/15/file.jpg`
- 访问URL：`https://cdn.example.com/uploads/2024/01/15/file.jpg`

### 场景3：历史完整URL数据
**结果：**
- 数据库存储：`https://old-domain.com/file.jpg`
- 访问URL：`https://old-domain.com/file.jpg`（直接返回）

## 🔄 配置切换效果

### 从本地切换到OSS
1. 修改配置：`Upload.Type: local` → `Upload.Type: oss`
2. 重启服务
3. **所有文件**（包括历史文件）立即使用OSS域名访问

### 从OSS切换到本地
1. 修改配置：`Upload.Type: oss` → `Upload.Type: local`
2. 重启服务
3. **所有文件**（包括历史文件）立即使用本地域名访问

## ✅ 优势

1. **配置驱动**：完全由配置文件控制，无需关心数据库字段
2. **简单明了**：逻辑清晰，容易理解和维护
3. **灵活切换**：修改配置即可全局切换域名
4. **向后兼容**：历史完整URL数据保持不变
5. **统一管理**：所有文件URL使用同一套规则

## 🎉 实际效果

现在系统中所有文件的URL构建都遵循这个简单规则：

- **当前配置是OSS** → 所有相对路径文件使用OSS域名
- **当前配置是本地** → 所有相对路径文件使用本地域名
- **历史完整URL** → 保持原样不变

这样既保证了一致性，又提供了灵活性！ 