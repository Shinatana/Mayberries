package repo

import (
	"catalog_service/internal/models"
	"context"
	"github.com/google/uuid"
)

// Migrator defines an abstraction for running migrations.
type Migrator interface {
	Migrate(ctx context.Context) error
}

type DB interface {
	GetProducts(ctx context.Context) ([]models.Products, error)
	GetProductsById(ctx context.Context, productId uuid.UUID) (*models.Products, error)
	PostProducts(ctx context.Context, products models.Products) error
	DeleteProducts(ctx context.Context, productId uuid.UUID) error
	Close()
}
