package config

type OssConfig struct {
	Endpoint        string `mapstructure:"endpoint"`          // OSS Endpoint
	AccessKeyID     string `mapstructure:"access_key_id"`     // Access Key ID
	SecretAccessKey string `mapstructure:"secret_access_key"` // Secret Access Key
	UseSSL          bool   `mapstructure:"use_ssl"`           // 是否使用SSL
}
