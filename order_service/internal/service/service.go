package service

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/mayberries/shared/pkg/val"
	"order_service/internal/models"
	"order_service/internal/repo"
)

type service struct {
	db repo.DB
}

func NewService(db repo.DB) Service {
	return &service{db: db}
}

func (s *service) CreateOrder(ctx context.Context, o models.Order) (uuid.UUID, error) {

	if err := val.ValidateStruct(o); err != nil {
		return uuid.Nil, fmt.Errorf("%w: %v", models.ErrValidation, err)
	}

	id, err := s.db.CreateOrder(ctx, o)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func (s *service) GetOrder(ctx context.Context, id uuid.UUID) (models.Order, error) {
	var order models.Order
	if id == uuid.Nil {
		return models.Order{}, models.ErrIDIsNil
	}
	order, err := s.db.FindOrder(ctx, id)
	if err != nil {
		return models.Order{}, err
	}
	return order, nil
}

func (s *service) DeleteOrder(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return models.ErrIDIsNil
	}
	err := s.db.DeleteOrder(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) PatchOrder(ctx context.Context, id uuid.UUID, order models.Order) error {

	if id == uuid.Nil {
		return models.ErrIDIsNil
	}

	if order.Status == "" {
		if err := val.ValidateStruct(order.Status); err != nil {
			return fmt.Errorf("%w: %v", models.ErrValidation, err)
		}
		return fmt.Errorf("status is empty")
	}

	err := s.db.PatchOrder(ctx, id, order.Status)
	if err != nil {
		return err
	}
	return nil
}
