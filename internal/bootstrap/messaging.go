package bootstrap

import (
	"campus/internal/messaging"
	"campus/internal/websocket"
	"errors"
	"log"
)

var (
	// 全局WebSocket管理器
	wsManager *websocket.Manager
)

// InitMessaging 初始化消息系统
func InitMessaging() error {
	// 获取配置
	config := GetConfig()
	if config == nil || config.RabbitMQ == nil {
		return errors.New("消息服务配置缺失")
	}

	// 创建WebSocket管理器
	wsManager = websocket.NewManager()

	// 启动WebSocket管理器
	go wsManager.Start()
	log.Println("WebSocket管理器已启动")

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

	log.Println("消息系统初始化成功")
	return nil
}

// GetWebSocketManager 获取WebSocket管理器
func GetWebSocketManager() *websocket.Manager {
	return wsManager
}
