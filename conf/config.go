package conf

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/nexptr/omnigram-server/store"
	"go.uber.org/zap/zapcore"
	"gopkg.in/yaml.v2"
)

// Config 定义 配置结构图
type Config struct {

	//APIAddr , default: 0.0.0.0:8080
	APIAddr string `yaml:"api_addr"`

	LogLevel zapcore.Level `yaml:"log_level"`

	LogDir string `yaml:"log_dir"`

	MetaDataPath string `yaml:"metadata_path"`

	DBOption *store.Opt `json:"db_options" yaml:"db_options"`

	ModelOptions []ModelOptions `yaml:"model_options"`

	EpubOptions EpubOptions `yaml:"epub_options"`
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

	if cf.DBOption == nil {
		cf.DBOption = &store.Opt{
			Driver:   "sqlite3",
			Host:     filepath.Join(cf.MetaDataPath, "db"),
			LogLevel: cf.LogLevel,
		}

		err = os.Mkdir(cf.DBOption.Host, 0755)
	}
	cf.DBOption.LogLevel = cf.LogLevel

	return cf, err
}

func defaultConfig() *Config {
	cf := &Config{
		APIAddr:      "0.0.0.0:8080",
		LogLevel:     zapcore.InfoLevel,
		LogDir:       "./logs",
		MetaDataPath: "./data",
	}
	return cf
}

type EpubOptions struct {
	DataPath           string `json:"data_path" yaml:"data_path"`
	SaveCoverBesideSrc bool   `json:"save_cover_beside_src" yaml:"save_cover_beside_src"`
	MaxEpubSize        int64  `json:"max_epub_size" yaml:"max_epub_size"`
}
