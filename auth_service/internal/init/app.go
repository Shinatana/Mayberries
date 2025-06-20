package init

import (
	"auth_service/internal/conf"
	"auth_service/internal/conf/loader"
	"auth_service/pkg/log"
	"fmt"
)

func App() error {
	cfg, err := loader.NewLoader().Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	Log(&cfg.Log)

	sanitizeConfig := func(cfg conf.Config) *conf.Config {
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

	return nil
}
