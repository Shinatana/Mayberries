package pgsql

import (
	"auth_service/internal/models"
	"auth_service/pkg/log"
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
)

type userRepository struct {
	db *Pgsql
}

// NewUserRepository — конструктор для репозитория пользователей PostgreSQL
func NewUserRepository(db *Pgsql) *userRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	query := `
        INSERT INTO users (id, email, password_hash, name, role_id, created_at)
        VALUES ($1, $2, $3, $4, $5, $6)
    `
	_, err := r.db.pool.Exec(ctx, query,
		user.ID,
		user.Email,
		user.PasswordHash,
		user.Name,
		user.RoleID,
		user.CreatedAt,
	)
	if err != nil {
		log.Error("failed to create user", "email", user.Email, "error", err)
		return err
	}

	log.Info("user created successfully", "email", user.Email, "id", user.ID)
	return nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
        SELECT id, email, password_hash, name, role_id, created_at
        FROM users WHERE email = $1
    `
	row := r.db.pool.QueryRow(ctx, query, email)

	var user models.User
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Name,
		&user.RoleID,
		&user.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Info("user not found by email", "email", email)
			return nil, nil // пользователь не найден
		}
		log.Error("failed to get user by email", "email", email, "error", err)
		return nil, err
	}

	log.Info("user retrieved by email", "email", email, "id", user.ID)
	return &user, nil
}

func (r *userRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
	query := `
        SELECT id, email, password_hash, name, role_id, created_at
        FROM users WHERE id = $1
    `
	row := r.db.pool.QueryRow(ctx, query, id)

	var user models.User
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Name,
		&user.RoleID,
		&user.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Info("user not found by id", "id", id)
			return nil, nil
		}
		log.Error("failed to get user by id", "id", id, "error", err)
		return nil, err
	}

	log.Info("user retrieved by id", "id", id, "email", user.Email)
	return &user, nil
}
