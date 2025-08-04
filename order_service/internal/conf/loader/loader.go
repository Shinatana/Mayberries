package loader

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/mayberries/shared/pkg/val"
	"order_service/internal/conf"
)

const (
	defaultConfigFile = ".env"
)

type loader struct{}

func NewLoader() conf.Loader {
	return &loader{}
}

func (l *loader) Load() (*conf.Config, error) {
	if err := godotenv.Load(defaultConfigFile); err != nil {
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("failed to load .env file: %w", err)
		}
	}

	var cfg conf.Config

	// HTTP
	cfg.Http.Host = os.Getenv("HTTP_HOST")
	cfg.Http.Port, _ = parseInt(os.Getenv("HTTP_PORT"))
	cfg.Http.MaxHeaderBytes, _ = parseInt(os.Getenv("HTTP_MAX_HEADER_BYTES"))
	cfg.Http.ReadTimeout, _ = parseDuration(os.Getenv("HTTP_READ_TIMEOUT"))
	cfg.Http.WriteTimeout, _ = parseDuration(os.Getenv("HTTP_WRITE_TIMEOUT"))
	cfg.Http.IdleTimeout, _ = parseDuration(os.Getenv("HTTP_IDLE_TIMEOUT"))
	cfg.Http.ReadHeaderTimeout, _ = parseDuration(os.Getenv("HTTP_READ_HEADER_TIMEOUT"))
	cfg.Http.ShutdownTimeout, _ = parseDuration(os.Getenv("HTTP_SHUTDOWN_TIMEOUT"))

	// Log
	cfg.Log.Format = os.Getenv("LOG_FORMAT")
	cfg.Log.Level = os.Getenv("LOG_LEVEL")

	// DB
	cfg.DB.Host = os.Getenv("DB_HOST")
	cfg.DB.Port, _ = parseInt(os.Getenv("DB_PORT"))
	cfg.DB.User = os.Getenv("DB_USER")
	cfg.DB.Pwd = os.Getenv("DB_PWD")
	cfg.DB.Database = os.Getenv("DB_DATABASE")
	cfg.DB.SSL = os.Getenv("DB_SSL")
	cfg.DB.MaxOpenConnections, _ = parseInt(os.Getenv("DB_MAX_OPEN"))
	cfg.DB.MaxIdleConnections, _ = parseInt(os.Getenv("DB_MAX_IDLE"))
	cfg.DB.ConnMaxLifetime, _ = parseDuration(os.Getenv("DB_MAX_LIFETIME"))
	cfg.DB.InitTimeout, _ = parseDuration(os.Getenv("DB_INIT_TIMEOUT"))

	if err := val.ValidateStruct(&cfg); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &cfg, nil
}

func parseInt(s string) (int, error) {
	if s == "" {
		return 0, nil
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("invalid integer: %w", err)
	}
	return i, nil
}

func parseDuration(s string) (time.Duration, error) {
	if s == "" {
		return 0, nil
	}
	s = strings.TrimSpace(s)
	s = strings.ToLower(s)
	return time.ParseDuration(s)
}
