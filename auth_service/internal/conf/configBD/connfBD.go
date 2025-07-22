package configBD

import (
	"auth_service/pkg/config"
	"database/sql"
)

func СonfigureDBPool(sqlDb *sql.DB, dbConfig *config.DatabaseOptions) {
	if dbConfig.MaxOpenConnections > 0 {
		sqlDb.SetMaxOpenConns(dbConfig.MaxOpenConnections)
	}
	if dbConfig.MaxIdleConnections >= 0 {
		sqlDb.SetMaxIdleConns(dbConfig.MaxIdleConnections)
	}
	if dbConfig.ConnMaxLifetime > 0 {
		sqlDb.SetConnMaxLifetime(dbConfig.ConnMaxLifetime)
	}
}
