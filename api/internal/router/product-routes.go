package router

import (
	"github.com/labstack/echo/v4"
	"github.com/ventry/internal/features/products"
	"github.com/ventry/internal/pkg/auth"
)

func ProductRoutes(e *echo.Echo, pc products.ProductController, authService auth.AuthService) {
	api := e.Group("/api/products")
	api.Use(auth.AuthMiddleware(&authService), auth.RoleMiddleware("user"))

	api.GET("/inventory/:inventoryId", pc.ListProducts)
	api.GET("/:id", pc.GetProduct)
	api.POST("", pc.CreateProduct)
	api.PUT("/:id", pc.EditProduct)
	api.DELETE("/:id", pc.DeleteProduct)
}
