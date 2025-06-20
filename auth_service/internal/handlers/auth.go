package handlers

import (
	"auth_service/internal/config"
	"auth_service/internal/services"
	"encoding/json"
	"net/http"
)

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	FullName string `json:"full_name"`
	Role     string `json:"role"`
}

func RegisterHandler(svc *services.AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RegisterRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err := svc.Register(r.Context(), req.Email, req.Password, req.FullName, req.Role)
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"message":"user registered successfully"}`))
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request, cfg *config.Config) {
	w.Write([]byte("login endpoint\n"))
}
