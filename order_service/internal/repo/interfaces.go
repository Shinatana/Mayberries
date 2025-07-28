package repo

import (
	"context"
	"order_service/internal/models"
)

type DB interface {
	CreateOrder(ctx context.Context, order models.Order) error
	Close()
}
