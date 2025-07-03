package register

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	ginImpl "auth_service/internal/http/gin"
	requestid "auth_service/internal/http/gin/middlewares/request-id"
	"auth_service/internal/models"
	"auth_service/internal/repo"
	"auth_service/pkg/log"
	"auth_service/pkg/val"
)

type handler struct {
	db repo.DB
}

func NewRegisterHandler(db repo.DB) ginImpl.Router {
	return &handler{db: db}
}

func (h *handler) Register(router gin.IRouter) {
	router.POST("/auth/register", h.post())
}

func (h *handler) post() func(c *gin.Context) {
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

		var user models.Users
		if err = c.ShouldBindJSON(&user); err != nil {
			lg.Error("failed to bind request body", "error", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid request body",
			})
			return
		}

		sUser := user
		sUser.Password = "********"
		lg.Debug("request body", "body", sUser)

		if err = val.ValidateStruct(user); err != nil {
			lg.Error("failed to validate request body", "error", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid body content",
			})
			return
		}

		if err = h.db.RegisterUser(c.Request.Context(), user); err != nil {
			if errors.Is(err, models.ErrDuplicateUser) {
				lg.Warn("user already exists", "email", user.Email)
				c.Status(http.StatusConflict)
				return
			}
			lg.Error("failed to register user", "error", err)
			c.Status(http.StatusInternalServerError)
			return
		}

		lg.Info("request completed")
		c.Status(http.StatusCreated)
	}
}
