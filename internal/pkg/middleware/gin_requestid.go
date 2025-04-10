package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport"
)

// GinRequestID 中间件用于检查请求头中是否有requestid，如果没有则自动生成一个
func GinRequestID(logger log.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 尝试从请求头中获取requestID
		requestID := c.GetHeader(RequestIDKey)

		// 如果没有requestID，则生成一个新的
		if requestID == "" {
			requestID = generateRequestID()
		}

		// 将requestID添加到响应头中
		c.Header(RequestIDKey, requestID)

		// 将requestID保存到上下文中，便于后续使用
		ctx := context.WithValue(c.Request.Context(), RequestIDKey, requestID)
		c.Request = c.Request.WithContext(ctx)

		// 将请求ID添加到 Kratos 的响应头中
		if tr, ok := transport.FromServerContext(ctx); ok {
			tr.ReplyHeader().Set(RequestIDKey, requestID)
		}

		// 继续处理请求
		c.Next()
	}
}
