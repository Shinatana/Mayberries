package models

import (
	"errors"
)

var (
	ErrDuplicateProducts = errors.New("products already exists")
	ErrProductsNotFound  = errors.New("products not found")
	ErrInvalidToken      = errors.New("invalid token")
	ErrInvalidTokenType  = errors.New("invalid token type")
)
