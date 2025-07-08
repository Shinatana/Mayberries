package login

import (
	"auth_service/internal/hash"
	"auth_service/internal/jwt"
	"auth_service/internal/repo"
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	ginImpl "auth_service/internal/http/gin"
	requestid "auth_service/internal/http/gin/middlewares/request-id"
	"auth_service/internal/models"
	"auth_service/pkg/log"
	"auth_service/pkg/val"
)

type handler struct {
	db     repo.DB
	jwt    jwt.Handler
	hasher hash.Hasher
}

func NewLoginHandler(db repo.DB, jwt jwt.Handler) ginImpl.Router {
	return &handler{db: db, jwt: jwt}
}
func (h *handler) Register(router gin.IRouter) {
	router.POST("/auth/login", h.post())
}

func (h *handler) post() func(c *gin.Context) {
	return func(c *gin.Context) {
		var err error

		requestID := c.GetHeader(requestid.Header)
		err = val.ValidateWithTag(requestID, "required,uuid4")
		if err != nil {
			log.Warn("invalid request id provided", requestid.Header, requestID)
		}
		c.Header(requestid.Header, requestID)

		lg := log.Copy().With(
			"requestID", requestID,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
		)
		lg.Debug("request received")

		email, pwd, err := basicAuth(c, lg)
		if err != nil {
			lg.Error("basic auth error", "error", err)
			c.Status(http.StatusUnauthorized)
			return
		}

		dbPasswordHash, err := h.db.GetUserPassword(c.Request.Context(), email)
		if err != nil {
			if errors.Is(err, models.ErrUserNotFound) {
				lg.Error("user not found")
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": "user not found",
				})
				return
			}

			lg.Error("db query error", "error", err, "email", email)

			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "internal server error",
			})
			return
		}

		if err = h.hasher.CheckHash(pwd, dbPasswordHash); err != nil {
			lg.Error("invalid password", "email", email)
			c.Status(http.StatusUnauthorized)
			return
		}
		lg.Debug("password validated successfully", "email", email)

		accessToken, refreshToken, err := h.jwt.GenerateTokenPair(email)
		if err != nil {
			lg.Error("failed to generate token pair", "error", err)
			c.Status(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"access_token":       accessToken,
			"token_type":         "bearer",
			"expires_in":         h.jwt.GetTokenLifetime(),
			"refresh_token":      refreshToken,
			"refresh_expires_in": h.jwt.GetRefreshTokenLifetime(),
		})

		lg.Info("login successful", "email", email)

	}

}

func basicAuth(c *gin.Context, lg *slog.Logger) (email string, pwd string, err error) {
	email, pwd, hasAuth := c.Request.BasicAuth()
	if !hasAuth {
		const errMsg = "authentication required"
		lg.Error(errMsg)
		return "", "", errors.New(errMsg)
	}

	if err = val.ValidateWithTag(email, "required,email"); err != nil {
		const errMsg = "invalid username"
		lg.Error(errMsg)
		return "", "", errors.New(errMsg)
	}

	if err = val.ValidateWithTag(pwd, "required,min=12"); err != nil {
		const errMsg = "bad password"
		lg.Error(errMsg)
		return "", "", errors.New(errMsg)
	}

	return email, pwd, nil
}
