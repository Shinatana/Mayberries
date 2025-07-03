package login

import (
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
	db repo.DB
}

func NewLoginHandler(db repo.DB) ginImpl.Router {
	return &handler{db: db}
}

func (h *handler) Register(router gin.IRouter) {
	router.GET("/auth/login", h.post())
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

		email, _, err := basicAuth(c, lg)
		if err != nil {
			lg.Error("basic auth error", "error", err)
			c.Status(http.StatusUnauthorized)
			return
		}

		_, err = h.db.GetUserPassword(c.Request.Context(), email)
		if err != nil {
			if errors.Is(err, models.ErrUserNotFound) {
				lg.Error("user not found")
				c.Status(http.StatusUnauthorized)
				return
			}

			lg.Error("db query error", "error", err, "email", email)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"token_type": "bearer",
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
