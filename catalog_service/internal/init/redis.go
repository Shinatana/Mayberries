package init

import (
	"catalog_service/pkg/config"
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

func init_redis(cfg *config.RedisOptions) (*redis.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), cfg.InitTimeout)
	defer cancel()

	addrRedis := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	rdb := redis.NewClient(&redis.Options{
		Addr:     addrRedis,
		Password: cfg.Password,
		DB:       cfg.DB,
	})
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to ping redis: %w", err)
	}

	return rdb, nil
}
