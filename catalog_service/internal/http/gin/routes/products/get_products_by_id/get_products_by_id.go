package get_products_by_id

import (
	ginImpl "catalog_service/internal/http/gin"
	requestid "catalog_service/internal/http/gin/middlewares/request-id"
	"catalog_service/internal/models"
	"catalog_service/internal/repo"
	"catalog_service/pkg/log"
	"catalog_service/pkg/val"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type handler struct {
	db repo.DB
}

func GetProducts(db repo.DB) ginImpl.Router {
	return &handler{db: db}
}

func (h *handler) Register(router gin.IRouter) {
	router.GET("/:id", h.get())
}

func (h *handler) get() func(c *gin.Context) {
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

		var products *models.Products

		productID := c.Param("id")
		if err := val.ValidateWithTag(productID, "required,uuid4"); err != nil {
			lg.Error("invalid product ID header", "product_id", productID)
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid or missing product ID"})
			return
		}

		productUUID, err := uuid.Parse(productID)

		products, err = h.db.GetProductsById(c.Request.Context(), productUUID)
		if err != nil {
			lg.Error("failed to fetch product", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch product"})
			return
		}

		// Ответ
		c.JSON(http.StatusOK, products)
		lg.Info("product fetched successfully", "id", productID)
	}
}
