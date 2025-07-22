package recovery

import (
	ginImpl "catalog_service/internal/http/gin"
	"github.com/gin-gonic/gin"
)

func Middleware() ginImpl.Middleware {
	return ginImpl.MiddlewareFunc(gin.Recovery())
}
