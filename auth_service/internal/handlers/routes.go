package handlers

import (
	"auth_service/internal/config"
	"auth_service/internal/services"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router, cfg *config.Config, svc *services.AuthService) {
	r.Post("/register", RegisterHandler(svc))

	r.Post("/login", func(w http.ResponseWriter, r *http.Request) {
		LoginHandler(w, r, cfg)
	})

	// Пример публичного хелсчека
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})
}
