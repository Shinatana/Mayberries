package conf

import (
	"catalog_service/pkg/config"
)

type Config struct {
	Http config.HttpOptions     `mapstructure:"http"`
	Log  config.LoggerOptions   `mapstructure:"log"`
	DB   config.DatabaseOptions `mapstructure:"db"`
}

type Loader interface {
	Load() (*Config, error)
}
