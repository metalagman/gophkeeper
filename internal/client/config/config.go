package config

import "gophkeeper/pkg/logger"

type Config struct {
	Server ServerConfig  `mapstructure:"server"`
	Logger logger.Config `mapstructure:"log"`
}

type ServerConfig struct {
	Addr string `mapstructure:"addr"`
}
