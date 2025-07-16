package createRole

import (
	ginImpl "auth_service/internal/http/gin"
	requestid "auth_service/internal/http/gin/middlewares/request-id"
	"auth_service/internal/models"
	"auth_service/internal/repo"
	"auth_service/pkg/log"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type handler struct {
	db repo.DB
}

func CreateRole(db repo.DB) ginImpl.Router {
	return &handler{db: db}
}

func (h *handler) Register(router gin.IRouter) {
	router.POST("/create", h.post())
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

		var role *models.Role

		if err = c.ShouldBindJSON(&role); err != nil {
			lg.Error("failed to bind request body", "error", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid request body",
			})
			return
		}

		err = h.db.CreateRole(c.Request.Context(), role)
		if err != nil {
			if errors.Is(err, models.ErrDuplicateRole) {
				lg.Error("duplicate role", "role", role.Name)
				c.JSON(http.StatusConflict, gin.H{"error": "role already exists"})
				return
			}
			lg.Error("failed to create role", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create role"})
			return
		}

		c.JSON(http.StatusOK, role)
		lg.Info("role created successfully", "id", role.ID)
	}
}
