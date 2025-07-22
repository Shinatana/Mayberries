package infoRole

import (
	ginImpl "auth_service/internal/http/gin"
	requestid "auth_service/internal/http/gin/middlewares/request-id"
	"auth_service/internal/models"
	"auth_service/internal/repo"
	"auth_service/pkg/log"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type handler struct {
	db repo.DB
}

func TakeInfoRole(db repo.DB) ginImpl.Router {
	return &handler{db: db}
}

func (h *handler) Register(router gin.IRouter) {
	router.GET("/:id", h.get())
}

func (h *handler) get() func(c *gin.Context) {
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

		roleId := c.Params.ByName("id")

		id, err := strconv.Atoi(roleId)
		if err != nil {
			lg.Error("invalid roleId", "roleId", roleId)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid roleId",
			})
			return
		}

		role, err = h.db.GetRoleByID(c.Request.Context(), id)
		if err != nil {
			if errors.Is(err, models.ErrRoleNotFound) {
				lg.Error("role not found", "roleid", id)
				c.JSON(http.StatusNotFound, gin.H{
					"error": "role not found",
				})
				return
			}
			lg.Error("failed to fetch role", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch role"})
			return
		}
		c.JSON(http.StatusOK, role)
		lg.Info("product fetched successfully", "id", roleId)
	}
}
