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
)

func init_db(dbOptions *config.DatabaseOptions) (repo.DB, error) {
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

	tx, err := sqlDB.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin migration transaction: %w", err)
	}

	log.Info("⏫ Starting migration")

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit migration transaction: %w", err)
	}

	log.Info("✅ Migration complete")

	db, err := pgsql.NewDB(ctx, dbOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to create pgx DB pool: %w", err)
	}

	return db, nil
}
