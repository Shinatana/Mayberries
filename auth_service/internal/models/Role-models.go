package models

type Role struct {
	ID          int    `json:"id" validate:"required"`
	Name        string `json:"name" validate:"required,min=3,max=255"`
	Description string `json:"description,omitempty" validate:"omitempty"`
}
