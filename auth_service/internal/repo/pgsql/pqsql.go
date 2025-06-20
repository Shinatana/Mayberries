package pgsql

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"auth_service/pkg/config"
	"auth_service/pkg/misc"
)

const initPingTimeout = 1 * time.Second

type Pgsql struct {
	pool *pgxpool.Pool
}

func NewDB(ctx context.Context, dbConfig *config.DatabaseOptions) (*Pgsql, error) {
	config, err := pgxpool.ParseConfig(misc.GetDSN(dbConfig, misc.WithPGXv5Format()))
	if err != nil {
		return nil, fmt.Errorf("unable to parse connection string: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create pool: %w", err)
	}

	ctsSec, cancel := context.WithTimeout(ctx, initPingTimeout)
	defer cancel()

	if err = pool.Ping(ctsSec); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping pool: %w", err)
	}

	return &Pgsql{pool: pool}, nil
}

func (p *Pgsql) Close() {
	p.pool.Close()
}
