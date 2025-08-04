package pqsql

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/mayberries/shared/pkg/config"
	"github.com/mayberries/shared/pkg/misc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"order_service/internal/conf/configBD"
	"order_service/internal/models"
	"order_service/internal/repo"
	"strings"
	"time"
)

const initPingTimeout = 1 * time.Second

type pgsql struct {
	pool *gorm.DB
}

func NewDB(ctx context.Context, dbConfig *config.DatabaseOptions) (repo.DB, error) {
	dsn := misc.GetDSN(dbConfig, misc.WithGormFormat())

	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to open gorm DB: %w", err)
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB from gorm DB: %w", err)
	}

	configBD.СonfigureDBPool(sqlDB, dbConfig)

	ctsSec, cancel := context.WithTimeout(ctx, initPingTimeout)
	defer cancel()

	if err := sqlDB.PingContext(ctsSec); err != nil {
		return nil, fmt.Errorf("failed to ping DB: %w", err)
	}

	return &pgsql{pool: gormDB}, nil
}

func (p *pgsql) Close() {
	sqlDB, err := p.pool.DB()
	if err == nil {
		_ = sqlDB.Close()
	}
}

func (p *pgsql) CreateOrder(ctx context.Context, order models.Order) (uuid.UUID, error) {
	result := p.pool.WithContext(ctx).Create(&order)
	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "duplicate key value violates unique constraint") {
			return uuid.Nil, models.ErrDuplicateOrder
		}
		return uuid.Nil, fmt.Errorf("failed to create order: %w", result.Error)
	}
	return order.ID, nil
}
