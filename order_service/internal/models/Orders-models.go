package models

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID            uuid.UUID `json:"id" db:"id"`
	UserID        uuid.UUID `json:"user_id" db:"user_id" validate:"required"`
	TotalPrice    float64   `json:"total_price" db:"total_price" validate:"required,gte=0"`
	DeliveryPrice float64   `json:"delivery_price" db:"delivery_price" validate:"required,gte=0"`
	Currency      string    `json:"currency" db:"currency" validate:"required,len=3"`
	Status        string    `json:"status" db:"status" validate:"required,oneof=created paid shipped delivered cancelled"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
}
