package middleware

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/google/uuid"
)

// RequestIDKey 是请求ID在上下文中的键
const RequestIDKey = "x-request-id"

// RequestID 中间件用于检查请求头中是否有requestid，如果没有则自动生成一个
func RequestID(logger log.Logger) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			// 尝试从请求头中获取requestID
			requestID := extractRequestID(ctx)

			// 如果没有requestID，则生成一个新的
			if requestID == "" {
				requestID = generateRequestID()
			}

			// 创建带有requestID的logger
			//logger = log.With(logger, RequestIDKey, requestID)

			// 如果是HTTP请求，将requestID添加到响应头中
			if tr, ok := transport.FromServerContext(ctx); ok {
				tr.ReplyHeader().Set(RequestIDKey, requestID)
			}

			// 将requestID保存到上下文中，便于后续使用
			ctx = context.WithValue(ctx, RequestIDKey, requestID)

			// 继续处理请求
			return handler(ctx, req)
		}
	}
}

// extractRequestID 从上下文中提取请求ID
func extractRequestID(ctx context.Context) string {
	if tr, ok := transport.FromServerContext(ctx); ok {
		// 从HTTP头或gRPC元数据中提取requestID
		requestID := tr.RequestHeader().Get(RequestIDKey)
		if requestID != "" {
			return requestID
		}
	}
	return ""
}

// generateRequestID 生成一个新的请求ID
func generateRequestID() string {
	return uuid.New().String()
}
