package init

import (
	"github.com/gin-gonic/gin"
	delete "order_service/internal/handlers/delete_orders"
	get "order_service/internal/handlers/get_orders"
	patch "order_service/internal/handlers/patch_orders"
	post "order_service/internal/handlers/post_orders"
	ginImpl "order_service/internal/http/gin"
	"order_service/internal/http/middlewares/recovery"
	requestid "order_service/internal/http/middlewares/request-id"
	"order_service/internal/service"
)

const (
	OrdersPath = "/orders"
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
		CreateGroupOrders(OrdersPath, svc),
	)

	return ginRouter.Build()
}

func CreateGroupOrders(path string, svc service.Service) *ginImpl.Group {
	orders := ginImpl.NewGroup(path)

	orders.AddRouters(
		post.PostOrders(svc),
		get.GetOrders(svc),
		delete.DeleteOrders(svc),
		patch.PatchOrders(svc),
	)

	return orders
}
