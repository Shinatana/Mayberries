package models

import (
	"github.com/google/uuid"
	"time"
)

type Users struct {
	ID           uuid.UUID `json:"id" gorm:"primaryKey" `
	Email        string    `json:"email" gorm:"unique;not null" validator:"required,email"`
	Password     string    `json:"-" gorm:"-" validator:"required,min=12"`
	PasswordHash string    `json:"passwordHash" gorm:"column:password_hash;not null"`
	Name         string    `json:"name" gorm:"not null" validator:"required,min=3"`
	RoleID       int       `json:"roleId"`
	CreatedAt    time.Time `json:"createdAt" gorm:"default:CURRENT_TIMESTAMP"`
}
