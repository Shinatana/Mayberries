package models

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID            uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserID        uuid.UUID `json:"user_id" gorm:"type:uuid;not null" validate:"required"`
	TotalPrice    float64   `json:"total_price" gorm:"not null" validate:"required,gte=0"`
	DeliveryPrice float64   `json:"delivery_price" gorm:"not null" validate:"required,gte=0"`
	Currency      string    `json:"currency" gorm:"type:char(3);not null" validate:"required,len=3"`
	Status        string    `json:"status" gorm:"type:varchar(20);not null" validate:"required,oneof=created paid shipped delivered cancelled"`
	CreatedAt     time.Time `json:"created_at" gorm:"autoCreateTime"`
}
