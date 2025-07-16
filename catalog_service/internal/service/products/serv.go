package products

import (
	"catalog_service/internal/models"
	"catalog_service/internal/repo"
	"context"
	"github.com/google/uuid"
)

type Service struct {
	DB repo.DB
}

func (s *Service) GetProducts(ctx context.Context) ([]models.Products, error) {
	return s.DB.GetProducts(ctx)
}

func (s *Service) GetProductById(ctx context.Context, id uuid.UUID) (*models.Products, error) {
	return s.DB.GetProductsById(ctx, id)
}

func (s *Service) CreateProduct(ctx context.Context, p models.Products) error {
	return s.DB.PostProducts(ctx, p)
}

func (s *Service) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	return s.DB.DeleteProducts(ctx, id)
}
