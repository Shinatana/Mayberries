package init

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/mayberries/shared/pkg/config"
	"github.com/mayberries/shared/pkg/log"
	"github.com/mayberries/shared/pkg/misc"
	"order_service/internal/repo"
	"order_service/internal/repo/pqsql"
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
