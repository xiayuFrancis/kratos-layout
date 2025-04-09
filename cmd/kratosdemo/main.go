package main

import (
	"context"
	"flag"
	"kratosdemo/internal/conf"
	"kratosdemo/internal/pkg/logger"
	"os"
	"strings"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	_ "go.uber.org/automaxprocs"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name string
	// Version is the version of the compiled software.
	Version string
	// flagconf is the config flag.
	flagconf string

	id, _ = os.Hostname()
)

func init() {
	flag.StringVar(&flagconf, "conf", "../../configs", "config path, eg: -conf config.yaml")
}

func newApp(logger log.Logger, gs *grpc.Server, hs *http.Server) *kratos.App {
	return kratos.New(
		kratos.ID(id),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(
			gs,
			hs,
		),
	)
}

func main() {
	flag.Parse()
	// 创建一个临时的标准日志记录器，用于初始化配置加载
	tempLogger := log.With(log.NewStdLogger(os.Stdout),
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"service.id", id,
		"service.name", Name,
		"service.version", Version,
		"trace.id", tracing.TraceID(),
		"span.id", tracing.SpanID(),
	)
	c := config.New(
		config.WithSource(
			file.NewSource(flagconf),
		),
	)
	defer c.Close()

	if err := c.Load(); err != nil {
		panic(err)
	}

	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		panic(err)
	}

	// 使用配置创建 Zap logger
	zapLogger := logger.NewLogger(bc.Logger)
	// 添加公共字段
	logger := log.With(zapLogger,
		//"trace.id", tracing.TraceID(),
		//"span.id", tracing.SpanID(),
		"requestid", log.Valuer(func(ctx context.Context) interface{} {
			requestID, _ := ctx.Value("x-request-id").(string)
			return requestID
		}),
	)

	// 创建日志目录
	if bc.Logger != nil && bc.Logger.Rotate != nil && bc.Logger.Rotate.Filename != "" {
		logDir := bc.Logger.Rotate.Filename
		if idx := strings.LastIndex(logDir, "/"); idx > 0 {
			logDir = logDir[:idx]
			if _, err := os.Stat(logDir); os.IsNotExist(err) {
				if err := os.MkdirAll(logDir, 0755); err != nil {
					tempLogger.Log(log.LevelError, "msg", "failed to create log directory", "error", err)
				}
			}
		}
	}

	app, cleanup, err := wireApp(bc.Server, bc.Data, logger)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}
}
