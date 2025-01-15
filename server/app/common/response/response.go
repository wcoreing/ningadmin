package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"wnadmin/app/provider/mylog"
)

type Response struct {
	ReturnCode int         `json:"returnCode"`
	Data       interface{} `json:"data"`
	ReturnMsg  string      `json:"returnMsg"`
}

const (
	ERROR   = 9999
	SUCCESS = 0000
)

func Result(returnCode int, data interface{}, returnMsg string, c *gin.Context) {
	// 开始时间
	mylog.Log.Info("response", returnCode, data, returnMsg)
	c.JSON(http.StatusOK, Response{
		returnCode,
		data,
		returnMsg,
	})
	return
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
