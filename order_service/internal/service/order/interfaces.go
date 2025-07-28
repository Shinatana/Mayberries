package order

import (
	"context"
	"order_service/internal/models"
)

type Service interface {
	CreateOrder(ctx context.Context, order models.Order) error
}
