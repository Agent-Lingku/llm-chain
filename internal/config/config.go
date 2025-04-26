package config

import (
	"fmt"
	"log"
	"strings"

	"github.com/spf13/viper"
)

// Config 系统配置
type Config struct {
	ApiBaseUrl string `mapstructure:"apiBaseUrl"`
	ApiBaseKey string `mapstructure:"apiBaseKey"`
	Prefix     string `mapstructure:"prefix"`
	LogLevel   string `mapstructure:"logLevel"`
}

var v *viper.Viper

func init() {
	v = viper.New()
	v.SetConfigName("config") 
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.AddConfigPath("./configs")
	v.AddConfigPath("/etc/appname/")
	
	// 设置环境变量前缀并自动绑定
	v.SetEnvPrefix("APP")
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	
	// 设置默认值
	v.SetDefault("apiBaseUrl", "http://127.0.0.1:11434")
	v.SetDefault("apiBaseKey", "sk-xxx")
	v.SetDefault("prefix", "/api/chat")
	v.SetDefault("logLevel", "info")
}

// LoadConfig 加载并验证配置
func LoadConfig() (*Config, error) {
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("读取配置文件失败: %w", err)
		}
		log.Println("使用默认配置和环境变量")
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("配置解析失败: %w", err)
	}

	if err := validateConfig(&cfg); err != nil {
		return nil, fmt.Errorf("配置验证失败: %w", err)
	}

	log.Printf("加载配置成功: ApiBaseUrl=%s, LogLevel=%s\n", cfg.ApiBaseUrl, cfg.LogLevel)
	return &cfg, nil
}

// validateConfig 验证配置合法性
func validateConfig(cfg *Config) error {
	if cfg.ApiBaseUrl == "" {
		return fmt.Errorf("apiBaseUrl 必须配置")
	}
	if !strings.HasPrefix(cfg.ApiBaseUrl, "http") {
		return fmt.Errorf("apiBaseUrl 必须以 http 或 https 开头")
	}
	return nil
}

// 保持原有常量兼容
const (
	OllamaUrl    = "http://127.0.0.1:11434"
	OllamaPrefix = "/api/chat"
)

// 兼容旧版DefaultConfig（建议逐步迁移）
var DefaultConfig = Config{
	ApiBaseUrl: "http://127.0.0.1:11434",
	ApiBaseKey: "sk-xxxx",
	Prefix:     "/api/chat",
	LogLevel:   "info",
}
