package conf

import (
	"fmt"
	"os"

	"github.com/nexptr/llmchain/llms"
	"go.uber.org/zap/zapcore"
	"gopkg.in/yaml.v2"
)

// Config 定义 配置结构图
type Config struct {

	//APIAddr , default: 0.0.0.0:8080
	APIAddr string `yaml:"api_addr"`

	LogLevel zapcore.Level `yaml:"log_level"`

	LogDir string `yaml:"log_dir"`

	ModelOptions []llms.ModelOptions `yaml:"model_options"`
}

func defaultConfig() *Config {
	cf := &Config{
		APIAddr:  "0.0.0.0:8080",
		LogLevel: zapcore.InfoLevel,
		LogDir:   "./logs",
	}
	return cf
}

func InitConfig(path string) (*Config, error) {

	f, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("cannot read config file: %w", err)
	}

	cf := defaultConfig()

	if err := yaml.Unmarshal(f, cf); err != nil {
		return nil, fmt.Errorf("cannot unmarshal config file: %w", err)
	}
	return cf, nil
}
