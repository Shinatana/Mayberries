package models

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	Password     string    `json:"-"`
	PasswordHash string    `json:"passwordHash"`
	Name         string    `json:"name"`
	RoleID       int       `json:"roleId"`
	CreatedAt    time.Time `json:"createdAt"`
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
