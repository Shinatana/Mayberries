package init

import (
	"auth_service/internal/repo"
	"auth_service/internal/repo/pgsql"
	"auth_service/pkg/config"
	"auth_service/pkg/log"
	"auth_service/pkg/misc"
	"context"
	"database/sql"
	"fmt"
	"github.com/pressly/goose/v3"
)

func init_db(dbOptions *config.DatabaseOptions, dbMigrate *config.MigrationOptions) (repo.DB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbOptions.InitTimeout)
	defer cancel()

	sqlDB, err := sql.Open("postgres", misc.GetDSN(dbOptions, misc.WithMigratorFormat()))
	if err != nil {
		return nil, fmt.Errorf("failed to open sql DB for migrations: %w", err)
	}
	defer func() {
		if errSQL := sqlDB.Close(); errSQL != nil {

			log.Error("failed to close sqlDB", "error", errSQL)
		}
	}()

	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping DB for migrations: %w", err)
	}

	if err := goose.SetDialect("postgres"); err != nil {
		return nil, fmt.Errorf("failed to set goose dialect: %w", err)
	}

	if err := goose.Up(sqlDB, dbMigrate.MigrationFiles); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Info("✅ Migration complete")

	db, err := pgsql.NewDB(ctx, dbOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to create pgx DB pool: %w", err)
	}

	return db, nil
}
