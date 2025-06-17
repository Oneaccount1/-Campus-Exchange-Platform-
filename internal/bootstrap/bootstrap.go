package bootstrap

import (
	"campus/internal/config"
	"campus/internal/utils/logger"
	"campus/internal/websocket"
	"fmt"
	"github.com/casbin/casbin/v2"
	"gorm.io/gorm"
)

// 全局变量
var (
	// DB 全局数据库连接
	db *gorm.DB

	// Enforcer 全局Casbin执行器
	enforcer *casbin.Enforcer

	// AppConfig 全局配置
	appConfig *config.Config

	// 全局WebSocket管理器
	wsManager *websocket.Manager
)

// Bootstrap 初始化应用
func Bootstrap(configPath string) error {
	// 先加载配置，但不使用日志记录
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		fmt.Printf("配置初始化失败: %v\n", err)
		return err
	}

	// 设置全局配置
	SetConfig(cfg)

	// 初始化日志系统
	InitLogger()

	// 现在日志系统已经初始化，可以记录信息
	logger.Info("配置和日志初始化成功")

	// 初始化数据库
	if err := InitDatabase(); err != nil {
		return err
	}

	//初始化Casbin
	if err := InitCasbin(); err != nil {
		return err
	}

	//初始化权限
	if err := InitPermissions(); err != nil {
		logger.Errorf("权限初始化失败: %v", err)
		// 权限初始化失败不应该阻止应用启动
		// 但我们应该记录错误
	}

	// 初始化消息系统
	if err := InitMessaging(); err != nil {
		logger.Errorf("消息系统初始化失败: %v", err)
		return err
	}

	logger.Info("应用初始化完成")
	return nil
}

// Shutdown 优雅关闭应用
func Shutdown() error {
	// 关闭数据库连接
	if err := CloseDatabase(); err != nil {
		logger.Errorf("关闭数据库连接失败: %v", err)
		return err
	}

	// 同步日志
	logger.Sync()

	logger.Info("应用已关闭")
	return nil
}

// GetDB 获取全局数据库连接
func GetDB() *gorm.DB {
	return db
}

// SetDB 设置全局数据库连接（内部使用）
func SetDB(database *gorm.DB) {
	db = database
}

// GetEnforcer 获取全局Casbin执行器
func GetEnforcer() *casbin.Enforcer {
	return enforcer
}

// SetEnforcer 设置全局Casbin执行器（内部使用）
func SetEnforcer(e *casbin.Enforcer) {
	enforcer = e
}

// GetConfig 获取全局配置
func GetConfig() *config.Config {
	return appConfig
}

// SetConfig 设置全局配置（内部使用）
func SetConfig(cfg *config.Config) {
	appConfig = cfg
}

// GetWebSocketManager 获取WebSocket管理器
func GetWebSocketManager() *websocket.Manager {
	return wsManager
}

func SetWebSocketManager(websocketManager *websocket.Manager) {
	wsManager = websocketManager
}
