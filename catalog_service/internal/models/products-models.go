package models

import (
	"github.com/google/uuid"
	"time"
)

type Products struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name        string    `json:"name" validate:"required,min=2,max=255" gorm:"type:varchar(255);not null"`
	Description string    `json:"description" validate:"max=1000" gorm:"type:text"` // лимит на длину, если хочешь
	Price       float64   `json:"price" validate:"required,gt=0" gorm:"type:decimal(10,2);not null"`
	CategoryID  int       `json:"categoryId" validate:"required,gt=0" gorm:"not null"`
	CreatedAt   time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
}
