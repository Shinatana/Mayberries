package init

import (
	"github.com/gin-gonic/gin"
	"order_service/internal/handlers"
	ginImpl "order_service/internal/http/gin"
	"order_service/internal/http/middlewares/recovery"
	requestid "order_service/internal/http/middlewares/request-id"
	"order_service/internal/service"
)

const (
	OrdersPath = "/orders/"
)

func Gin(svc service.Service) *gin.Engine {
	ginRouter := ginImpl.NewGinServer()

	// Middleware
	ginRouter.AddMiddleware(
		recovery.Middleware(),
		requestid.Middleware(),
	)

	// Register route groups
	ginRouter.AddRouters(
		postOrders(OrdersPath, svc),
	)

	return ginRouter.Build()
}

func postOrders(path string, svc service.Service) *ginImpl.Group {
	Orders := ginImpl.NewGroup(path)

	router := handlers.PostOrders(svc)

	Orders.AddRouters(router)

	return Orders
}
