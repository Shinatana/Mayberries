package models

import (
	"errors"
)

var (
	ErrDuplicateProducts = errors.New("handlers_products already exists")
	ErrProductsNotFound  = errors.New("handlers_products not found")
	ErrInvalidToken      = errors.New("invalid token")
	ErrInvalidTokenType  = errors.New("invalid token type")
	ErrFetchCategories   = errors.New("failed to fetch categories")
)
