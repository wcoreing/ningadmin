package response

import (
	"net/http"
	"wnadmin/app/provider/mylog"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Response struct {
	Code      int         `json:"code"`
	Data      interface{} `json:"data"`
	Msg       string      `json:"msg"`
	RequestID string      `json:"request_id,omitempty"`
}

const (
	ERROR   = 9999
	SUCCESS = 200
)

func Result(Code int, data interface{}, Msg string, c *gin.Context) {
	// 获取请求ID
	requestID, exists := c.Get("request_id")
	requestIDStr := ""
	if exists {
		requestIDStr = requestID.(string)
	}

	// 使用带请求ID的日志记录
	var logger *logrus.Entry
	if requestIDStr != "" {
		logger = mylog.WithRequestID(requestIDStr)
	} else {
		logger = mylog.Log.WithFields(logrus.Fields{})
	}

	// 记录响应信息
	logger.WithField("status_code", Code).
		WithField("message", Msg).
		Info("响应结果")

	// 返回响应
	c.JSON(http.StatusOK, Response{
		Code:      Code,
		Data:      data,
		Msg:       Msg,
		RequestID: requestIDStr,
	})
}

func Ok(c *gin.Context) {
	Result(SUCCESS, map[string]interface{}{}, "success", c)
}

func OkWithMessage(message string, c *gin.Context) {
	Result(SUCCESS, map[string]interface{}{}, message, c)
}

func OkWithData(data interface{}, c *gin.Context) {
	Result(SUCCESS, data, "success", c)
}

func OkWithDetailed(data interface{}, message string, c *gin.Context) {
	Result(SUCCESS, data, message, c)
}

func Fail(c *gin.Context) {
	Result(ERROR, map[string]interface{}{}, "fail", c)
}

func FailWithMessage(code int, message string, c *gin.Context) {
	Result(code, map[string]interface{}{}, message, c)
}

func FailWithDetailed(data interface{}, message string, c *gin.Context) {
	Result(ERROR, data, message, c)
}
