package deleteUser

import (
	ginImpl "auth_service/internal/http/gin"
	requestid "auth_service/internal/http/gin/middlewares/request-id"
	"auth_service/internal/models"
	"auth_service/internal/repo"
	"auth_service/pkg/log"
	"auth_service/pkg/val"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type handler struct {
	db repo.DB
}

func DeleteUser(db repo.DB) ginImpl.Router {
	return &handler{db: db}
}

func (h *handler) Register(router gin.IRouter) {
	router.DELETE("/:id", h.delete())
}

func (h *handler) delete() func(c *gin.Context) {
	return func(c *gin.Context) {
		requestID := c.GetHeader(requestid.Header)

		lg := log.Copy().With(
			"requestID", requestID,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
		)
		lg.Debug("request received")

		userID := c.Param("id")

		if err := val.ValidateWithTag(userID, "required,uuid4"); err != nil {
			lg.Error("invalid product ID header", "product_id", userID)
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid or missing product ID"})
			return
		}

		ID, err := uuid.Parse(userID)

		err = h.db.DeleteUser(c.Request.Context(), ID)
		if err != nil {
			if errors.Is(err, models.ErrUserNotFound) {
				lg.Error("user not found", "user_id", userID)
				c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
				return
			}
			lg.Error("failed to delete user", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete user"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "user deleted successfully", "id": userID})
		lg.Info("user deleted successfully", "id", userID)
	}
}
