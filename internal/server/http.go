package server

import (
	"blog/api/common"
	v1 "blog/api/helloworld/v1"
	userv1 "blog/api/user/v1"
	"blog/internal/conf"
	"blog/internal/service"
	"encoding/json"
	"github.com/go-kratos/kratos/v2/encoding"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/google/uuid"
	"golang.org/x/net/context"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, greeter *service.GreeterService, user *service.UserServiceService, logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			MiddlewareRequestID(),
		),
		http.ResponseEncoder(CustomResponseEncoder()),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}

	srv := http.NewServer(opts...)
	v1.RegisterGreeterHTTPServer(srv, greeter)
	userv1.RegisterUserServiceHTTPServer(srv, user)
	return srv
}

// CustomResponseEncoder 自定义的 ResponseEncoder
func CustomResponseEncoder() http.EncodeResponseFunc {
	return func(w http.ResponseWriter, r *http.Request, v interface{}) error {
		if v == nil {
			return nil
		}
		if resp, ok := v.(*common.Response); ok && len(resp.Data) > 0 {
			// 创建一个新的 map 来保存最终的响应
			result := make(map[string]interface{})
			result["code"] = resp.Code
			result["msg"] = resp.Msg

			// 将 Data 字段解析为 interface{}
			var jsonData interface{}
			if err := json.Unmarshal(resp.Data, &jsonData); err != nil {
				return err
			}
			result["data"] = jsonData

			w.Header().Set("Content-Type", "application/json")
			return json.NewEncoder(w).Encode(result)
		}
		// 对于其他类型的响应，使用默认 JSON 编码器
		codec := encoding.GetCodec("json")
		data, err := codec.Marshal(v)
		if err != nil {
			return err
		}
		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(data)
		return err
	}
}

// MiddlewareRequestID 自定义中间件，用于获取request_id
func MiddlewareRequestID() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			if tr, ok := transport.FromServerContext(ctx); ok {
				if htr, ok := tr.(*http.Transport); ok {
					// 获取http.Request
					httpReq := htr.Request()
					// 从header中获取request_id
					requestID := httpReq.Header.Get("X-Request-ID")
					if requestID == "" {
						// 如果没有request_id，则生成一个新的UUID
						requestID = uuid.New().String()
						// 可以将生成的request_id添加到响应header中，以便客户端可以使用
						httpReq.Header.Set("X-Request-ID", requestID)
					}
					// 可以将request_id存储在上下文中，以便后续处理
					ctx = context.WithValue(ctx, "request_id", requestID)
				}
			}
			// 继续执行下一个中间件或handler
			return handler(ctx, req)
		}
	}
}
