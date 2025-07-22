package init

import (
	authHttp "catalog_service/internal/http"
	"catalog_service/pkg/config"
	"catalog_service/pkg/log"
	"catalog_service/pkg/misc"

	"errors"
	"fmt"
	"net/http"
)

func Http(httpOptions *config.HttpOptions, handler http.Handler) (closer func()) {
	httpServer := authHttp.NewHttpServer(handler, httpOptions)

	go func() {
		err := httpServer.Start()
		if err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				log.Error("Failed to start server", "error", err)
				misc.GracefulStop()
			}
		}
	}()

	log.Debug(fmt.Sprintf("starting server on %s", httpServer.Addr()))
	log.Info("server started")

	return func() {
		err := httpServer.Close(httpOptions.ShutdownTimeout)
		if err != nil {
			log.Warn("failed to shutdown http server", "error", err)
		} else {
			log.Info("http server stopped")
		}
	}
}
