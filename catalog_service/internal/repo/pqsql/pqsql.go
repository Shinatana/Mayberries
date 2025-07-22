package pqsql

import (
	"catalog_service/internal/conf/configBD"
	"catalog_service/internal/models"
	"catalog_service/internal/repo"
	"catalog_service/pkg/config"
	"catalog_service/pkg/misc"
	"context"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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

func (p *pgsql) GetProducts(ctx context.Context) ([]models.Products, error) {
	var products []models.Products

	err := p.pool.WithContext(ctx).
		Model(&models.Products{}).
		Find(&products).
		Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch handlers_products: %w", err)
	}
	return products, nil
}

func (p *pgsql) PostProducts(ctx context.Context, products models.Products) error {
	result := p.pool.WithContext(ctx).Create(&products)
	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "duplicate key value violates unique constraint") {
			return models.ErrDuplicateProducts
		}
		return fmt.Errorf("failed to infoUser handlers_products: %w", result.Error)
	}
	return nil
}

func (p *pgsql) GetProductsById(ctx context.Context, productId uuid.UUID) (*models.Products, error) {
	var product models.Products
	err := p.pool.WithContext(ctx).
		First(&product, "id = ?", productId).
		Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch product by id: %w", err)
	}
	return &product, nil
}

func (p *pgsql) DeleteProducts(ctx context.Context, productId uuid.UUID) error {
	result := p.pool.WithContext(ctx).Delete(&models.Products{
		ID: productId,
	})
	if result.Error != nil {
		return fmt.Errorf("failed to deleteRole handlers_products: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("product not found")
	}
	return nil
}

func (p *pgsql) GetCategories(ctx context.Context) ([]models.Categories, error) {
	var categories []models.Categories

	err := p.pool.WithContext(ctx).
		Model(&models.Categories{}).
		Find(&categories).
		Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch categories: %w", err)
	}

	return categories, nil
}
