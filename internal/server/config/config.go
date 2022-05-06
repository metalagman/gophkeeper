package config

import (
	"gophkeeper/pkg/grpcserver"
	"gophkeeper/pkg/logger"
)

type Config struct {
	GRPC   grpcserver.Config `mapstructure:"grpc"`
	DB     DatabaseConfig    `mapstructure:"db"`
	Logger logger.Config     `mapstructure:"log"`
}

type DatabaseConfig struct {
	DSN string `mapstructure:"dsn"`
}
