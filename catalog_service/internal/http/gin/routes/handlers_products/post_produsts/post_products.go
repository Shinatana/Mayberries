package post_products

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

const (
	defaultCategory = 1
)

type handler struct {
	db repo.DB
}

func PostProducts(db repo.DB) ginImpl.Router {
	return &handler{db: db}
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

		var product models.Products

		if err = c.ShouldBindJSON(&product); err != nil {
			lg.Error("failed to bind request body", "error", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid request body",
			})
			return
		}
		product.ID = uuid.New()

		if err = val.ValidateStruct(product); err != nil {
			lg.Error("failed to validate request body", "error", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid request body",
			})
			return
		}
		//Until the mechanics of creating and removing roles are implemented, a constant is used
		product.CategoryID = defaultCategory

		if err = h.db.PostProducts(c.Request.Context(), product); err != nil {
			if errors.Is(err, models.ErrDuplicateProducts) {
				lg.Warn("product already exists", defaultCategory, product.ID)
				c.JSON(http.StatusConflict, gin.H{
					"error": "product already exists",
				})
				return
			}
			lg.Error("failed to post handlers_products", "error", err)
			c.Status(http.StatusInternalServerError)
			return
		}
		lg.Info("request completed")

		c.Status(http.StatusCreated)
	}
}
