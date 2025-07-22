package deleteRole

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

func DeleteRole(db repo.DB) ginImpl.Router {
	return &handler{db: db}
}

func (h *handler) Register(router gin.IRouter) {
	router.DELETE("/:id", h.delete())
}

func (h *handler) delete() func(c *gin.Context) {
	return func(c *gin.Context) {
		var err error

		requestID := c.GetHeader(requestid.Header)

		lg := log.Copy().With(
			"requestID", requestID,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
		)

		lg.Debug("request received")

		roleID := c.Param("id")

		id, err := strconv.Atoi(roleID)
		if err != nil {
			lg.Error("invalid roleId", "roleId", roleID)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid roleId",
			})
			return
		}

		err = h.db.DeleteRole(c.Request.Context(), id)

		if err != nil {

			if errors.Is(err, models.ErrRoleNotFound) {
				lg.Error("role not found", "roleID", id)
				c.JSON(http.StatusNotFound, gin.H{"error": "role not found"})
				return
			}
			lg.Error("failed to deleteRole role", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to deleteRole role"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "role deleted successfully", "id": id})
		lg.Info("role deleted successfully", "id", id)
	}
}
