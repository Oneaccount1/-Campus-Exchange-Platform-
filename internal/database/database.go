package database

import (
	"campus/internal/utils/config"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"sync"
	"time"
)

// Database 数据库实例结构体
type Database struct {
	DB *gorm.DB
}

// 全局数据库实例
var (
	instance *Database
	once     sync.Once
	mu       sync.RWMutex
)

// Init 初始化数据库连接
func Init() error {
	var err error
	once.Do(func() {
		dbConfig := config.GetDatabaseConfig()
		instance, err = connect(dbConfig)
	})
	return err
}

// connect 连接数据库（内部函数）
func connect(dbConfig config.DatabaseConfig) (*Database, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		dbConfig.Username,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.DBName,
		dbConfig.Charset,
	)

	// 设置日志级别
	logMode := logger.Error
	if config.GetServerConfig().Mode == "debug" {
		logMode = logger.Info
	}

	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logMode),
	}

	db, err := gorm.Open(mysql.Open(dsn), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("数据库连接失败: %v", err)
	}

	// 配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("配置数据库连接池失败: %w", err)
	}

	sqlDB.SetMaxIdleConns(dbConfig.MaxIdleConn)
	sqlDB.SetMaxOpenConns(dbConfig.MaxOpenConn)
	sqlDB.SetConnMaxLifetime(time.Hour) // 连接生命周期固定为1小时

	log.Println("数据库连接成功")
	return &Database{DB: db}, nil
}

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	mu.RLock()
	defer mu.RUnlock()

	if instance == nil {
		log.Panic("数据库未初始化，请先调用 Init")
	}
	return instance.DB
}

// Close 关闭数据库连接
func Close() error {
	if instance == nil {
		return nil
	}

	mu.Lock()
	defer mu.Unlock()

	sqlDB, err := instance.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// Transaction 执行事务
func Transaction(fn func(tx *gorm.DB) error) error {
	return GetDB().Transaction(fn)
}

// AutoMigrate 自动迁移数据库表结构
func AutoMigrate(models ...interface{}) error {
	return GetDB().AutoMigrate(models...)
}

// Ping 检查数据库连接是否正常
func Ping() error {
	sqlDB, err := GetDB().DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}

// Stats 获取数据库连接池统计信息
func Stats() interface{} {
	sqlDB, err := GetDB().DB()
	if err != nil {
		return nil
	}
	return sqlDB.Stats()
}
