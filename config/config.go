package config

import (
	"fmt"
	"log/slog"

	"github.com/spf13/viper"
)

type Config struct {
	Server ServerConfig `mapstructure:"server"`
}

type ServerConfig struct {
	Port int    `mapstructure:"port"`
	Host string `mapstructure:"host"`
}

func MustLoad() (Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		return Config{}, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return Config{}, fmt.Errorf("config parsing error: %w", err)
	}

	if cfg.Server.Port <= 0 || cfg.Server.Port > 65535 {
		return Config{}, fmt.Errorf("invalid server port: %d (must be 1-65535)", cfg.Server.Port)
	}
	if cfg.Server.Host == "" {
		return Config{}, fmt.Errorf("server host cannot be empty")
	}

	slog.Info("config loaded", "host", cfg.Server.Host, "port", cfg.Server.Port)

	return cfg, nil
}
