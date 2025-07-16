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

// Group represents a group of routes with optional middlewares
type Group struct {
	path        string
	middlewares []Middleware
	routers     []Router
}

// NewGroup creates a new route group
func NewGroup(path string) *Group {
	return &Group{
		path:        path,
		middlewares: make([]Middleware, 0),
		routers:     make([]Router, 0),
	}
}

// AddMiddleware adds middlewares to the group
func (g *Group) AddMiddleware(middlewares ...Middleware) *Group {
	g.middlewares = append(g.middlewares, middlewares...)
	return g
}

// AddRouters adds routers to the group
func (g *Group) AddRouters(routers ...Router) *Group {
	g.routers = append(g.routers, routers...)
	return g
}

// Register implements the Router interface
func (g *Group) Register(router gin.IRouter) {
	group := router.Group(g.path)

	for _, m := range g.middlewares {
		group.Use(m.Handler())
	}

	for _, r := range g.routers {
		r.Register(group)
	}
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
	for _, m := range s.middlewares {
		s.engine.Use(m.Handler())
	}

	for _, r := range s.routers {
		r.Register(s.engine)
	}

	return s.engine
}
