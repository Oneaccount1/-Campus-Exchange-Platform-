package config

import (
	"fmt"
	"github.com/spf13/viper"
	"time"
)

// Config 应用配置结构体
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Upload   UploadConfig
	RabbitMQ *RabbitMQConfig
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port            int
	Mode            string
	InitPermissions bool `mapstructure:"init_permissions"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Driver          string
	Host            string
	Port            string
	Username        string
	Password        string
	DBName          string
	Charset         string
	MaxIdleConn     int
	MaxOpenConn     int
	ConnMaxLifetime time.Duration
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

// RabbitMQConfig RabbitMQ配置
type RabbitMQConfig struct {
	URL      string
	Username string
	Password string
	Host     string
	Port     string
}

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
	config.Server.InitPermissions = v.GetBool("server.init_permissions")

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
	config.Database.ConnMaxLifetime = time.Duration(v.GetInt("database.conn_max_lifetime")) * time.Second

	// JWT配置
	config.JWT.Secret = v.GetString("jwt.secret")
	config.JWT.Expiration = time.Duration(v.GetInt("jwt.expiration")) * time.Hour

	// 上传配置
	config.Upload.SavePath = v.GetString("upload.save_path")
	config.Upload.AllowedTypes = v.GetStringSlice("upload.allowed_types")
	config.Upload.MaxSize = v.GetInt("upload.max_size")

	// RabbitMQ配置
	if v.IsSet("rabbitmq") {
		config.RabbitMQ = &RabbitMQConfig{
			URL:      v.GetString("rabbitmq.url"),
			Username: v.GetString("rabbitmq.username"),
			Password: v.GetString("rabbitmq.password"),
			Host:     v.GetString("rabbitmq.host"),
			Port:     v.GetString("rabbitmq.port"),
		}

		// 如果URL未设置，但设置了其他参数，则构造URL
		if config.RabbitMQ.URL == "" && config.RabbitMQ.Host != "" {
			user := "guest"
			pass := "guest"
			if config.RabbitMQ.Username != "" {
				user = config.RabbitMQ.Username
			}
			if config.RabbitMQ.Password != "" {
				pass = config.RabbitMQ.Password
			}
			port := "5672"
			if config.RabbitMQ.Port != "" {
				port = config.RabbitMQ.Port
			}
			config.RabbitMQ.URL = fmt.Sprintf("amqp://%s:%s@%s:%s/", user, pass, config.RabbitMQ.Host, port)
		}

		// 如果没有设置URL和其他参数，使用默认本地URL
		if config.RabbitMQ.URL == "" {
			config.RabbitMQ.URL = "amqp://guest:guest@localhost:5672/"
		}
	} else {
		// 默认RabbitMQ配置
		config.RabbitMQ = &RabbitMQConfig{
			URL: "amqp://guest:guest@localhost:5672/",
		}
	}

	return config, nil
}
