package bootstrap

import (
	"campus/internal/utils/logger"
)

// InitLogger 初始化日志系统
// 这个函数不会生成错误，所以移除了错误返回值
func InitLogger() {
	cfg := GetConfig()
	logger.Init(&cfg.Log)
}
