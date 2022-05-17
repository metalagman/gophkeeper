package config

import (
	"gophkeeper/pkg/grpcserver"
	"gophkeeper/pkg/logger"
)

type Config struct {
	GRPC     grpcserver.Config `mapstructure:"grpc"`
	DB       DatabaseConfig    `mapstructure:"db"`
	Security SecurityConfig    `mapstructure:"security"`
	Logger   logger.Config     `mapstructure:"log"`
}

type DatabaseConfig struct {
	DSN string `mapstructure:"dsn"`
}

type SecurityConfig struct {
	SecretKey string `mapstructure:"secret_key"`
}
