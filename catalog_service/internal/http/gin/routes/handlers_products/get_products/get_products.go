package get_products

import (
	ginImpl "catalog_service/internal/http/gin"
	requestid "catalog_service/internal/http/gin/middlewares/request-id"
	"catalog_service/internal/service/products"
	"catalog_service/pkg/log"
	"github.com/gin-gonic/gin"
	"net/http"
)

type handler struct {
	productService products.Service
}

func GetProducts(service products.Service) ginImpl.Router {
	return &handler{productService: service}
}

func (h *handler) Register(router gin.IRouter) {
	router.GET("/", h.get())
}

func (h *handler) get() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader(requestid.Header)

		lg := log.Copy().With(
			"requestID", requestID,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
		)
		lg.Debug("request received")

		productsAll, err := h.productService.GetProducts(c.Request.Context())
		if err != nil {
			lg.Error("Failed to fetch handlers_products", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch handlers_products"})
			return
		}

		c.JSON(http.StatusOK, productsAll)
		lg.Info("Products fetched successfully", "count", len(productsAll))
	}
}
