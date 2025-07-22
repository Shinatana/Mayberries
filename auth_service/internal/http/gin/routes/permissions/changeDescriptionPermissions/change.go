package changeDescriptionPermissions

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

func NewDescriptionPermissions(db repo.DB) ginImpl.Router {
	return &handler{db: db}
}

func (h *handler) Register(router gin.IRouter) {
	router.POST("/id", h.post())
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

		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			lg.Error("invalid id param", "id", idParam)
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id parameter"})
			return
		}
		var permission models.Permission
		if err := c.ShouldBindJSON(&permission); err != nil {
			lg.Error("invalid request body", "error", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}

		if permission.ID != id {
			lg.Error("id in URL and body do not match", "url_id", id, "body_id", permission.ID)
			c.JSON(http.StatusBadRequest, gin.H{"error": "id in URL and body do not match"})
			return
		}

		err = h.db.ChangeDescriptionPermissions(c.Request.Context(), &permission)
		if err != nil {
			if errors.Is(err, models.ErrPermissionNotFound) {
				lg.Error("permission not found", "error", err)
				c.JSON(http.StatusNotFound, gin.H{"error": "permission not found"})
				return
			}
			lg.Error("error while changing description permissions", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}

		lg.Info("permission changed")
		c.JSON(http.StatusOK, gin.H{
			"permission": permission,
		})
	}
}
