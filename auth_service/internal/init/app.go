package init

import (
	internalConf "auth_service/internal/conf"
	"auth_service/internal/conf/loader"
	hashImpl "auth_service/internal/hash/hasher"
	"auth_service/internal/http/gin"
	"auth_service/internal/http/gin/middlewares/recovery"
	requestid "auth_service/internal/http/gin/middlewares/request-id"
	"auth_service/internal/http/gin/routes/auth/login"
	"auth_service/internal/http/gin/routes/v1/auth/register"
	"auth_service/internal/http/gin/routes/v2/auth/refresh"
	"auth_service/internal/http/gin/routes/v3/auth/info"
	gojwt "auth_service/internal/jwt/go-jwt"
	"auth_service/pkg/log"
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

	jwtHandler, err := gojwt.NewJwtHandler(&cfg.JWT)
	if err != nil {
		return err
	}
	log.Info("jwt keys loaded")

	db, err := init_db(&cfg.DB)
	if err != nil {
		return err
	}

	defer func() {
		db.Close()
		log.Info("closed database connection")
	}()
	log.Info("connected to database")

	ginServer := gin.NewGinServer()

	ginServer.AddMiddleware(
		recovery.Middleware(),
		requestid.Middleware(),
	)

	hasher := hashImpl.NewHasher(hashImpl.MinHashCost)

	ginServer.AddRouters(
		login.NewLoginHandler(db, hasher, jwtHandler),
		register.NewRegisterHandler(db, hasher),
		refresh.NewRefreshHandler(db, jwtHandler),
		info.TakeInfoMe(db, jwtHandler),
	)

	handler := ginServer.Build()

	closer := Http(&cfg.Http, handler)

	// Graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()

	closer() // Shutdown сервер
	return nil
}
