package pgsql

import (
	"auth_service/internal/repo"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"auth_service/internal/models"
	"auth_service/pkg/config"
	"auth_service/pkg/misc"
)

const initPingTimeout = 1 * time.Second
const registerUserQuery = `
	INSERT INTO users (email, name, password_hash)
	VALUES ($1, $2, $3);
`
const queryUserPwd = `
	SELECT password_hash FROM users WHERE email = $1;
`

type pgsql struct {
	pool *pgxpool.Pool
}

func NewDB(ctx context.Context, dbConfig *config.DatabaseOptions) (repo.DB, error) {
	config, err := pgxpool.ParseConfig(misc.GetDSN(dbConfig, misc.WithPGXv5Format()))
	if err != nil {
		return nil, fmt.Errorf("unable to parse connection string: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create pool: %w", err)
	}

	ctsSec, cancel := context.WithTimeout(ctx, initPingTimeout)
	defer cancel()

	err = pool.Ping(ctsSec)
	if err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping pool: %w", err)
	}

	return &pgsql{pool: pool}, nil
}

func (p *pgsql) Close() {
	p.pool.Close()
}

func (p *pgsql) RegisterUser(ctx context.Context, user models.RegisterUser) error {
	_, err := p.pool.Exec(ctx, registerUserQuery, user.Email, user.Name, user.Password)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return models.ErrDuplicateUser
		}
		return fmt.Errorf("failed to register user: %w", err)
	}

	return nil
}

func (p *pgsql) GetUserPassword(ctx context.Context, email string) (string, error) {
	var passwordHash string
	err := p.pool.QueryRow(ctx, queryUserPwd, email).Scan(&passwordHash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", models.ErrUserNotFound
		}
		return "", fmt.Errorf("failed to get user's password: %w", err)
	}

	return passwordHash, nil
}
