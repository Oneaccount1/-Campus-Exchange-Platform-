package config

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"time"
)

// Config 应用配置结构体
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Upload   UploadConfig
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port int
	Mode string
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Driver      string
	Host        string
	Port        string
	Username    string
	Password    string
	DBName      string
	Charset     string
	MaxIdleConn int
	MaxOpenConn int
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret     string
	Expiration time.Duration
}

// UploadConfig 文件上传配置
type UploadConfig struct {
	SavePath     string
	AllowedTypes []string
	MaxSize      int
}

var (
	// AppConfig 全局配置实例
	AppConfig *Config
)

// LoadConfig 加载配置
func LoadConfig(configPath string) (*Config, error) {
	v := viper.New()
	v.SetConfigFile(configPath)
	v.SetConfigType("yaml")

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	config := &Config{}

	// 服务器配置
	config.Server.Port = v.GetInt("server.port")
	config.Server.Mode = v.GetString("server.mode")

	// 数据库配置
	config.Database.Driver = v.GetString("database.driver")
	config.Database.Host = v.GetString("database.host")
	config.Database.Port = v.GetString("database.port")
	config.Database.Username = v.GetString("database.username")
	config.Database.Password = v.GetString("database.password")
	config.Database.DBName = v.GetString("database.dbname")
	config.Database.Charset = v.GetString("database.charset")
	config.Database.MaxIdleConn = v.GetInt("database.max_idle_conns")
	config.Database.MaxOpenConn = v.GetInt("database.max_open_conns")

	// JWT配置
	config.JWT.Secret = v.GetString("jwt.secret")
	config.JWT.Expiration = time.Duration(v.GetInt("jwt.expiration")) * time.Hour

	// 上传配置
	config.Upload.SavePath = v.GetString("upload.save_path")
	config.Upload.AllowedTypes = v.GetStringSlice("upload.allowed_types")
	config.Upload.MaxSize = v.GetInt("upload.max_size")

	AppConfig = config

	return config, nil
}

// Init 初始化配置
func Init(configPath string) error {
	config, err := LoadConfig(configPath)
	if err != nil {
		return err
	}

	log.Println("配置加载成功")
	AppConfig = config

	return nil
}

// GetConfig 获取配置
func GetConfig() *Config {
	if AppConfig == nil {
		log.Panic("配置未初始化，请先调用 Init")
	}
	return AppConfig
}

// GetDatabaseConfig 获取数据库配置
func GetDatabaseConfig() DatabaseConfig {
	return GetConfig().Database
}

// GetJWTConfig 获取JWT配置
func GetJWTConfig() JWTConfig {
	return GetConfig().JWT
}

// GetServerConfig 获取服务器配置
func GetServerConfig() ServerConfig {
	return GetConfig().Server
}

// GetUploadConfig 获取上传配置
func GetUploadConfig() UploadConfig {
	return GetConfig().Upload
}
