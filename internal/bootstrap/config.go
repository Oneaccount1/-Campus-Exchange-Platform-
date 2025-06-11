package bootstrap

import (
	"campus/internal/config"
	"log"
)

// InitConfig 初始化配置
func InitConfig(configPath string) error {
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Printf("配置初始化失败: %v", err)
		return err
	}

	// 设置全局配置
	SetConfig(cfg)

	log.Println("配置初始化成功")
	return nil
}
