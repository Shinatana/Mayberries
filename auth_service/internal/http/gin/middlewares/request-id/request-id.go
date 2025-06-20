package requestid

import (
	"auth_service/internal/http/gin"
	"github.com/gin-contrib/requestid"
)

const Header = "X-Request-Id"

func Middleware() gin.Middleware {
	return gin.MiddlewareFunc(requestid.New())
}
