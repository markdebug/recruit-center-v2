package logger

import (
	"go.uber.org/zap"
	"org.thinkinai.com/recruit-center/pkg/config"
)

var (
	L *zap.Logger // 全局日志实例
)

// Init 初始化日志
func Init(cfg *config.LogConfig) error {
	logger, err := config.InitLogger(cfg)
	if err != nil {
		return err
	}
	L = logger
	return nil
}

// Sync 同步日志缓存
func Sync() {
	if L != nil {
		L.Sync()
	}
}
