package init

import (
	ginImpl "catalog_service/internal/http/gin"
	"catalog_service/internal/http/gin/middlewares/recovery"
	requestid "catalog_service/internal/http/gin/middlewares/request-id"
	deleteproducts "catalog_service/internal/http/gin/routes/products/delete_products"
	getproducts "catalog_service/internal/http/gin/routes/products/get_products"
	getproductsbyId "catalog_service/internal/http/gin/routes/products/get_products_by_id"
	postproducts "catalog_service/internal/http/gin/routes/products/post_produsts"
	"catalog_service/internal/repo"
	"github.com/gin-gonic/gin"
)

const v1ProductsPath = "/products"

func Gin(db repo.DB) *gin.Engine {
	ginRouter := ginImpl.NewGinServer()

	// Middleware
	ginRouter.AddMiddleware(
		recovery.Middleware(),
		requestid.Middleware(),
	)

	// Register route groups
	ginRouter.AddRouters(
		getV1Products(v1ProductsPath, db),
	)

	return ginRouter.Build()
}

func getV1Products(path string, db repo.DB) *ginImpl.Group {
	v1Products := ginImpl.NewGroup(path)
	v1Products.AddRouters(
		getproducts.GetProducts(db),
		postproducts.PostProducts(db),
		getproductsbyId.GetProducts(db),
		deleteproducts.DeleteProducts(db),
	)
	return v1Products
}
