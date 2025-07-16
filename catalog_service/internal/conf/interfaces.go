package conf

import (
	"catalog_service/pkg/config"
)

type Config struct {
	Http  config.HttpOptions     `mapstructure:"http"`
	Log   config.LoggerOptions   `mapstructure:"log"`
	DB    config.DatabaseOptions `mapstructure:"db"`
	Redis config.RedisOptions    `mapstructure:"redis"`
}

type Loader interface {
	Load() (*Config, error)
}
