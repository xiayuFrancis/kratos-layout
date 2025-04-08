package logger

import (
	"kratosdemo/internal/conf"
	"os"
	"strings"

	"github.com/go-kratos/kratos/v2/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// ZapLogger 是基于 zap 的日志记录器
type ZapLogger struct {
	log  *zap.Logger
	Sync func() error
}

// Log 实现 log.Logger 接口
func (l *ZapLogger) Log(level log.Level, keyvals ...interface{}) error {
	if len(keyvals) == 0 {
		return nil
	}
	if len(keyvals)%2 != 0 {
		keyvals = append(keyvals, "KEYVALS UNPAIRED")
	}

	// 提取消息字段
	message := ""
	var data []zap.Field
	for i := 0; i < len(keyvals); i += 2 {
		key := toString(keyvals[i])
		value := keyvals[i+1]
		
		// 如果是消息字段，提取为主消息
		if key == "msg" || key == log.DefaultMessageKey {
			if str, ok := value.(string); ok {
				message = str
				continue // 跳过添加到字段列表
			}
		}
		
		data = append(data, zap.Any(key, value))
	}

	switch level {
	case log.LevelDebug:
		l.log.Debug(message, data...)
	case log.LevelInfo:
		l.log.Info(message, data...)
	case log.LevelWarn:
		l.log.Warn(message, data...)
	case log.LevelError:
		l.log.Error(message, data...)
	case log.LevelFatal:
		l.log.Fatal(message, data...)
	}
	return nil
}

// toString 将任意类型转换为字符串
func toString(v interface{}) string {
	if v == nil {
		return ""
	}
	switch v := v.(type) {
	case string:
		return v
	case error:
		return v.Error()
	case []byte:
		return string(v)
	default:
		return strings.TrimSpace(strings.TrimPrefix(strings.TrimSuffix(log.DefaultMessageKey, "}"), "{"))
	}
}

// NewZapLogger 创建一个新的基于 zap 的日志记录器
func NewZapLogger(c *conf.Logger) log.Logger {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// 设置日志级别
	var level zapcore.Level
	switch strings.ToLower(c.Level) {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	default:
		level = zapcore.InfoLevel
	}

	// 配置日志输出
	var cores []zapcore.Core

	// 配置日志格式
	var encoder zapcore.Encoder
	if c.Format == "json" {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	// 配置标准输出
	if strings.Contains(c.OutputPaths, "stdout") {
		stdoutCore := zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), level)
		cores = append(cores, stdoutCore)
	}

	if strings.Contains(c.ErrorOutputPaths, "stderr") {
		stderrCore := zapcore.NewCore(encoder, zapcore.AddSync(os.Stderr), zapcore.ErrorLevel)
		cores = append(cores, stderrCore)
	}

	// 配置文件输出和日志轮转
	if c.Rotate != nil && c.Rotate.Filename != "" {
		// 配置日志轮转
		rotateLogger := &lumberjack.Logger{
			Filename:   c.Rotate.Filename,
			MaxSize:    int(c.Rotate.MaxSize),    // 单位：MB
			MaxBackups: int(c.Rotate.MaxBackups), // 保留的旧日志文件数量
			MaxAge:     int(c.Rotate.MaxAge),     // 保留天数
			Compress:   c.Rotate.Compress,        // 是否压缩
		}

		fileCore := zapcore.NewCore(encoder, zapcore.AddSync(rotateLogger), level)
		cores = append(cores, fileCore)

		// 错误日志单独输出
		if c.Rotate.Filename != "" {
			errorFilename := strings.Replace(c.Rotate.Filename, ".log", ".error.log", 1)
			errorRotateLogger := &lumberjack.Logger{
				Filename:   errorFilename,
				MaxSize:    int(c.Rotate.MaxSize),
				MaxBackups: int(c.Rotate.MaxBackups),
				MaxAge:     int(c.Rotate.MaxAge),
				Compress:   c.Rotate.Compress,
			}
			errorFileCore := zapcore.NewCore(encoder, zapcore.AddSync(errorRotateLogger), zapcore.ErrorLevel)
			cores = append(cores, errorFileCore)
		}
	}

	// 合并所有 core
	core := zapcore.NewTee(cores...)

	// 创建 logger
	zapLogger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(2))
	if c.Development {
		zapLogger = zapLogger.WithOptions(zap.Development())
	}

	return &ZapLogger{
		log:  zapLogger,
		Sync: zapLogger.Sync,
	}
}
