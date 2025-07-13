package delete_products

import (
	ginImpl "catalog_service/internal/http/gin"
	requestid "catalog_service/internal/http/gin/middlewares/request-id"
	"catalog_service/internal/models"
	"catalog_service/internal/repo"
	"catalog_service/pkg/log"
	"catalog_service/pkg/val"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type handler struct {
	db repo.DB
}

func DeleteProducts(db repo.DB) ginImpl.Router {
	return &handler{db: db}
}

func (h *handler) Register(router gin.IRouter) {
	router.DELETE("/:id", h.delete())
}

func (h *handler) delete() func(c *gin.Context) {
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

		productIDStr := c.Param("id")
		if err := val.ValidateWithTag(productIDStr, "required,uuid4"); err != nil {
			lg.Error("invalid product id param", "product_id", productIDStr)
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
			return
		}

		productID, err := uuid.Parse(productIDStr)
		if err != nil {
			lg.Error("failed to parse product id", "error", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id format"})
			return
		}

		err = h.db.DeleteProducts(c.Request.Context(), productID)
		if err != nil {
			if errors.Is(err, models.ErrProductsNotFound) {
				lg.Warn("product not found", "id", productID)
				c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
				return
			}
			lg.Error("failed to delete product", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete product"})
			return
		}

		lg.Info("product deleted successfully", "id", productID)
		c.Status(http.StatusNoContent)
	}
}
