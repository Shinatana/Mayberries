package config

import (
	"errors"
	"time"
)

var ErrUnexpectedType = errors.New("unexpected type for flag")

// HttpOptions expected usage:
//
//	type Config struct {
//	   ...
//	   Http HttpOptions `mapstructure:"http"`
//	   ...
//	}
type HttpOptions struct {
	Host              string        `mapstructure:"host" validate:"required,hostname_rfc1123|ip_addr"`
	Port              int           `mapstructure:"port" validate:"required,gt=1023,lt=65536"`
	MaxHeaderBytes    int           `mapstructure:"max_header_bytes" validate:"min=1,max=1048576"`
	ReadTimeout       time.Duration `mapstructure:"read_timeout" validate:"min=500ms,max=30s"`
	ReadHeaderTimeout time.Duration `mapstructure:"read_header_timeout" validate:"min=500ms,max=30s"`
	WriteTimeout      time.Duration `mapstructure:"write_timeout" validate:"min=1s,max=60s"`
	IdleTimeout       time.Duration `mapstructure:"idle_timeout" validate:"min=1m,max=2m"`
	ShutdownTimeout   time.Duration `mapstructure:"shutdown_timeout" validate:"min=1s,max=60s"`
}

type LoggerOptions struct {
	Format string `mapstructure:"format" validate:"oneof=json text"`
	Level  string `mapstructure:"level" validate:"oneof=debug info warn error"`
}

// MigrationOptions expected usage:pwd
//
//	type Config struct {
//	   ...
//	   Miration MigrationOptions `mapstructure:"migrate"`
//	   ...
//	}
type MigrationOptions struct {
	Version          int           `mapstructure:"version" validate:"gt=-2"`
	MigrationFiles   string        `mapstructure:"dir" validate:"omitempty,dir"`
	MigrationTimeout time.Duration `mapstructure:"timeout" validate:"min=1s,max=60s"`
}

// DatabaseOptions expected usage:
//
//	type Config struct {
//	   ...
//	   DB DatabaseOptions `mapstructure:"db"`
//	   ...
//	}
type DatabaseOptions struct {
	Host               string        `mapstructure:"host" validate:"required,hostname_rfc1123|ip_addr"`
	Port               int           `mapstructure:"port" validate:"gt=1023,lt=65536"`
	User               string        `mapstructure:"user" validate:"required"`
	Pwd                string        `mapstructure:"pwd" validate:"required"`
	Database           string        `mapstructure:"database" validate:"required"`
	SSL                string        `mapstructure:"ssl" validate:"oneof=require disable verify-full verify-ca allow prefer"`
	Schema             string        `mapstructure:"schema" validate:"omitempty,min=1"`
	MaxOpenConnections int           `mapstructure:"max_open" validate:"gt=0"`
	MaxIdleConnections int           `mapstructure:"max_idle" validate:"gte=0"`
	ConnMaxLifetime    time.Duration `mapstructure:"max_lifetime" validate:"min=1s,max=1h"`
	InitTimeout        time.Duration `mapstructure:"init_timeout" validate:"min=1s,max=60s"`
}
