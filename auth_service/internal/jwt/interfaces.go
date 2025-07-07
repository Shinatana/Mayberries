package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AccessClaims struct {
	jwt.RegisteredClaims
	TokenType string `json:"type"`
}

type RefreshClaims struct {
	jwt.RegisteredClaims
	TokenType string `json:"type"`
}

type Handler interface {
	GenerateTokenPair(sub string) (accessToken, refreshToken string, err error)
	VerifyRefreshToken(tokenString string) (*RefreshClaims, error)
	VerifyAccessToken(tokenString string) (*AccessClaims, error)
	GetTokenLifetime() time.Duration
	GetRefreshTokenLifetime() time.Duration
}
