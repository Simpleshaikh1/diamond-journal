package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	AppName    string `mapstructure:"app_name"`
	Version    string `mapstructure:"version"`
	StorageDir string `mapstructure:"storage_dir"`
	DBFile     string `mapstructure:"db_file"`
	Editor     string `mapstructure:"editor"`
}

func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath(filepath.Join(os.Getenv("HOME"), ".diamond-journal"))

	viper.SetDefault("app_name", "Diamond Journal")
	viper.SetDefault("version", "0.1.0")
	viper.SetDefault("storage_dir", "storage")
	viper.SetDefault("db_file", "journal.db")
	viper.SetDefault("editor", "code")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			//create default config
			if err := os.MkdirAll(filepath.Join(os.Getenv("HOME"), ".diamond-journal"), 0775); err != nil {
				return nil, fmt.Errorf("failed to create default config: %w", err)
			}
		} else {
			return nil, err
		}
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	//ensure storage directory exist
	storagePath := cfg.StorageDir
	if !filepath.IsAbs(storagePath) {
		storagePath = filepath.Join(".", storagePath)
	}
	if err := os.MkdirAll(storagePath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create storage directory: %w", err)
	}

	fmt.Printf("✅ Config loaded. Storage: %s\n", storagePath)
	return &cfg, nil
}
