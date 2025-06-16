package bootstrap

import (
	"campus/internal/database"
	"campus/internal/models"
	"log"
)

// InitDatabase 初始化数据库
func InitDatabase() error {
	// 获取配置
	config := GetConfig()

	// 初始化数据库连接，传递服务器模式
	db, err := database.NewDatabase(config.Database, config.Server.Mode)
	if err != nil {
		log.Printf("数据库连接初始化失败: %v", err)
		return err
	}

	// 设置全局数据库连接
	SetDB(db)

	// 自动迁移数据库表结构
	//if err := AutoMigrateModels(); err != nil {
	//	log.Printf("数据库迁移失败: %v", err)
	//	return err
	//}

	log.Println("数据库初始化成功")
	return nil
}

// AutoMigrateModels 自动迁移数据库表结构
func AutoMigrateModels() error {
	return database.AutoMigrate(
		GetDB(),
		&models.User{},
		&models.Role{},
		&models.Permission{},
		&models.UserRole{},
		&models.RolePermission{},
		&models.Product{},
		&models.Order{},
		&models.Review{},
		&models.Message{},
		&models.ProductImage{},
	)
}

// CloseDatabase 关闭数据库连接
func CloseDatabase() error {
	if db == nil {
		return nil
	}
	return database.CloseDB(GetDB())
}
