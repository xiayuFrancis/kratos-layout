package logger

import (
	"kratosdemo/internal/conf"
	"os"

	"github.com/go-kratos/kratos/v2/log"
)

// NewLogger 创建一个新的日志记录器
func NewLogger(c *conf.Logger) log.Logger {
	if c == nil {
		// 如果没有配置，使用默认的标准输出日志记录器
		return log.NewStdLogger(os.Stdout)
	}
	
	// 使用 Zap 日志记录器
	return NewZapLogger(c)
}
