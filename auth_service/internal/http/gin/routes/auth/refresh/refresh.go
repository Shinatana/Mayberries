package refresh

import (
	ginImpl "auth_service/internal/http/gin"
	requestid "auth_service/internal/http/gin/middlewares/request-id"
	"auth_service/internal/jwt"
	"auth_service/internal/jwt/codec"
	"auth_service/internal/repo"
	"auth_service/pkg/log"
	"github.com/gin-gonic/gin"
	"net/http"
)

type handler struct {
	db  repo.DB
	jwt jwt.Handler
}

func NewRefreshHandler(db repo.DB, jwt jwt.Handler) ginImpl.Router {
	return &handler{db: db, jwt: jwt}
}

func (h *handler) Register(router gin.IRouter) {
	router.POST("/refresh", h.post())
}

func (h *handler) post() gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error

		requestID := c.GetHeader(requestid.Header)

		lg := log.Copy().With(
			"requestID", requestID,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
		)
		lg.Debug("request received")

		var req struct {
			RefreshToken string `json:"refresh_token" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			lg.Error("invalid request body", "error", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}

		tokenString, err := codec.Decode(req.RefreshToken)
		if err != nil {
			lg.Error("invalid refresh token", "error", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
			return
		}

		claims, err := h.jwt.VerifyRefreshToken(tokenString)
		if err != nil {
			lg.Error("invalid refresh token", "error", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
			return
		}

		userID, err := h.db.GetUserIDByEmail(c.Request.Context(), claims.Subject)
		if err != nil {
			lg.Error("failed to get userID by email", "error", err, "email", claims.Subject)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}

		roles, err := h.db.GetUserRoles(c.Request.Context(), userID)
		if err != nil {
			lg.Error("failed to get roles", "error", err, "userID", userID)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}

		permissions, err := h.db.GetUserPermissions(c.Request.Context(), userID)
		if err != nil {
			lg.Error("failed to get permissions", "error", err, "userID", userID)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}

		accessToken, refreshToken, err := h.jwt.GenerateTokenPair(claims.Subject, roles, permissions)
		if err != nil {
			lg.Error("invalid refresh token", "error", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
			return
		}

		accessToken = codec.Encode(accessToken)
		refreshToken = codec.Encode(refreshToken)

		c.JSON(200, gin.H{
			"access_token":       accessToken,
			"token_type":         "bearer",
			"expires_in":         h.jwt.GetTokenLifetime().Seconds(),
			"refresh_token":      refreshToken,
			"refresh_expires_in": h.jwt.GetRefreshTokenLifetime().Seconds(),
		})
		lg.Info("refresh token successful", "subject", claims.Subject)
	}
}
