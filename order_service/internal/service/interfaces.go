package service

import (
	"context"
	"github.com/google/uuid"
	"order_service/internal/models"
)

type Service interface {
	CreateOrder(ctx context.Context, order models.Order) (uuid.UUID, error)
}
