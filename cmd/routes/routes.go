package routes

import (
	"net/http"
	"project/internal/item_detail/repo/datasource/dal"
	"project/internal/item_detail/repo/datasource/dao"
	"project/internal/item_detail/rest"
	"project/internal/item_detail/service"

	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
)

var (
	productDal  = dal.NewProductDAL()
	sellerDal   = dal.NewSellerDAL()
	categoryDal = dal.NewCategoryDAL()
	imageDal    = dal.NewImageDAL()
)

func crud[T any](group *echo.Group, dal dao.CrudDAO[T]) {
	crudService := service.NewCrudService(dal)
	crudHandler := rest.NewCrudHandler(crudService)

	group.POST("", crudHandler.CreateEntity)
	group.GET("", crudHandler.GetAllEntities)
	group.GET("/:id", crudHandler.GetEntityByID)
	group.PATCH("/:id", crudHandler.UpdateEntity)
	group.DELETE("/:id", crudHandler.DeleteEntity)
}

func productRouter(r *echo.Group) {
	productGroup := r.Group("/products")

	crud(productGroup, productDal)

	productService := service.NewProductService(
		productDal,
		sellerDal,
		categoryDal,
		imageDal,
	)

	productHandler := rest.NewProductHandler(productService)

	productGroup.GET("/:id/category", productHandler.GetCategories)
	productGroup.GET("/:id/seller", productHandler.GetSellers)
	productGroup.GET("/:id/characteristic", productHandler.GetCharacteristic)
	productGroup.GET("/:id/details", productHandler.GetDetails)
	productGroup.GET("/:id/images", productHandler.GetImages)
	productGroup.POST("/:id/images", productHandler.AddImages)

	productGroup.PATCH("/:id/category", productHandler.ChangeCategories)
	productGroup.PATCH("/:id/seller", productHandler.ChangeSellers)
}

func imageRouter(r *echo.Group) {
	imageGroup := r.Group("/images")

	crud(imageGroup, imageDal)
}

func categoryRouter(r *echo.Group) {
	categoryGroup := r.Group("/categories")

	crud(categoryGroup, categoryDal)
}

func sellerRouter(r *echo.Group) {
	sellerGroup := r.Group("/sellers")

	crud(sellerGroup, sellerDal)
}

func Routes(r *echo.Echo) *echo.Echo {
	r.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "UP")
	})

	api := r.Group("/api/v1")

	// Metricas de prometheus
	r.Use(echoprometheus.NewMiddleware("item_detail"))
	r.GET("/metrics", echoprometheus.NewHandler())

	api.Use(ApiKeyMiddleware)

	// Routes
	productRouter(api)
	categoryRouter(api)
	sellerRouter(api)
	imageRouter(api)

	return r
}
