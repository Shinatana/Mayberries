package init

import (
	"github.com/gin-gonic/gin"
	"order_service/internal/handlers/handlers_orders/post_orders"
	ginImpl "order_service/internal/http/gin"
	"order_service/internal/http/middlewares/recovery"
	"order_service/internal/http/middlewares/request-id"
	"order_service/internal/service/order"
)

const (
	OrdersPath = "/orders/"
)

func Gin(svc order.Service) *gin.Engine {
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

func postOrders(path string, svc order.Service) *ginImpl.Group {
	Orders := ginImpl.NewGroup(path)

	router := post_orders.PostOrders(svc)

	Orders.AddRouters(router)

	return Orders
}
