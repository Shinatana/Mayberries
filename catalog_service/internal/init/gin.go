package init

import (
	ginImpl "catalog_service/internal/http/gin"
	"catalog_service/internal/http/gin/middlewares/recovery"
	requestid "catalog_service/internal/http/gin/middlewares/request-id"
	getcategories "catalog_service/internal/http/gin/routes/handlers_categories/get_categories"
	deleteproducts "catalog_service/internal/http/gin/routes/handlers_products/delete_products"
	getproducts "catalog_service/internal/http/gin/routes/handlers_products/get_products"
	getproductsbyId "catalog_service/internal/http/gin/routes/handlers_products/get_products_by_id"
	postproducts "catalog_service/internal/http/gin/routes/handlers_products/post_produsts"
	"catalog_service/internal/repo"
	"catalog_service/internal/service/products"
	"github.com/gin-gonic/gin"
)

const (
	ProductsPath   = "/products"
	CategoriesPath = "/categories"
)

func Gin(db repo.DB) *gin.Engine {
	ginRouter := ginImpl.NewGinServer()

	// Middleware
	ginRouter.AddMiddleware(
		recovery.Middleware(),
		requestid.Middleware(),
	)

	// Register route groups
	ginRouter.AddRouters(
		getProducts(ProductsPath, db),
		getCategories(CategoriesPath, db),
	)

	return ginRouter.Build()
}

func getProducts(path string, db repo.DB) *ginImpl.Group {
	Products := ginImpl.NewGroup(path)

	Products.AddRouters(
		getproducts.GetProducts(products.Service{DB: db}),
		postproducts.PostProducts(db),
		getproductsbyId.GetProducts(db),
		deleteproducts.DeleteProducts(db),
	)
	return Products
}

func getCategories(path string, db repo.DB) *ginImpl.Group {
	Categories := ginImpl.NewGroup(path)

	Categories.AddRouters(
		getcategories.GetCategories(products.Service{DB: db}),
	)

	return Categories
}
