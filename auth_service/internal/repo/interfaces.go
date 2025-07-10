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
	RegisterUser(ctx context.Context, user models.Users) error
	GetUserPassword(ctx context.Context, email string) (string, error)
	GetUserRoles(ctx context.Context, userID string) ([]string, error)
	GetUserPermissions(ctx context.Context, userID string) ([]string, error)
	GetRoleByName(ctx context.Context, name string, role *models.Role) error
	CheckPermission(ctx context.Context, userID, permissionCode string) (bool, error)
	GetUserIDByEmail(ctx context.Context, email string) (string, error)
	Close()
}
