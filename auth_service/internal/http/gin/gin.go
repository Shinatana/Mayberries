package gin

import (
	"github.com/gin-gonic/gin"
)

// RouterFunc is a function that implements the Router interface
type RouterFunc func(router gin.IRouter)

// Register implements the Router interface
func (f RouterFunc) Register(router gin.IRouter) {
	f(router)
}

// MiddlewareFunc is a function that implements the Middleware interface
type MiddlewareFunc gin.HandlerFunc

// Handler implements the Middleware interface
func (f MiddlewareFunc) Handler() gin.HandlerFunc {
	return gin.HandlerFunc(f)
}

// Server represents a gin HTTP server
type Server struct {
	engine      *gin.Engine
	middlewares []Middleware
	routers     []Router
}

// NewGinServer creates a new Server
func NewGinServer() *Server {
	return &Server{
		engine:      gin.New(),
		middlewares: make([]Middleware, 0),
		routers:     make([]Router, 0),
	}
}

// AddMiddleware adds global middlewares to the server
func (s *Server) AddMiddleware(middlewares ...Middleware) *Server {
	s.middlewares = append(s.middlewares, middlewares...)
	return s
}

// AddRouters adds routers to the server
func (s *Server) AddRouters(routers ...Router) *Server {
	s.routers = append(s.routers, routers...)
	return s
}

// Build configures and returns the underlying gin.Engine
func (s *Server) Build() *gin.Engine {
	// Apply global middlewares
	for _, middleware := range s.middlewares {
		s.engine.Use(middleware.Handler())
	}

	// Register routers
	for _, router := range s.routers {
		router.Register(s.engine)
	}

	return s.engine
}
