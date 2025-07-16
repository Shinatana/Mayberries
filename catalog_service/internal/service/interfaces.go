package service

import (
	"catalog_service/internal/models"
	"context"
	"github.com/google/uuid"
)

type Service interface {
	GetProducts(ctx context.Context) ([]models.Products, error)
	GetProductById(ctx context.Context, id uuid.UUID) (*models.Products, error)
	CreateProduct(ctx context.Context, p models.Products) error
	DeleteProduct(ctx context.Context, id uuid.UUID) error
	GetCategories(ctx context.Context) ([]models.Categories, error)
}
