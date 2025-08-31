package repo

import (
	"context"
	"github.com/google/uuid"
	"order_service/internal/models"
)

type DB interface {
	CreateOrder(ctx context.Context, order models.Order) (uuid.UUID, error)
	FindOrder(ctx context.Context, id uuid.UUID) (models.Order, error)
	DeleteOrder(ctx context.Context, id uuid.UUID) error
	PatchOrder(ctx context.Context, id uuid.UUID, status string) error
	Close()
}
