package models

import (
	"errors"
)

var (
	ErrDuplicateOrder = errors.New("order already exists")
	ErrValidation     = errors.New("validation failed")
)
