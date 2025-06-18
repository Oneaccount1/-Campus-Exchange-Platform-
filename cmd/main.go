package main

import (
	"campus/internal/bootstrap"
	"campus/internal/middleware"
	"campus/internal/router"
	"campus/internal/utils/logger"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"os"
	"path/filepath"
)

func main() {
	// 获取当前工作目录
	workDir, err := os.Getwd()

	if err != nil {
		// 还没有初始化日志系统，使用标准日志
		fmt.Printf("获取工作目录失败: %v\n", err)
		os.Exit(1)
	}

	// 配置文件路径
	configPath := filepath.Join(workDir, "configs", "config.yaml")

	// 初始化应用
	if err := bootstrap.Bootstrap(configPath); err != nil {
		// bootstrap过程中如果出错，可能已经初始化了日志系统
		logger.Fatalf("应用初始化失败: %v", err)
	}

	// 确保应用优雅关闭
	defer bootstrap.Shutdown()

	// 设置Gin模式
	serverConfig := bootstrap.GetConfig().Server
	gin.SetMode(serverConfig.Mode)

	// 创建Gin引擎
	r := gin.Default()

	// 应用CORS中间件
	r.Use(middleware.CORS())

	// 全局panic捕获
	r.Use(gin.Recovery())

	// 注册路由
	router.RegisterRoutes(r)

	// 启动服务器
	addr := fmt.Sprintf(":%d", serverConfig.Port)
	logger.Info("服务器启动", zap.String("地址", addr))
	if err := r.Run(addr); err != nil {
		logger.Fatalf("服务器启动失败: %v", err)
	}
}
