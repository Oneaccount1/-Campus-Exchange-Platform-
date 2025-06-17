package logger

import (
	"campus/internal/config"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	// Logger 全局日志实例
	Logger *zap.Logger
	// Sugar 简化API的日志实例
	Sugar *zap.SugaredLogger
)

// Init 初始化日志系统
func Init(cfg *config.LogConfig) {
	// 确保日志目录存在
	if cfg.Output.File.Path != "" {
		dir := filepath.Dir(cfg.Output.File.Path)
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			os.MkdirAll(dir, 0755)
		}
	}

	// 配置encoder
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

	// 根据配置选择编码器
	var encoder zapcore.Encoder
	if cfg.Format == "json" {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	// 创建输出源
	var cores []zapcore.Core
	level := getLogLevel(cfg.Level)

	// 控制台输出
	if cfg.Output.Console {
		cores = append(cores, zapcore.NewCore(
			encoder,
			zapcore.AddSync(os.Stdout),
			level,
		))
	}

	// 文件输出
	if cfg.Output.File.Path != "" {
		// 这里简化实现，实际应该使用lumberjack等库处理日志轮转
		file, err := os.OpenFile(cfg.Output.File.Path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err == nil {
			cores = append(cores, zapcore.NewCore(
				encoder,
				zapcore.AddSync(file),
				level,
			))
		}
	}

	// 创建zap logger
	core := zapcore.NewTee(cores...)
	Logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	Sugar = Logger.Sugar()

	// 替换全局logger
	zap.ReplaceGlobals(Logger)
}

// 获取日志级别
func getLogLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}

// 以下是日志快捷方法

func Debug(msg string, fields ...zap.Field) {
	Logger.Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	Logger.Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	Logger.Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	Logger.Error(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	Logger.Fatal(msg, fields...)
}

// 带格式的日志

func Debugf(format string, args ...interface{}) {
	Sugar.Debugf(format, args...)
}

func Infof(format string, args ...interface{}) {
	Sugar.Infof(format, args...)
}

func Warnf(format string, args ...interface{}) {
	Sugar.Warnf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	Sugar.Errorf(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	Sugar.Fatalf(format, args...)
}

// WithFields 创建带有字段的日志记录器
func WithFields(fields map[string]interface{}) *zap.SugaredLogger {
	var zapFields []zap.Field
	for k, v := range fields {
		zapFields = append(zapFields, zap.Any(k, v))
	}
	return Logger.With(zapFields...).Sugar()
}

// Sync 同步日志，应在程序退出前调用
func Sync() {
	Logger.Sync()
}
