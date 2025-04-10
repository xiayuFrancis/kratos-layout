package server

import (
	connectv1 "kratosdemo/api/connect/v1"
	"kratosdemo/internal/conf"
	"kratosdemo/internal/pkg/middleware"
	"kratosdemo/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	kratoshttp "github.com/go-kratos/kratos/v2/transport/http"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, connect *service.ConnectService, logger log.Logger) *kratoshttp.Server {
	// 创建 Gin 引擎
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	// 添加中间件
	router.Use(gin.Recovery())
	router.Use(middleware.GinRequestID(logger))

	// 注册路由
	router.GET("/v1/connect/test", func(c *gin.Context) {
		req := &connectv1.TestConnectRequest{}
		resp, err := connect.TestConnect(c.Request.Context(), req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": 500,
				"msg":  err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, resp)
	})

	// 创建 Kratos HTTP 服务器
	var opts = []kratoshttp.ServerOption{
		kratoshttp.Middleware(
			recovery.Recovery(),
			middleware.RequestID(logger),
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, kratoshttp.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, kratoshttp.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, kratoshttp.Timeout(c.Http.Timeout.AsDuration()))
	}

	// 使用 Gin 作为 HTTP 处理器
	srv := kratoshttp.NewServer(opts...)
	srv.HandlePrefix("/", router)

	return srv
}
