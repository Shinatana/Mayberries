package order

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/mayberries/shared/pkg/val"
	"order_service/internal/models"
	"order_service/internal/repo"
)

const defaultStatus = "created"

type service struct {
	db repo.DB
}

func NewService(db repo.DB) Service {
	return &service{db: db}
}

func (s *service) CreateOrder(ctx context.Context, o models.Order) error {
	o.ID = uuid.New()
	o.Status = defaultStatus

	if err := val.ValidateStruct(o); err != nil {
		return fmt.Errorf("%w: %v", models.ErrValidation, err)
	}

	return s.db.CreateOrder(ctx, o)
}
