package recovery

import (
	"github.com/gin-gonic/gin"
	ginImpl "order_service/internal/http/gin"
)

func Middleware() ginImpl.Middleware {
	return ginImpl.MiddlewareFunc(gin.Recovery())
}
