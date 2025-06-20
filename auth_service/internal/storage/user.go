package storage

import (
	"context"
	"database/sql"
)

type User struct {
	ID       int64
	Email    string
	Password string
	FullName string
	Role     string
}

type UserStore interface {
	Exists(ctx context.Context, email string) (bool, error)
	Create(ctx context.Context, user *User) error
}

type PostgresUserStore struct {
	db *sql.DB
}

func NewPostgresUserStore(db *sql.DB) *PostgresUserStore {
	return &PostgresUserStore{db: db}
}

func (s *PostgresUserStore) Exists(ctx context.Context, email string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email=$1)`
	err := s.db.QueryRowContext(ctx, query, email).Scan(&exists)
	return exists, err
}

func (s *PostgresUserStore) Create(ctx context.Context, user *User) error {
	query := `INSERT INTO users (email, password, full_name, role) VALUES ($1, $2, $3, $4)`
	_, err := s.db.ExecContext(ctx, query, user.Email, user.Password, user.FullName, user.Role)
	return err
}
