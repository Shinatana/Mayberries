package pgsql

import (
	"auth_service/internal/repo"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
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

	// вот тут вопрос! Я раньше делала это в misc.WithPGXv5Format. В Gorm это делать нельзя
	// вопрос: куда это лучше убрать, не думаю, что этим conf тут место

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

func (p *pgsql) RegisterUser(ctx context.Context, user models.RegisterUser) error {
	result := p.pool.WithContext(ctx).Create(&user)
	if result.Error != nil {
		var pgErr *pgconn.PgError
		if errors.As(result.Error, &pgErr) && pgErr.Code == "23505" {
			return models.ErrDuplicateUser
		}
		return fmt.Errorf("failed to register user: %w", result.Error)
	}

	return nil
}

func (p *pgsql) GetUserPassword(ctx context.Context, email string) (string, error) {
	var user models.RegisterUser
	err := p.pool.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", models.ErrUserNotFound
		}
		return "", fmt.Errorf("failed to get user's password: %w", err)
	}

	return user.Password, nil
}
