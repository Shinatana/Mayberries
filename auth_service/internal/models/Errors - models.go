package models

import (
	"errors"
)

var (
	ErrDuplicateUser    = errors.New("user already exists")
	ErrUserNotFound     = errors.New("user not found")
	ErrInvalidToken     = errors.New("invalid token")
	ErrInvalidTokenType = errors.New("invalid token type")
	ErrRoleNotFound     = errors.New("role not found")
)
