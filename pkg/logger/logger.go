package logger

import (
	"blog/internal/conf"
	"fmt"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"strings"
)
import "github.com/go-kratos/kratos/v2/log"

var _ log.Logger = (*Logger)(nil)

type Logger struct {
	log    *zap.Logger
	msgKey string
}

type Option func(*Logger)

// WithMessageKey with message key.
func WithMessageKey(key string) Option {
	return func(l *Logger) {
		l.msgKey = key
	}
}

func NewLogger(c *conf.Bootstrap) *Logger {
	return &Logger{
		log:    getLogger(c),
		msgKey: log.DefaultMessageKey,
	}
}

func (l *Logger) Log(level log.Level, keyvals ...interface{}) error {
	// If logging at this level is completely disabled, skip the overhead of
	// string formatting.
	if zapcore.Level(level) < zapcore.DPanicLevel && !l.log.Core().Enabled(zapcore.Level(level)) {
		return nil
	}
	var (
		msg    = ""
		keylen = len(keyvals)
	)
	if keylen == 0 || keylen%2 != 0 {
		l.log.Warn(fmt.Sprint("Keyvalues must appear in pairs: ", keyvals))
		return nil
	}

	data := make([]zap.Field, 0, (keylen/2)+1)
	for i := 0; i < keylen; i += 2 {
		if keyvals[i].(string) == l.msgKey {
			msg, _ = keyvals[i+1].(string)
			continue
		}
		data = append(data, zap.Any(fmt.Sprint(keyvals[i]), keyvals[i+1]))
	}
	switch level {
	case log.LevelDebug:
		l.log.Debug(msg, data...)
	case log.LevelInfo:
		l.log.Info(msg, data...)
	case log.LevelWarn:
		l.log.Warn(msg, data...)
	case log.LevelError:
		l.log.Error(msg, data...)
	case log.LevelFatal:
		l.log.Fatal(msg, data...)
	}
	return nil
}

func (l *Logger) Sync() error {
	return l.log.Sync()
}

func (l *Logger) Close() error {
	return l.Sync()
}

func getLogger(c *conf.Bootstrap) *zap.Logger {
	encoder := getEncoder(c)
	syncer := getWriterSyncer(c)
	core := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(syncer), getLevel(c.Log.Level))
	logger := zap.New(core, zap.AddCaller())
	return logger
}
func getLevel(level string) zapcore.Level {
	switch strings.ToLower(level) {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "dpanic":
		return zapcore.DPanicLevel
	case "panic":
		return zapcore.PanicLevel
	case "fatal":
		return zapcore.FatalLevel
	}
	return zapcore.InfoLevel
}

func getEncoder(c *conf.Bootstrap) zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder   //修改时间编码器
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder //在日志文件中使用大写字母记录日志级别
	//return zapcore.NewConsoleEncoder(zap.NewProductionEncoderConfig())
	//return zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	return zapcore.NewJSONEncoder(encoderConfig)
}

func getWriterSyncer(c *conf.Bootstrap) zapcore.WriteSyncer {
	logger := &lumberjack.Logger{
		Filename:   fmt.Sprintf("./%s/%s.log", c.Log.Path, c.Log.Name),
		MaxSize:    int(c.Log.MaxSize),
		MaxBackups: int(c.Log.MaxBackups),
		MaxAge:     int(c.Log.MaxAge),
		Compress:   c.Log.Compress,
	}
	return zapcore.AddSync(logger)
}
