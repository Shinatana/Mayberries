package categories

import (
	"catalog_service/internal/models"
	"catalog_service/internal/repo"
	"context"
	"github.com/redis/go-redis/v9"
)

type Service struct {
	db    repo.DB
	redis *redis.Client
}

func (s *Service) GetCategories(ctx context.Context) ([]models.Products, error) {
	return s.DB.GetCategories(ctx)
}
