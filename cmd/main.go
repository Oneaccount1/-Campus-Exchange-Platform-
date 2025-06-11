package main

import (
	"campus/internal/bootstrap"
	"campus/internal/router"
	"fmt"
	"github.com/gin-gonic/gin"

	"log"
	"os"
	"path/filepath"
)

func main() {
	// 获取当前工作目录
	workDir, err := os.Getwd()

	if err != nil {
		log.Fatalf("获取工作目录失败: %v", err)
	}

	// 配置文件路径
	configPath := filepath.Join(workDir, "configs", "config.yaml")

	// 初始化应用
	if err := bootstrap.Bootstrap(configPath); err != nil {
		log.Fatalf("应用初始化失败: %v", err)
	}

	// 确保应用优雅关闭
	defer bootstrap.Shutdown()

	// 设置Gin模式
	serverConfig := bootstrap.GetConfig().Server
	gin.SetMode(serverConfig.Mode)

	// 创建Gin引擎
	r := gin.Default()

	// 注册路由
	router.RegisterRoutes(r)

	// 启动服务器
	addr := fmt.Sprintf(":%d", serverConfig.Port)
	log.Printf("服务器启动在 %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
