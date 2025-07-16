package get_categories

import (
	ginImpl "catalog_service/internal/http/gin"
	requestid "catalog_service/internal/http/gin/middlewares/request-id"
	"catalog_service/internal/service/categories"
	"catalog_service/internal/service/products"
	"catalog_service/pkg/log"
	"github.com/gin-gonic/gin"
	"net/http"
)

type handler struct {
	categories categories.Service
}

func GetCategories(service products.Service) ginImpl.Router {
	return &handler{categories: service}
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

		categoriesAll, err := h.categories.GetCategories(c.Request.Context())
		if err != nil {
			lg.Error("Failed to fetch categories", "error", err)
			ErrFetchCategories
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch categories"})
			return
		}

		c.JSON(http.StatusOK, categoriesAll)
		lg.Info("Categories fetched successfully")
	}
}
