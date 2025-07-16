package init

import (
	internalConf "catalog_service/internal/conf"
	"catalog_service/internal/conf/loader"
	"catalog_service/pkg/log"
	"context"
	"fmt"
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

	redisClient, err := init_redis(&cfg.Redis)
	if err != nil {
		return err
	}
	defer func() {
		//я не уверена, как тут правильно обработать эту ошибку.
		//И не уверена тут ли вообще я должна ее обрабатывать
		redisClient.Close()
		log.Info("closed redis connection")
	}()
	log.Info("connected to redis")

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
