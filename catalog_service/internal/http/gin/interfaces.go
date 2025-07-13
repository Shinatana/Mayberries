package gin

import (
	"github.com/gin-gonic/gin"
)

// Router is the main interface for registering routes
type Router interface {
	Register(router gin.IRouter)
}

// Middleware represents gin middleware
type Middleware interface {
	Handler() gin.HandlerFunc
}
