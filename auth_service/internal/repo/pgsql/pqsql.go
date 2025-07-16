package pgsql

import (
	"auth_service/internal/repo"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"strings"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"auth_service/internal/conf/configBD"
	"auth_service/internal/models"
	"auth_service/pkg/config"
	"auth_service/pkg/misc"
)

const initPingTimeout = 1 * time.Second

type pgsql struct {
	pool *gorm.DB
}

func NewDB(ctx context.Context, dbConfig *config.DatabaseOptions) (repo.DB, error) {
	dsn := misc.GetDSN(dbConfig, misc.WithGormFormat())

	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to open gorm DB: %w", err)
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB from gorm DB: %w", err)
	}

	configBD.СonfigureDBPool(sqlDB, dbConfig)

	ctsSec, cancel := context.WithTimeout(ctx, initPingTimeout)
	defer cancel()

	if err := sqlDB.PingContext(ctsSec); err != nil {
		return nil, fmt.Errorf("failed to ping DB: %w", err)
	}

	return &pgsql{pool: gormDB}, nil
}

func (p *pgsql) Close() {
	sqlDB, err := p.pool.DB()
	if err == nil {
		_ = sqlDB.Close()
	}
}

func (p *pgsql) RegisterUser(ctx context.Context, user models.Users) error {
	result := p.pool.WithContext(ctx).Create(&user)
	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "duplicate key value violates unique constraint") {
			return models.ErrDuplicateUser
		}
		return fmt.Errorf("failed to infoUser user: %w", result.Error)
	}
	return nil
}

func (p *pgsql) GetUserPassword(ctx context.Context, email string) (string, error) {
	var user models.Users
	err := p.pool.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", models.ErrUserNotFound
		}
		return "", fmt.Errorf("failed to get user's password: %w", err)
	}

	return user.PasswordHash, nil
}

func (p *pgsql) GetUserRoles(ctx context.Context, userID string) ([]string, error) {
	var roles []string
	err := p.pool.WithContext(ctx).
		Table("roles").
		Select("roles.name").
		Joins("JOIN users ON users.role_id = roles.id").
		Where("users.id = ?", userID).
		Pluck("name", &roles).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get user's roles: %w", err)
	}
	return roles, nil
}

func (p *pgsql) GetUserPermissions(ctx context.Context, userID string) ([]string, error) {
	var permissions []string
	err := p.pool.WithContext(ctx).
		Table("permissions").
		Select("permissions.code").
		Joins("JOIN roles_permissions ON roles_permissions.permission_id = permissions.id").
		Joins("JOIN roles ON roles.id = roles_permissions.role_id").
		Where("roles.id = (?)", p.pool.WithContext(ctx).
			Table("users").
			Select("role_id").
			Where("id = ?", userID),
		).
		Pluck("code", &permissions).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get user's permissions: %w", err)
	}
	return permissions, nil
}

func (p *pgsql) GetRoleByName(ctx context.Context, name string, role *models.Role) error {
	err := p.pool.WithContext(ctx).Where("name = ?", name).First(role).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.ErrRoleNotFound
		}
		return fmt.Errorf("failed to get role by name: %w", err)
	}
	return nil
}
func (p *pgsql) GetUserIDByEmail(ctx context.Context, email string) (string, error) {
	var user struct {
		ID string
	}
	err := p.pool.WithContext(ctx).
		Table("users").
		Select("id").
		Where("email = ?", email).
		First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", models.ErrUserNotFound
		}
		return "", fmt.Errorf("failed to get user ID by email: %w", err)
	}
	return user.ID, nil
}

func (p *pgsql) GetUserByID(ctx context.Context, userID string) (*models.Users, error) {
	var user models.Users

	err := p.pool.WithContext(ctx).
		Where("id = ?", userID).
		First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, models.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}
	return &user, nil
}

func (p *pgsql) GetUserByEmail(ctx context.Context, email string) (*models.Users, error) {
	var user models.Users

	err := p.pool.WithContext(ctx).
		Preload("Roles").
		Preload("Permissions").
		Where("email = ?", email).
		First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, models.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (p *pgsql) GetRoleByID(ctx context.Context, roleID int) (*models.Role, error) {
	var role models.Role

	err := p.pool.WithContext(ctx).
		Where("id = ?", roleID).
		First(&role).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, models.ErrRoleNotFound
		}
		return nil, err
	}
	return &role, nil
}

func (p *pgsql) ChangeDescriptionPermissions(ctx context.Context, permission *models.Permission) error {

	err := p.pool.WithContext(ctx).
		Model(&models.Permission{}).
		Where("id = ?", permission.ID).
		Update("description", permission.Description).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.ErrPermissionNotFound
		}
		return err
	}
	return nil
}

func (p *pgsql) CreateRole(ctx context.Context, role *models.Role) error {
	result := p.pool.WithContext(ctx).Create(role)
	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "duplicate key value violates unique constraint") {
			return models.ErrDuplicateUser
		}
		return fmt.Errorf("failed create role: %w", result.Error)
	}
	return nil
}

func (p *pgsql) DeleteRole(ctx context.Context, roleID int) error {
	result := p.pool.WithContext(ctx).Delete(&models.Role{
		ID: roleID,
	})
	if result.Error != nil {
		return fmt.Errorf("failed to deleteRole role: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("role not found")
	}
	return nil
}

func (p *pgsql) GetAllRoles(ctx context.Context) ([]models.Role, error) {
	var roles []models.Role
	err := p.pool.WithContext(ctx).Find(&roles).Error
	if err != nil {
		return nil, err
	}
	return roles, nil
}

func (p *pgsql) FindUsersByRole(ctx context.Context, roleName string) ([]string, error) {
	var userNames []string

	err := p.pool.WithContext(ctx).
		Model(&models.Users{}).
		Joins("JOIN roles ON roles.id = users.role_id").
		Where("roles.name = ?", roleName).
		Pluck("users.name", &userNames).Error

	if err != nil {
		return nil, err
	}
	return userNames, nil
}

func (p *pgsql) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	result := p.pool.WithContext(ctx).Delete(&models.Users{ID: userID})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return models.ErrUserNotFound
	}
	return nil
}
