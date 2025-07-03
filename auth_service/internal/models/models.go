package models

import (
	"errors"
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID           uuid.UUID `json:"id" gorm:"primaryKey"`
	Email        string    `json:"email" gorm:"unique;not null"`
	Password     string    `json:"-" gorm:"-"`
	PasswordHash string    `json:"passwordHash" gorm:"column:password_hash;not null"`
	Name         string    `json:"name" gorm:"not null"`
	RoleID       int       `json:"roleId"`
	CreatedAt    time.Time `json:"createdAt" gorm:"default:CURRENT_TIMESTAMP"`
}

type Role struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

type Permission struct {
	ID          int    `json:"id"`
	Code        string `json:"code"`
	Description string `json:"description"`
}

type RolePermission struct {
	RoleID       int `json:"roleId"`
	PermissionID int `json:"permissionId"`
}

var (
	ErrDuplicateUser    = errors.New("user already exists")
	ErrUserNotFound     = errors.New("user not found")
	ErrInvalidToken     = errors.New("invalid token")
	ErrInvalidTokenType = errors.New("invalid token type")
)

type RegisterUser struct {
	Email    string `json:"email" validator:"required,email"`
	Password string `json:"password" validator:"required,min=12"`
	Name     string `json:"name" validator:"required,min=3"`
}
