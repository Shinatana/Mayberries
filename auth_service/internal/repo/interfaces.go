package repo

import (
	"context"
	"github.com/google/uuid"

	"auth_service/internal/models"
)

type DB interface {
	RegisterUser(ctx context.Context, user models.Users) error
	GetUserPassword(ctx context.Context, email string) (string, error)
	GetUserRoles(ctx context.Context, userID string) ([]string, error)
	GetUserByID(ctx context.Context, userID string) (*models.Users, error)
	GetUserPermissions(ctx context.Context, userID string) ([]string, error)
	ChangeDescriptionPermissions(ctx context.Context, permission *models.Permission) error
	GetRoleByName(ctx context.Context, name string, role *models.Role) error
	GetRoleByID(ctx context.Context, roleID int) (*models.Role, error)
	GetUserIDByEmail(ctx context.Context, email string) (string, error)
	CreateRole(ctx context.Context, role *models.Role) error
	DeleteRole(ctx context.Context, role int) error
	GetAllRoles(ctx context.Context) ([]models.Role, error)
	DeleteUser(ctx context.Context, userID uuid.UUID) error
	FindUsersByRole(ctx context.Context, roleName string) ([]string, error)
	Close()
}
