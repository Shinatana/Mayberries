package config

import (
	"errors"
	"os"
	"strconv"
)

type Config struct {
	Port      string
	DBUrl     string
	RedisAddr string
	JWTSecret string
	TokenTTL  int // в минутах
}

func Load() (*Config, error) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dburl := os.Getenv("DB_URL")
	if dburl == "" {
		return nil, errors.New("DB_URL env var is required")
	}

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return nil, errors.New("JWT_SECRET env var is required")
	}

	tokenTTLStr := os.Getenv("TOKEN_TTL")
	tokenTTL := 60 // по умолчанию 60 минут
	if tokenTTLStr != "" {
		if ttl, err := strconv.Atoi(tokenTTLStr); err == nil {
			tokenTTL = ttl
		}
	}

	return &Config{
		Port:      port,
		DBUrl:     dburl,
		RedisAddr: redisAddr,
		JWTSecret: jwtSecret,
		TokenTTL:  tokenTTL,
	}, nil
}
