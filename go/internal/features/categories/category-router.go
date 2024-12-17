package categories

import (
	"github.com/labstack/echo/v4"
	"github.com/ventry/internal/features/auth"
)

func CategoryRouter(e *echo.Echo, handler *CategoryHandler, authService auth.AuthService) {
	api := e.Group("/api/categories")
	api.Use(auth.AuthMiddleware(&authService), auth.RoleMiddleware("user"))

	api.GET("/inventory/:inventoryId", handler.ListCategories)
	api.GET("/:id", handler.GetCategory)
	api.POST("", handler.CreateCategory)
	api.PUT("/:id", handler.UpdateCategory)
	api.DELETE("/:id", handler.DeleteCategory)
}
