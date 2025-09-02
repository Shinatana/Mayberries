package delete_orders

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mayberries/shared/pkg/log"

	"net/http"
	ginImpl "order_service/internal/http/gin"
	requestid "order_service/internal/http/middlewares/request-id"
	"order_service/internal/models"
	"order_service/internal/service"
)

type handler struct {
	svc service.Service
}

func DeleteOrders(svc service.Service) ginImpl.Router {
	return &handler{svc: svc}
}

func (h *handler) Register(router gin.IRouter) {
	router.DELETE("/:id", h.delete())
}

func (h *handler) delete() func(c *gin.Context) {
	return func(c *gin.Context) {
		var err error

		requestID := c.GetHeader(requestid.Header)

		c.Header(requestid.Header, requestID)

		lg := log.Copy().With(
			"requestID", requestID,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
		)
		lg.Debug("request received")

		idStr := c.Param("id")
		uid, err := uuid.Parse(idStr)
		if err != nil {
			lg.Warn("invalid id", "id", idStr, "error", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}

		err = h.svc.DeleteOrder(c.Request.Context(), uid)
		if err != nil {
			if errors.Is(err, models.ErrOrderNotFound) {
				lg.Warn("order not found", "id", uid)
				c.JSON(http.StatusConflict, gin.H{"error": "order not found"})
				return
			}
			if errors.Is(err, models.ErrIDIsNil) {
				lg.Warn("empty id", "error", err)
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			lg.Error("failed to get order", "error", err)
			c.Status(http.StatusInternalServerError)
			return
		}
		lg.Info("order delete", uid)
		c.JSON(http.StatusOK, gin.H{"message": "order deleted"})
	}
}
