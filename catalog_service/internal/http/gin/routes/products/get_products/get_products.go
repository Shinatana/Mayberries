package get_products

import (
	ginImpl "catalog_service/internal/http/gin"
	requestid "catalog_service/internal/http/gin/middlewares/request-id"
	"catalog_service/internal/repo"
	"catalog_service/pkg/log"
	"catalog_service/pkg/val"
	"github.com/gin-gonic/gin"
	"net/http"
)

type handler struct {
	db repo.DB
}

func GetProducts(db repo.DB) ginImpl.Router {
	return &handler{db: db}
}

func (h *handler) Register(router gin.IRouter) {
	router.GET("/", h.get())
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

		products, err := h.db.GetProducts(c.Request.Context())
		if err != nil {
			lg.Error("Failed to fetch products", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch products"})
			return
		}

		c.JSON(http.StatusOK, products)
		lg.Info("Products fetched successfully", "count", len(products))
	}
}
