package bootstrap

import (
	"campus/internal/messaging"
	"campus/internal/utils/logger"
	"campus/internal/websocket"
	"errors"
)

// InitMessaging 初始化消息系统
func InitMessaging() error {
	// 获取配置
	config := GetConfig()
	if config == nil || config.RabbitMQ == nil {
		return errors.New("消息服务配置缺失")
	}

	// 创建WebSocket管理器
	NewWsManager := websocket.NewManager()

	// 启动WebSocket管理器
	go NewWsManager.Start()
	logger.Info("WebSocket管理器已启动")

	// 测试RabbitMQ连接
	rabbitURL := config.RabbitMQ.URL
	testRMQ, err := messaging.NewRabbitMQ(
		rabbitURL,
		"test_exchange",
		"test_queue",
		"test_route",
	)
	if err != nil {
		return errors.New("RabbitMQ连接失败: " + err.Error())
	}
	defer testRMQ.Close()

	logger.Info("消息系统初始化成功")
	SetWebSocketManager(NewWsManager)
	return nil
}
