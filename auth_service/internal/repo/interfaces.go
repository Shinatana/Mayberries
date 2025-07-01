package repo

import (
	"context"

	"auth_service/internal/models"
)

// Migrator defines an abstraction for running migrations.
type Migrator interface {
	Migrate(ctx context.Context) error
}

type DB interface {
	RegisterUser(ctx context.Context, user models.RegisterUser) error
	GetUserPassword(ctx context.Context, email string) (string, error)
	Close()
}
