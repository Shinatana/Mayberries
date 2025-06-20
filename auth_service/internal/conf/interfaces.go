package conf

import (
	"auth_service/pkg/config"
)

type Config struct {
	Http    config.HttpOptions      `mapstructure:"http"`
	Log     config.LoggerOptions    `mapstructure:"log"`
	DB      config.DatabaseOptions  `mapstructure:"db"`
	Migrate config.MigrationOptions `mapstructure:"migrate"`
}

type Loader interface {
	Load() (*Config, error)
}
