package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
	"time"
	"wnadmin/app/model/dto"
	"wnadmin/app/provider/mylog"
	"wnadmin/app/service/system"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// 响应体写入器
type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// Write 写入响应体并保存副本
func (r responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

// 生成请求ID
func generateRequestID() string {
	return uuid.New().String()
}

// 获取操作信息
func getOperationInfo(c *gin.Context, method, path string) dto.OperationInfo {
	// 优先从上下文中获取操作信息（中间件设置）
	if operationInfo, exists := c.Get("operation_info"); exists {
		if info, ok := operationInfo.(dto.OperationInfo); ok {
			return info
		}
	}

	// 如果上下文中没有，则使用原来的映射表方式
	operationLogService := &system.OperationLogService{}
	operationMapping := operationLogService.GetOperationMapping()

	// 构建查找键
	key := method + ":" + path

	// 直接匹配
	if info, exists := operationMapping[key]; exists {
		return info
	}

	// 模糊匹配（处理带参数的路径）
	for pattern, info := range operationMapping {
		if matchPath(pattern, key) {
			return info
		}
	}

	// 默认返回
	return dto.OperationInfo{
		Module: "未知模块",
		Name:   fmt.Sprintf("%s %s", method, path),
		Type:   getOperationType(method),
	}
}

// 路径匹配（支持参数）
func matchPath(pattern, path string) bool {
	// 将路径中的数字ID替换为通配符进行匹配
	// 例如：DELETE:/system/user/delete/123 -> DELETE:/system/user/delete
	re := regexp.MustCompile(`/\d+$`)
	normalizedPath := re.ReplaceAllString(path, "")
	return pattern == normalizedPath
}

// 根据HTTP方法获取操作类型
func getOperationType(method string) string {
	switch method {
	case "GET":
		return "SELECT"
	case "POST":
		return "CREATE"
	case "PUT", "PATCH":
		return "UPDATE"
	case "DELETE":
		return "DELETE"
	default:
		return "OTHER"
	}
}

// 获取用户信息
func getUserInfo(c *gin.Context) (int64, string, string) {
	// 从JWT或session中获取用户信息
	userID := int64(0)
	username := "anonymous"
	userType := "admin"

	if userIDStr, exists := c.Get("user_id"); exists {
		if uid, ok := userIDStr.(string); ok {
			if id, err := strconv.ParseInt(uid, 10, 64); err == nil {
				userID = id
			}
		}
	}

	if usernameStr, exists := c.Get("userName"); exists {
		if un, ok := usernameStr.(string); ok {
			username = un
		}
	}

	return userID, username, userType
}

// 判断是否需要记录操作日志
func shouldRecordLog(method, path string) bool {
	// 排除OPTIONS请求
	if method == "OPTIONS" {
		return false
	}

	// 排除不需要记录的路径
	excludePaths := []string{
		"/health",
		"/ping",
		"/metrics",
		"/static",
		"/assets",
	}

	for _, excludePath := range excludePaths {
		if strings.Contains(path, excludePath) {
			return false
		}
	}

	// 只记录API接口
	return strings.HasPrefix(path, "/api/") || strings.HasPrefix(path, "/system/") || strings.HasPrefix(path, "/auth/")
}

// Logger 日志中间件，记录请求和响应信息，并生成请求ID
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始计时
		startTime := time.Now()

		// 生成请求ID
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = generateRequestID()
			c.Header("X-Request-ID", requestID)
		}

		// 将请求ID保存到上下文
		c.Set("request_id", requestID)

		// 创建带有请求ID的日志条目
		logger := mylog.WithRequestID(requestID)

		// 获取请求体
		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// 记录请求信息
		logger.WithField("method", c.Request.Method).
			WithField("path", c.Request.URL.Path).
			WithField("client_ip", c.ClientIP()).
			WithField("user_agent", c.Request.UserAgent()).
			Info("请求开始")

		// 记录请求体（仅当为POST/PUT/PATCH请求时）
		if c.Request.Method == "POST" || c.Request.Method == "PUT" || c.Request.Method == "PATCH" {
			if len(requestBody) > 0 {
				// 限制日志大小，避免日志过大
				if len(requestBody) > 1024 {
					logger.WithField("request_body", string(requestBody[:1024])+"...").
						Debug("请求体")
				} else {
					logger.WithField("request_body", string(requestBody)).
						Debug("请求体")
				}
			}
		}

		// 捕获响应体
		responseBodyWriter := &responseBodyWriter{
			ResponseWriter: c.Writer,
			body:           bytes.NewBufferString(""),
		}
		c.Writer = responseBodyWriter

		// 处理请求
		c.Next()

		// 结束计时
		latency := time.Since(startTime)

		// 记录响应信息
		logger.WithField("status", c.Writer.Status()).
			WithField("latency", fmt.Sprintf("%s", latency)).
			WithField("error", c.Errors.String()).
			Info("请求结束")

		// 记录响应体
		if responseBodyWriter.body.Len() > 0 {
			respBody := responseBodyWriter.body.String()
			// 限制日志大小，避免日志过大
			if len(respBody) > 1024 {
				logger.WithField("response_body", respBody[:1024]+"...").
					Debug("响应体")
			} else {
				logger.WithField("response_body", respBody).
					Debug("响应体")
			}
		}

		// 记录操作日志到数据库
		if shouldRecordLog(c.Request.Method, c.Request.URL.Path) {
			go recordOperationLog(c, requestID, requestBody, responseBodyWriter.body.Bytes(), startTime, latency)
		}
	}
}

// 异步记录操作日志
func recordOperationLog(c *gin.Context, traceID string, requestBody, responseBody []byte, startTime time.Time, latency time.Duration) {
	defer func() {
		if r := recover(); r != nil {
			mylog.Log.Error("记录操作日志失败:", r)
		}
	}()

	// 获取用户信息
	userID, username, userType := getUserInfo(c)

	// 获取操作信息
	operationInfo := getOperationInfo(c, c.Request.Method, c.Request.URL.Path)

	// 解析请求参数
	var requestParams interface{}
	if len(requestBody) > 0 {
		var params map[string]interface{}
		if err := json.Unmarshal(requestBody, &params); err == nil {
			requestParams = params
		} else {
			requestParams = string(requestBody)
		}
	} else {
		// GET请求的查询参数
		if c.Request.URL.RawQuery != "" {
			requestParams = c.Request.URL.RawQuery
		}
	}

	// 解析响应数据
	var responseData interface{}
	if len(responseBody) > 0 {
		var resp map[string]interface{}
		if err := json.Unmarshal(responseBody, &resp); err == nil {
			responseData = resp
		} else {
			responseData = string(responseBody)
		}
	}

	// 确定操作状态
	status := "1" // 成功
	errorMsg := ""
	if c.Writer.Status() >= 400 {
		status = "0" // 失败
		if len(c.Errors) > 0 {
			errorMsg = c.Errors.String()
		}
	}

	// 构建日志数据
	logData := &dto.OperationLogData{
		TraceID:       traceID,
		UserID:        userID,
		Username:      username,
		UserType:      userType,
		OperationType: operationInfo.Type,
		Module:        operationInfo.Module,
		OperationName: operationInfo.Name,
		Method:        c.Request.Method,
		RequestURL:    c.Request.URL.String(),
		RequestParams: requestParams,
		ResponseData:  responseData,
		IPAddress:     c.ClientIP(),
		UserAgent:     c.Request.UserAgent(),
		Location:      "", // 可以集成IP地址解析服务
		Status:        status,
		ErrorMsg:      errorMsg,
		CostTime:      latency.Milliseconds(),
	}

	// 记录到数据库
	operationLogService := &system.OperationLogService{}
	if err := operationLogService.RecordLog(logData); err != nil {
		mylog.Log.Error("保存操作日志失败:", err)
	}
}
