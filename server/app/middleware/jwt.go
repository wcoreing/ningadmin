package middleware

import (
	"strings"
	"wnadmin/app/common/response"
	"wnadmin/app/common/utils"

	"github.com/gin-gonic/gin"
)

// JWT 中间件，用于验证令牌
func JWT() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 从请求头中获取令牌
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			response.FailWithMessage(401, "未提供令牌", ctx)
			ctx.Abort()
			return
		}

		// 检查令牌格式
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			response.FailWithMessage(401, "令牌格式错误", ctx)
			ctx.Abort()
			return
		}

		// 解析令牌
		claims, err := utils.ParseToken(parts[1])
		if err != nil {
			response.FailWithMessage(401, "无效的令牌", ctx)
			ctx.Abort()
			return
		}

		// 将用户信息存储到上下文中
		ctx.Set("userId", claims.UserID)
		ctx.Set("userName", claims.UserName)

		ctx.Next()
	}
}
