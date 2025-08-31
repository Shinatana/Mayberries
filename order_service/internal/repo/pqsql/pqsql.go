package pqsql

import (
	"context"
	"errors"
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
func (p *pgsql) FindOrder(ctx context.Context, id uuid.UUID) (models.Order, error) {
	var order models.Order
	result := p.pool.WithContext(ctx).First(&order, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return models.Order{}, models.ErrRecordNotFound
		}
		return models.Order{}, fmt.Errorf("failed to find order: %w", result.Error)
	}
	return order, nil
}

func (p *pgsql) DeleteOrder(ctx context.Context, id uuid.UUID) error {
	var order models.Order
	result := p.pool.WithContext(ctx).First(&order, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return models.ErrRecordNotFound
		}
		return fmt.Errorf("failed to find order: %w", result.Error)
	}
	if err := p.pool.WithContext(ctx).Delete(&order, id).Error; err != nil {
		return fmt.Errorf("failed to delete order: %w", err)
	}
	return nil
}

func (p *pgsql) PatchOrder(ctx context.Context, id uuid.UUID, status string) error {
	res := p.pool.WithContext(ctx).
		Model(&models.Order{}).
		Where("id = ?", id).
		Update("status", status)

	if res.Error != nil {
		return fmt.Errorf("failed to update order status: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return models.ErrRecordNotFound
	}
	return nil
}
