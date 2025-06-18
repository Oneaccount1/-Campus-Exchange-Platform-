package bootstrap

import (
	"campus/internal/rabbitMQ"
	"campus/internal/utils/logger"
	"campus/internal/websocket"
	"errors"
)

// InitMessaging initializes the rabbitMQ system, including WebSocket manager and the RabbitMQ consumer.
func InitMessaging() error {
	config := GetConfig()
	if config == nil || config.RabbitMQ == nil {
		return errors.New("消息服务配置缺失")
	}

	// 1. Create and start the WebSocket manager
	wsManager := websocket.NewManager()
	go wsManager.Start()
	SetWebSocketManager(wsManager)
	logger.Info("WebSocket管理器已启动")

	// 2. Start the background message consumer
	rabbitURL := config.RabbitMQ.URL
	go rabbitMQ.StartConsumer(rabbitURL, wsManager)

	logger.Info("消息系统初始化成功")
	return nil
}
