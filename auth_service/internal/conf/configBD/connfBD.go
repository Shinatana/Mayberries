package configBD

import (
	"auth_service/pkg/config"
	"database/sql"
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
