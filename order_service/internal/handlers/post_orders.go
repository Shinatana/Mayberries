package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
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

func PostOrders(svc service.Service) ginImpl.Router {
	return &handler{svc: svc}
}

func (h *handler) Register(router gin.IRouter) {
	router.POST("/", h.post())
}

func (h *handler) post() func(c *gin.Context) {
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

		var order models.Order

		if err := c.ShouldBindJSON(&order); err != nil {
			lg.Error("failed to bind request body", "error", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}
		id, err := h.svc.CreateOrder(c.Request.Context(), order)
		if err != nil {
			if errors.Is(err, models.ErrDuplicateOrder) {
				lg.Warn("order already exists", "id", order.ID)
				c.JSON(http.StatusConflict, gin.H{"error": "order already exists"})
				return
			}
			if errors.Is(err, models.ErrValidation) {
				lg.Warn("validation failed", "error", err)
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			lg.Error("failed to create order", "error", err)
			c.Status(http.StatusInternalServerError)
			return
		}
		lg.Info("order created", "order_id", id)
		c.JSON(http.StatusCreated, gin.H{"id": id})

	}
}
