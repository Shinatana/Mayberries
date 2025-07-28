package configBD

import (
	"database/sql"
	"github.com/mayberries/shared/pkg/config"
)

func СonfigureDBPool(sqlDB *sql.DB, dbConfig *config.DatabaseOptions) {
	if dbConfig.MaxOpenConnections > 0 {
		sqlDB.SetMaxOpenConns(dbConfig.MaxOpenConnections)
	}
	if dbConfig.MaxIdleConnections >= 0 {
		sqlDB.SetMaxIdleConns(dbConfig.MaxIdleConnections)
	}
	if dbConfig.ConnMaxLifetime > 0 {
		sqlDB.SetConnMaxLifetime(dbConfig.ConnMaxLifetime)
	}
}
