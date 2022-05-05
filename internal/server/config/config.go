package config

import "gophkeeper/pkg/logger"

type Config struct {
	DB     DatabaseConfig `mapstructure:"db"`
	Logger logger.Config  `mapstructure:"log"`
}

type DatabaseConfig struct {
	DSN string `mapstructure:"dsn"`
}
