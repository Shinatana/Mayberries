package getAllRoles

import (
	ginImpl "auth_service/internal/http/gin"
	requestid "auth_service/internal/http/gin/middlewares/request-id"
	"auth_service/internal/repo"
	"auth_service/pkg/log"
	"github.com/gin-gonic/gin"
	"net/http"
)

type handler struct {
	db repo.DB
}

func GetAllRoles(db repo.DB) ginImpl.Router {
	return &handler{db: db}
}

func (h *handler) Register(router gin.IRouter) {
	router.GET("/all", h.get())
}

func (h *handler) get() func(c *gin.Context) {
	return func(c *gin.Context) {
		requestID := c.GetHeader(requestid.Header)
		lg := log.Copy().With(
			"requestID", requestID,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
		)
		lg.Debug("request received")

		roles, err := h.db.GetAllRoles(c.Request.Context())
		if err != nil {
			lg.Error("failed to fetch roles", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch roles"})
			return
		}

		c.JSON(http.StatusOK, roles)
		lg.Info("roles fetched successfully", "count", len(roles))
	}
}
