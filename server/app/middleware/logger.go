package middleware

import (
	"bytes"
	"fmt"
	"io"
	"time"
	"wnadmin/app/provider/mylog"

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
	}
}
