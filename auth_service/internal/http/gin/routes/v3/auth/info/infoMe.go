package info

import (
	ginImpl "auth_service/internal/http/gin"
	middleAuth "auth_service/internal/http/gin/middlewares/getAuth"
	requestid "auth_service/internal/http/gin/middlewares/request-id"
	"auth_service/internal/jwt"
	"auth_service/internal/models"
	"auth_service/internal/repo"
	"auth_service/pkg/log"
	"auth_service/pkg/val"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type handler struct {
	db  repo.DB
	jwt jwt.Handler
}

func TakeInfoMe(db repo.DB, jwt jwt.Handler) ginImpl.Router {
	return &handler{db: db, jwt: jwt}
}

func (h *handler) Register(router gin.IRouter) {
	router.GET("/auth/me", middleAuth.GetAuthMiddleware(h.jwt), h.get())
}

func (h *handler) get() func(c *gin.Context) {
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

		emailVal, exists := c.Get("email")
		if !exists {
			lg.Error("email not found in context")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		email, ok := emailVal.(string)
		if !ok {
			lg.Error("email in context is not a string")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}

		userID, err := h.db.GetUserIDByEmail(c.Request.Context(), email)
		if err != nil {
			if errors.Is(err, models.ErrUserNotFound) {
				lg.Warn("user not found by email", "email", email)
				c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
				return
			}
			lg.Error("failed to get userID by email", "error", err, "email", email)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}

		var user *models.Users

		user, err = h.db.GetUserByID(c.Request.Context(), userID)
		if err != nil {
			if errors.Is(err, models.ErrUserNotFound) {
				lg.Warn("user not found by ID", "userID", userID)
				c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
				return
			}
			lg.Error("failed to get user by ID", "error", err, "userID", userID)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}

		role, err := h.db.GetUserRoles(c.Request.Context(), userID)
		if err != nil {
			if errors.Is(err, models.ErrUserNotFound) {
				lg.Warn("user not found by userID", "userID", userID)
				c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
				return
			}
			lg.Error("failed to get user roles", "error", err, "userID", userID)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"createdAt": user.CreatedAt,
			"roleId":    user.RoleID,
			"roles":     role,
			"email":     user.Email,
			"name":      user.Name,
			"id":        userID,
		})

		lg.Info("user info returned successfully", "userID", userID)

	}
}
