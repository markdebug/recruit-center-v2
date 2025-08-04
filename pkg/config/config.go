package config

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// Config holds the application configuration.
type Config struct {
	DB      DB           `mapstructure:"db"`      // Database configuration
	Name    string       `mapstructure:"name"`    // Application name
	Port    int          `mapstructure:"port"`    // Application port
	Env     string       `mapstructure:"env"`     // Application environment (e.g., development, production)
	Log     LogConfig    `mapstructure:"log"`     // Logging configuration
	Version string       `mapstructure:"version"` // Application version
	v       *viper.Viper `mapstructure:"-"`
}

type DB struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
}

// LoadConfig loads the configuration from a file or environment variables.
func LoadConfig(filePath string) (*Config, error) {
	if filePath == "" {
		filePath = "config.yaml" // Default config file path
	}
	v := viper.New()
	v.SetConfigFile(filePath)
	v.SetConfigType("yaml")
	// 读取配置文件
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}
	config := &Config{}
	// 将配置解析到结构体中
	if err := v.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}
	// 监听配置文件变化
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		// 配置文件发生变更时重新解析
		if err := v.Unmarshal(config); err != nil {
			fmt.Printf("配置重新加载错误: %v\n", err)
			return
		}
		fmt.Println("配置已更新")
	})

	return config, nil
}

func (c *Config) Validate() error {
	if c.Name == "" {
		return fmt.Errorf("application name is required")
	}
	if c.Port <= 0 {
		return fmt.Errorf("invalid port number")
	}
	if c.DB.Host == "" {
		return fmt.Errorf("database host is required")
	}
	if c.DB.Port <= 0 {
		return fmt.Errorf("invalid database port number")
	}
	if c.DB.Username == "" {
		return fmt.Errorf("database username is required")
	}
	if c.DB.Password == "" {
		return fmt.Errorf("database password is required")
	}
	if c.DB.Database == "" {
		return fmt.Errorf("database name is required")
	}
	return nil
}

// Reload 重新加载配置
func (c *Config) Reload() error {
	if err := c.v.ReadInConfig(); err != nil {
		return fmt.Errorf("重新加载配置失败: %w", err)
	}
	if err := c.v.Unmarshal(c); err != nil {
		return fmt.Errorf("重新解析配置失败: %w", err)
	}
	return nil
}
