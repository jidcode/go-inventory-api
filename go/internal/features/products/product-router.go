package products

import (
	"github.com/labstack/echo/v4"
	"github.com/ventry/internal/features/auth"
)

func ProductRouter(e *echo.Echo, productHandler *ProductHandler, imagesHandler *ImagesHandler, authService auth.AuthService) {
	api := e.Group("/api/products")
	api.Use(auth.AuthMiddleware(&authService), auth.RoleMiddleware("user"))

	api.GET("/inventory/:inventoryId", productHandler.ListProducts)
	api.GET("/:id", productHandler.GetProduct)
	api.POST("", productHandler.CreateProduct)
	api.PUT("/:id", productHandler.UpdateProduct)
	api.DELETE("/:id", productHandler.DeleteProduct)

	img := e.Group("api/images")
	img.GET("/product/:productId", imagesHandler.ListImages)
	img.GET("/:id", imagesHandler.GetImage)
	img.POST("", imagesHandler.CreateImage)
	img.PUT("/:id", imagesHandler.UpdateImage)
	img.DELETE("/:id", imagesHandler.DeleteImage)
}
