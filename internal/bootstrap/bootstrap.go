package bootstrap

import (
	"campus/internal/config"
	"github.com/casbin/casbin/v2"
	"gorm.io/gorm"
	"log"
)

// 全局变量
var (
	// DB 全局数据库连接
	db *gorm.DB

	// Enforcer 全局Casbin执行器
	enforcer *casbin.Enforcer

	// AppConfig 全局配置
	appConfig *config.Config
)

// Bootstrap 初始化应用
func Bootstrap(configPath string) error {
	// 初始化配置
	if err := InitConfig(configPath); err != nil {
		return err
	}

	// 初始化数据库
	if err := InitDatabase(); err != nil {
		return err
	}

	// 初始化Casbin
	if err := InitCasbin(); err != nil {
		return err
	}

	// 初始化权限
	if err := InitPermissions(); err != nil {
		log.Printf("权限初始化失败: %v", err)
		// 权限初始化失败不应该阻止应用启动
		// 但我们应该记录错误
	}

	// 初始化消息系统
	if err := InitMessaging(); err != nil {
		log.Printf("消息系统初始化失败: %v", err)
		return err
	}

	log.Println("应用初始化完成")
	return nil
}

// Shutdown 优雅关闭应用
func Shutdown() error {
	// 关闭数据库连接
	if err := CloseDatabase(); err != nil {
		log.Printf("关闭数据库连接失败: %v", err)
		return err
	}

	log.Println("应用已关闭")
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
