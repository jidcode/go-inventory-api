package router

import (
	"github.com/labstack/echo/v4"
	"github.com/ventry/internal/features/categories"
	"github.com/ventry/internal/pkg/auth"
)

func CategoryRoutes(e *echo.Echo, cc categories.CategoryController, authService auth.AuthService) {
	api := e.Group("/api/categories")
	api.Use(auth.AuthMiddleware(&authService), auth.RoleMiddleware("user"))

	api.GET("/inventory/:inventoryId", cc.ListCategories)
	api.GET("/:id", cc.GetCategory)
	api.POST("", cc.CreateCategory)
	api.PUT("/:id", cc.EditCategory)
	api.DELETE("/:id", cc.DeleteCategory)
}
