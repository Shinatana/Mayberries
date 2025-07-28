package init

import (
	"context"
	"fmt"
	"github.com/mayberries/shared/pkg/log"
	internalConf "order_service/internal/conf"
	"order_service/internal/conf/loader"
	"os/signal"
	"syscall"
)

func App() error {
	cfg, err := loader.NewLoader().Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	Log(&cfg.Log)

	sanitizeConfig := func(cfg internalConf.Config) *internalConf.Config {
		sanitized := cfg

		// if sanitized.Redis.Pwd != "" {
		// 	sanitized.Redis.Pwd = "********"
		// }
		// if sanitized.init_db.Pwd != "" {
		// 	sanitized.init_db.Pwd = "********"
		// }

		return &sanitized
	}
	log.Debug(fmt.Sprintf("config: %+v", sanitizeConfig(*cfg)))

	db, err := init_db(&cfg.DB)
	if err != nil {
		return err
	}

	defer func() {
		db.Close()
		log.Info("closed database connection")
	}()
	log.Info("connected to database")

	httpClose := Http(
		&cfg.Http,
		Gin(db),
	)

	defer httpClose()

	// Graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()

	return nil
}
