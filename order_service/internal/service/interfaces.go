package service

import (
	"context"
	"github.com/google/uuid"
	"order_service/internal/models"
)

type Service interface {
	CreateOrder(ctx context.Context, order models.Order) (uuid.UUID, error)
	GetOrder(ctx context.Context, id uuid.UUID) (models.Order, error)
	DeleteOrder(ctx context.Context, id uuid.UUID) error
	PatchOrder(ctx context.Context, id uuid.UUID, order models.Order) error
}
