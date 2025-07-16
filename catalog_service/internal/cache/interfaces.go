package cache

import (
	"catalog_service/internal/models"
	"context"
)

type Redis interface {
	GetCategories(ctx context.Context) ([]models.Categories, error)
}
