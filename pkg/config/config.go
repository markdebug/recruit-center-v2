package config

import (
	"fmt"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var (
	globalConfig *Config
	once         sync.Once
)

// AIConfig AI服务配置
type AIConfig struct {
	BaseURL    string        `yaml:"baseURL"`
	APIKey     string        `yaml:"apiKey"`
	ModelName  string        `yaml:"modelName"`
	Timeout    time.Duration `yaml:"timeout"`
	MaxRetries int           `yaml:"maxRetries"`
}

// Config holds the application configuration.
type Config struct {
	DB               DB               `mapstructure:"db"`          // Database configuration
	Oss              OssConfig        `mapstructure:"oss"`         // OSS configuration
	JWTConfig        JWTConfig        `mapstructure:"jwt"`         // JWT configuration
	FileUploadConfig FileUploadConfig `mapstructure:"file_upload"` // File upload configuration
	Name             string           `mapstructure:"name"`        // Application name
	Port             int              `mapstructure:"port"`        // Application port
	Host             string           `mapstructure:"host"`        // Application host
	Env              string           `mapstructure:"env"`         // Application environment (e.g., development, production)
	Log              LogConfig        `mapstructure:"log"`         // Logging configuration
	Version          string           `mapstructure:"version"`     // Application version
	System           SystemConfig     `mapstructure:"system"`      // System configuration
	AI               AIConfig         `mapstructure:"ai"`          // AI configuration
	v                *viper.Viper     `mapstructure:"-"`
}

type SystemConfig struct {
	UserTagLimit int `mapstructure:"user_tag_limit"` // 用户标签限制
}

type JWTConfig struct {
	Secret string `mapstructure:"secret"` // JWT密钥
}

type DB struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
}

type FileUploadConfig struct {
	MaxSize      int64    `mapstructure:"max_size"`      // 最大文件大小（字节）
	AllowedTypes []string `mapstructure:"allowed_types"` // 允许的MIME类型
	MaxFiles     int      `mapstructure:"max_files"`     // 最大文件数量
	AllowedExts  []string `mapstructure:"allowed_exts"`  // 允许的文件扩展名
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
	globalConfig = config
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

// 校验AI配置
func (c *Config) ValidateAI() error {
	if c.AI.BaseURL == "" {
		return fmt.Errorf("AI服务BaseURL不能为空")
	}
	if c.AI.APIKey == "" {
		return fmt.Errorf("AI服务APIKey不能为空")
	}
	if c.AI.ModelName == "" {
		return fmt.Errorf("AI服务ModelName不能为空")
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

// InitGlobalConfig 初始化全局配置
func InitGlobalConfig(configPath string) error {
	var err error
	once.Do(func() {
		globalConfig, err = LoadConfig(configPath)
	})
	return err
}

// GetConfig 获取全局配置实例
func GetConfig() *Config {
	if globalConfig == nil {
		panic("配置未初始化，请先调用 InitGlobalConfig")
	}
	return globalConfig
}
