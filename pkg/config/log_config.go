package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	L *zap.Logger // 全局日志实例
)

// LogConfig 日志配置
type LogConfig struct {
	Level      string `mapstructure:"level"`      // 日志级别
	Filename   string `mapstructure:"filename"`   // 日志文件
	MaxSize    int    `mapstructure:"maxsize"`    // 单个文件最大大小，单位MB
	MaxAge     int    `mapstructure:"maxage"`     // 文件最大保存天数
	MaxBackups int    `mapstructure:"maxbackups"` // 最大保留文件数
}

// InitLogger 初始化日志
func InitLogger(cfg *LogConfig) (*zap.Logger, error) {
	// 创建日志目录
	logDir := filepath.Dir(cfg.Filename)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, fmt.Errorf("创建日志目录失败: %v", err)
	}

	// 设置日志级别
	level := zap.InfoLevel
	if err := level.UnmarshalText([]byte(cfg.Level)); err != nil {
		return nil, fmt.Errorf("解析日志级别失败: %v", err)
	}

	// 配置编码器
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     customTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// 配置输出
	jsonEncoder := zapcore.NewJSONEncoder(encoderConfig)

	// 写入文件
	fileWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   cfg.Filename,
		MaxSize:    cfg.MaxSize,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge,
	})

	// 同时输出到控制台和文件
	core := zapcore.NewTee(
		zapcore.NewCore(jsonEncoder, zapcore.AddSync(os.Stdout), level),
		zapcore.NewCore(jsonEncoder, fileWriter, level),
	)

	// 创建logger
	logger := zap.New(core,
		zap.AddCaller(),                   // 添加调用者信息
		zap.AddCallerSkip(1),              // 跳过封装函数
		zap.AddStacktrace(zap.ErrorLevel), // Error级别添加堆栈跟踪
	)

	return logger, nil
}

// customTimeEncoder 自定义时间编码器
func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}
