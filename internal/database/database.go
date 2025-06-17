package database

import (
	"campus/internal/config"
	"campus/internal/utils/logger"
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"time"
)

// GormLogger 自定义GORM日志适配器
type GormLogger struct {
	SlowThreshold        time.Duration
	IgnoreRecordNotFound bool
	LogLevel             gormlogger.LogLevel
	ModuleName           string
}

// LogMode 设置日志级别
func (l *GormLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	newLogger := *l
	newLogger.LogLevel = level
	return &newLogger

}

// Info 实现gorm日志接口
func (l *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormlogger.Info {
		logger.WithFields(map[string]interface{}{
			"module": l.ModuleName,
		}).Infof(msg, data...)
	}
}

// Warn 实现gorm日志接口
func (l *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormlogger.Warn {
		logger.WithFields(map[string]interface{}{
			"module": l.ModuleName,
		}).Warnf(msg, data...)
	}
}

// Error 实现gorm日志接口
func (l *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormlogger.Error {
		logger.WithFields(map[string]interface{}{
			"module": l.ModuleName,
		}).Errorf(msg, data...)
	}
}

// Trace 实现gorm日志接口，记录SQL执行
func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= gormlogger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	fields := map[string]interface{}{
		"module":  l.ModuleName,
		"elapsed": elapsed,
		"rows":    rows,
	}

	// 记录错误
	if err != nil && (!l.IgnoreRecordNotFound || !errors.Is(err, gorm.ErrRecordNotFound)) {
		fields["error"] = err
		logger.WithFields(fields).Errorf("SQL错误: %s", sql)
		return
	}

	// 记录慢查询
	if elapsed > l.SlowThreshold && l.SlowThreshold > 0 {
		logger.WithFields(fields).Warnf("慢查询 SQL: %s", sql)
		return
	}

	// 记录正常查询
	if l.LogLevel >= gormlogger.Info {
		logger.WithFields(fields).Debugf("SQL: %s", sql)
	}
}

// NewDatabase 创建数据库连接
func NewDatabase(dbConfig config.DatabaseConfig, serverMode string) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		dbConfig.Username,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.DBName,
		dbConfig.Charset,
	)

	// 根据服务器模式设置日志级别
	logMode := gormlogger.Error // 默认使用Error级别

	switch serverMode {
	case "debug":
		// 调试模式：显示所有SQL
		logMode = gormlogger.Info
	case "test":
		// 测试模式：显示慢查询和错误
		logMode = gormlogger.Warn
	case "production":
		// 生产模式：只显示错误
		logMode = gormlogger.Error
	}

	// 创建自定义日志配置
	gormLogger := &GormLogger{
		SlowThreshold:        time.Second, // 慢查询阈值
		LogLevel:             logMode,     // 日志级别
		IgnoreRecordNotFound: true,        // 忽略记录未找到错误
		ModuleName:           "gorm",      // 模块名称
	}

	gormConfig := &gorm.Config{
		Logger:                                   gormLogger,
		DisableForeignKeyConstraintWhenMigrating: true,
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
	sqlDB.SetConnMaxLifetime(dbConfig.ConnMaxLifetime) // 使用配置的连接生命周期

	// 将日志级别转换为字符串形式
	var logLevelStr string
	switch logMode {
	case gormlogger.Silent:
		logLevelStr = "Silent"
	case gormlogger.Error:
		logLevelStr = "Error"
	case gormlogger.Warn:
		logLevelStr = "Warn"
	case gormlogger.Info:
		logLevelStr = "Info"
	}

	logger.Info("数据库连接成功", zap.String("level", logLevelStr))
	return db, nil
}

// AutoMigrate 自动迁移数据库表结构
func AutoMigrate(db *gorm.DB, models ...interface{}) error {
	if err := db.AutoMigrate(models...); err != nil {
		return fmt.Errorf("自动迁移数据库表结构失败: %w", err)
	}
	return nil
}

// CloseDB 关闭数据库连接
func CloseDB(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("获取SQL DB失败: %w", err)
	}
	return sqlDB.Close()
}
