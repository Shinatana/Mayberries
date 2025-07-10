package gojwt

import (
	"crypto/ed25519"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"

	localjwt "auth_service/internal/jwt"
	"auth_service/internal/models"
	"auth_service/pkg/config"
	"auth_service/pkg/misc"
)

const (
	leeway           = 30 * time.Second
	TokenTypeAccess  = "access"
	TokenTypeRefresh = "refresh"
)

type jwtHandler struct {
	ed25519key           ed25519.PrivateKey
	ed25519pub           ed25519.PublicKey
	tokenLifetime        time.Duration
	refreshTokenLifetime time.Duration
	issuer               string
}

func NewJwtHandler(cfg *config.JwtOptions) (localjwt.Handler, error) {
	keyData, err := os.ReadFile(cfg.ED25519KeyFile)
	if err != nil {
		return nil, fmt.Errorf("error reading private key file: %w", err)
	}

	pubData, err := os.ReadFile(cfg.ED25519PubFile)
	if err != nil {
		return nil, fmt.Errorf("error reading public key file: %w", err)
	}

	keyBlock, _ := pem.Decode(keyData)
	if keyBlock == nil {
		return nil, fmt.Errorf("failed to parse PEM block containing the private key")
	}

	pubBlock, _ := pem.Decode(pubData)
	if pubBlock == nil {
		return nil, fmt.Errorf("failed to parse PEM block containing the public key")
	}

	privateKey, err := x509.ParsePKCS8PrivateKey(keyBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	publicKey, err := x509.ParsePKIXPublicKey(pubBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}

	ed25519Key, ok := privateKey.(ed25519.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("private key is not an ED25519 key")
	}

	ed25519Pub, ok := publicKey.(ed25519.PublicKey)
	if !ok {
		return nil, fmt.Errorf("public key is not an ED25519 key")
	}

	return &jwtHandler{
		ed25519key:           ed25519Key,
		ed25519pub:           ed25519Pub,
		tokenLifetime:        cfg.TokenLifetime,
		refreshTokenLifetime: cfg.RefreshTokenLifetime,
		issuer:               cfg.Issuer,
	}, nil
}

func (j *jwtHandler) GenerateTokenPair(sub string, roles, permissions []string) (accessToken, refreshToken string, err error) {
	accessTokenID := misc.GenerateUUID()
	refreshTokenID := misc.GenerateUUID()
	timeNow := time.Now()

	accessClaims := localjwt.AccessClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   sub,
			Issuer:    j.issuer,
			IssuedAt:  jwt.NewNumericDate(timeNow),
			ExpiresAt: jwt.NewNumericDate(timeNow.Add(j.tokenLifetime)),
			NotBefore: jwt.NewNumericDate(timeNow),
			ID:        accessTokenID,
		},
		TokenType:   TokenTypeAccess,
		Roles:       roles,
		Permissions: permissions,
	}

	refreshClaims := localjwt.RefreshClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   sub,
			Issuer:    j.issuer,
			IssuedAt:  jwt.NewNumericDate(timeNow),
			ExpiresAt: jwt.NewNumericDate(timeNow.Add(j.refreshTokenLifetime)),
			NotBefore: jwt.NewNumericDate(timeNow),
			ID:        refreshTokenID,
		},
		TokenType: TokenTypeRefresh,
	}

	accessJWT := jwt.NewWithClaims(jwt.SigningMethodEdDSA, accessClaims)
	refreshJWT := jwt.NewWithClaims(jwt.SigningMethodEdDSA, refreshClaims)

	accessToken, err = accessJWT.SignedString(j.ed25519key)
	if err != nil {
		return "", "", err
	}

	refreshToken, err = refreshJWT.SignedString(j.ed25519key)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (j *jwtHandler) GetTokenLifetime() time.Duration {
	return j.tokenLifetime
}

func (j *jwtHandler) GetRefreshTokenLifetime() time.Duration {
	return j.refreshTokenLifetime
}

func (j *jwtHandler) VerifyAccessToken(tokenString string) (*localjwt.AccessClaims, error) {
	var parseOptions []jwt.ParserOption

	parseOptions = append(parseOptions, jwt.WithValidMethods([]string{"EdDSA"}))
	parseOptions = append(parseOptions, jwt.WithIssuer(j.issuer))
	parseOptions = append(parseOptions, jwt.WithLeeway(leeway))
	parseOptions = append(parseOptions, jwt.WithExpirationRequired())
	parseOptions = append(parseOptions, jwt.WithIssuedAt())

	token, err := jwt.ParseWithClaims(
		tokenString,
		&localjwt.AccessClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return j.ed25519pub, nil
		},
		parseOptions...,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	// Check if the token is valid
	if !token.Valid {
		return nil, models.ErrInvalidToken
	}

	// Type assert to get claims
	claims, ok := token.Claims.(*localjwt.AccessClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	// Additional validation
	if claims.TokenType != TokenTypeAccess {
		return nil, models.ErrInvalidTokenType
	}

	return claims, nil
}

func (j *jwtHandler) VerifyRefreshToken(tokenString string) (*localjwt.RefreshClaims, error) {
	var parseOptions []jwt.ParserOption

	parseOptions = append(parseOptions, jwt.WithValidMethods([]string{"EdDSA"}))
	parseOptions = append(parseOptions, jwt.WithIssuer(j.issuer))
	parseOptions = append(parseOptions, jwt.WithLeeway(leeway))
	parseOptions = append(parseOptions, jwt.WithExpirationRequired())
	parseOptions = append(parseOptions, jwt.WithIssuedAt())

	token, err := jwt.ParseWithClaims(
		tokenString,
		&localjwt.RefreshClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return j.ed25519pub, nil
		},
		parseOptions...,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	// Check if the token is valid
	if !token.Valid {
		return nil, models.ErrInvalidToken
	}

	// Type assert to get claims
	claims, ok := token.Claims.(*localjwt.RefreshClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	// Additional validation
	if claims.TokenType != TokenTypeRefresh {
		return nil, models.ErrInvalidTokenType
	}

	return claims, nil
}
