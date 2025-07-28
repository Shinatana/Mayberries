package conf

import (
	"github.com/mayberries/shared/pkg/config"
)

type Config struct {
	Http config.HttpOptions     `mapstructure:"http"`
	Log  config.LoggerOptions   `mapstructure:"log"`
	DB   config.DatabaseOptions `mapstructure:"db"`
}

type Loader interface {
	Load() (*Config, error)
}
