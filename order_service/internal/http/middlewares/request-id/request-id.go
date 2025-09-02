package requestid

import (
	"github.com/gin-contrib/requestid"
	"order_service/internal/http/gin"
)

const Header = "X-Request-Id"

func Middleware() gin.Middleware {
	return gin.MiddlewareFunc(requestid.New())
}
