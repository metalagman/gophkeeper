package config

import (
	"gophkeeper/pkg/logger"
)

type Config struct {
	GRPC     GRPCConfig     `mapstructure:"grpc"`
	DB       DatabaseConfig `mapstructure:"db"`
	Security SecurityConfig `mapstructure:"security"`
	Logger   logger.Config  `mapstructure:"log"`
}

type GRPCConfig struct {
	ListenAddr string `mapstructure:"listen_addr"`
}

type DatabaseConfig struct {
	DSN string `mapstructure:"dsn"`
}

type SecurityConfig struct {
	SecretKey string `mapstructure:"secret_key"`
}
