package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AccessClaims struct {
	jwt.RegisteredClaims
	TokenType   string   `json:"type"`
	Roles       []string `json:"roles,omitempty"`
	Permissions []string `json:"permissions,omitempty"`
}

type RefreshClaims struct {
	jwt.RegisteredClaims
	TokenType string `json:"type"`
}

type Handler interface {
	GenerateTokenPair(sub string, roles, permissions []string) (accessToken, refreshToken string, err error)
	VerifyRefreshToken(tokenString string) (*RefreshClaims, error)
	VerifyAccessToken(tokenString string) (*AccessClaims, error)
	GetTokenLifetime() time.Duration
	GetRefreshTokenLifetime() time.Duration
}
