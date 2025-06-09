package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github/oneaccount1/-Campus-Exchange-Platform-/internal/database"
	"github/oneaccount1/-Campus-Exchange-Platform-/internal/models"
	"github/oneaccount1/-Campus-Exchange-Platform-/internal/router"
	"github/oneaccount1/-Campus-Exchange-Platform-/internal/utils/config"
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

	// 加载配置
	configPath := filepath.Join(workDir, "configs", "config.yaml")
	if err := config.Init(configPath); err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 初始化数据库
	if err := database.Init(); err != nil {
		log.Fatalf("初始化数据库失败: %v", err)
	}
	defer database.Close()

	// 自动迁移数据库表结构
	if err := database.AutoMigrate(
		&models.User{},
		&models.Product{},
		&models.Order{},
		&models.Review{},
		&models.Message{},
	); err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}

	// 设置Gin模式
	serverConfig := config.GetServerConfig()
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
