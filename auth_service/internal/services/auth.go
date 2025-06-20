package services

import (
	"auth_service/internal/storage"
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	store storage.UserStore
}

func NewAuthService(store storage.UserStore) *AuthService {
	return &AuthService{store: store}
}

func (s *AuthService) Register(ctx context.Context, email, password, fullName, role string) error {
	if email == "" || password == "" || fullName == "" || role == "" {
		return errors.New("all fields are required")
	}

	exists, err := s.store.Exists(ctx, email)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("user already exists!")
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := storage.User{
		Email:    email,
		Password: string(hashedPass),
		FullName: fullName,
		Role:     role,
	}

	return s.store.Create(ctx, &user)
}
