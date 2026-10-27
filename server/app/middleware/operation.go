package middleware

import (
	"wnadmin/app/model/dto"

	"github.com/gin-gonic/gin"
)

// WithOperation 为路由添加操作信息的中间件
func WithOperation(module, name, opType string) gin.HandlerFunc {
	return func(c *gin.Context) {
		operationInfo := dto.OperationInfo{
			Module: module,
			Name:   name,
			Type:   opType,
		}
		c.Set("operation_info", operationInfo)
		c.Next()
	}
}

// WithSelect 为SELECT操作添加操作信息的中间件
func WithSelect(module, name string) gin.HandlerFunc {
	return WithOperation(module, name, "SELECT")
}

// WithCreate 为CREATE操作添加操作信息的中间件
func WithCreate(module, name string) gin.HandlerFunc {
	return WithOperation(module, name, "CREATE")
}

// WithUpdate 为UPDATE操作添加操作信息的中间件
func WithUpdate(module, name string) gin.HandlerFunc {
	return WithOperation(module, name, "UPDATE")
}

// WithDelete 为DELETE操作添加操作信息的中间件
func WithDelete(module, name string) gin.HandlerFunc {
	return WithOperation(module, name, "DELETE")
}
