package http

import (
	"auth_service/pkg/config"
	"context"
	"net/http"
	"strconv"
	"time"
)

type Server struct {
	server *http.Server
}

func NewHttpServer(handler http.Handler, config *config.HttpOptions) *Server {
	return &Server{
		server: &http.Server{
			Addr:              config.Host + ":" + strconv.Itoa(config.Port),
			Handler:           handler,
			ReadTimeout:       config.ReadTimeout,
			ReadHeaderTimeout: config.ReadHeaderTimeout,
			WriteTimeout:      config.WriteTimeout,
			IdleTimeout:       config.IdleTimeout,
			MaxHeaderBytes:    config.MaxHeaderBytes,
		},
	}
}

func (s *Server) Start() error {
	return s.server.ListenAndServe()
}

func (s *Server) Addr() string {
	return s.server.Addr
}

func (s *Server) Close(t time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), t)
	defer cancel()

	return s.server.Shutdown(ctx)
}
