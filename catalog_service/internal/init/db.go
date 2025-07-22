package init

import (
	"catalog_service/internal/repo"
	"catalog_service/internal/repo/pqsql"
	"catalog_service/pkg/config"
	"catalog_service/pkg/log"
	"catalog_service/pkg/misc"
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

	db, err := pqsql.NewDB(ctx, dbOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to create pgx DB pool: %w", err)
	}

	return db, nil
}
