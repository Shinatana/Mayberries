package register

import (
	"auth_service/internal/hash"
	"errors"
	"github.com/google/uuid"
	"net/http"

	"github.com/gin-gonic/gin"

	ginImpl "auth_service/internal/http/gin"
	requestid "auth_service/internal/http/gin/middlewares/request-id"
	"auth_service/internal/models"
	"auth_service/internal/repo"
	"auth_service/pkg/log"
	"auth_service/pkg/val"
)

type handler struct {
	db     repo.DB
	hasher hash.Hasher
}

func NewRegisterHandler(db repo.DB, hasher hash.Hasher) ginImpl.Router {
	return &handler{db: db, hasher: hasher}
}

func (h *handler) Register(router gin.IRouter) {
	router.POST("/auth/register", h.post())
}

func (h *handler) post() func(c *gin.Context) {
	return func(c *gin.Context) {
		var err error

		requestID := c.GetHeader(requestid.Header)

		lg := log.Copy().With(
			"requestID", requestID,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
		)
		lg.Debug("request received")

		var user models.Users
		if err = c.ShouldBindJSON(&user); err != nil {
			lg.Error("failed to bind request body", "error", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid request body",
			})
			return
		}
		user.ID = uuid.New()

		if err = val.ValidateStruct(user); err != nil {
			lg.Error("failed to validate request body", "error", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid body content",
			})
			return
		}

		hashedPassword, err := h.hasher.Hash(user.Password)
		if err != nil {
			lg.Error("failed to hash password", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "internal server error",
			})
			return
		}
		user.PasswordHash = hashedPassword

		lg.Debug("password hashed successfully", "email", user.Email)

		var role models.Role
		err = h.db.GetRoleByName(c.Request.Context(), "user", &role)
		if err != nil {
			lg.Error("default role 'user' not found", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to assign default role",
			})
			return
		}
		user.RoleID = role.ID

		if err = h.db.RegisterUser(c.Request.Context(), user); err != nil {
			if errors.Is(err, models.ErrDuplicateUser) {
				lg.Warn("user already exists", "email", user.Email)
				c.Status(http.StatusConflict)
				return
			}
			lg.Error("failed to register user", "error", err)
			c.Status(http.StatusInternalServerError)
			return
		}

		lg.Info("request completed")
		c.Status(http.StatusCreated)
	}
}
