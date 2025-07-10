package pgsql

import (
	"auth_service/internal/repo"
	"context"
	"errors"
	"fmt"
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

	sqlDB, err := gormDB.DB() // Получаем sql.DB для низкоуровневых операций, например ping
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
		return fmt.Errorf("failed to info user: %w", result.Error)
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

func (p *pgsql) CheckPermission(ctx context.Context, userID, permissionCode string) (bool, error) {
	var count int64
	err := p.pool.WithContext(ctx).
		Table("permissions").
		Joins("JOIN roles_permissions ON roles_permissions.permission_id = permissions.id").
		Joins("JOIN roles ON roles.id = roles_permissions.role_id").
		Joins("JOIN users ON users.role_id = roles.id").
		Where("users.id = ? AND permissions.code = ?", userID, permissionCode).
		Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("failed to check permission: %w", err)
	}
	return count > 0, nil
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
