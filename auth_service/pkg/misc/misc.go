package misc

import (
	"auth_service/pkg/config"
	"auth_service/pkg/log"
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"syscall"

	"github.com/google/uuid"
)

func GracefulStop() {
	_ = syscall.Kill(os.Getpid(), syscall.SIGTERM)
}

func WaitForSignal(sigs ...os.Signal) {
	quit := make(chan os.Signal, 1)
	defer close(quit)
	signal.Notify(quit, sigs...)
	defer signal.Stop(quit)
	<-quit
	log.Info("server stopping...")

}

func GenerateUUID() string {
	return uuid.New().String()
}

type GetDSNFormat func(dsn *url.URL, db *config.DatabaseOptions) string

func WithMigratorFormat() GetDSNFormat {
	return func(dsn *url.URL, db *config.DatabaseOptions) string {
		dsn.Scheme = "pgx5"

		query := dsn.Query()

		query.Set("sslmode", db.SSL)
		if db.Schema != "" {
			query.Set("search_path", db.Schema)
		}

		dsn.RawQuery = query.Encode()

		return dsn.String()
	}
}

func WithPGXv5Format() GetDSNFormat {
	return func(dsn *url.URL, db *config.DatabaseOptions) string {
		dsn.Scheme = "postgres"

		query := dsn.Query()

		query.Set("sslmode", db.SSL)
		if db.Schema != "" {
			query.Set("search_path", db.Schema)
		}

		if db.MaxOpenConnections > 0 {
			query.Set("pool_max_conns", fmt.Sprintf("%d", db.MaxOpenConnections))
		}

		if db.MaxIdleConnections > 0 {
			query.Set("pool_min_conns", fmt.Sprintf("%d", db.MaxIdleConnections))
		}

		if db.ConnMaxLifetime > 0 {
			query.Set("pool_max_conn_lifetime", db.ConnMaxLifetime.String())
		}

		dsn.RawQuery = query.Encode()

		return dsn.String()
	}
}

func GetDSN(db *config.DatabaseOptions, format GetDSNFormat) string {
	dsn := url.URL{
		User: url.UserPassword(db.User, db.Pwd),
		Host: fmt.Sprintf("%s:%d", db.Host, db.Port),
		Path: db.Database,
	}

	return format(&dsn, db)
}
