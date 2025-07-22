package repo

import (
	"catalog_service/internal/models"
	"context"
	"github.com/google/uuid"
)

type DB interface {
	GetProducts(ctx context.Context) ([]models.Products, error)
	GetProductsById(ctx context.Context, productId uuid.UUID) (*models.Products, error)
	PostProducts(ctx context.Context, products models.Products) error
	DeleteProducts(ctx context.Context, productId uuid.UUID) error
	GetCategories(ctx context.Context) ([]models.Categories, error)
	Close()
}
