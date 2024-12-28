package router

import (
	"github.com/labstack/echo/v4"
	"github.com/ventry/internal/features/sales"
	"github.com/ventry/internal/pkg/auth"
)

func SaleRoutes(e *echo.Echo, sc sales.SaleController, authService auth.AuthService) {
	api := e.Group("/api/sales")
	api.Use(auth.AuthMiddleware(&authService), auth.RoleMiddleware("user"))

	api.GET("/inventory/:inventoryId", sc.ListSales)
	api.GET("/:id", sc.GetSale)
	api.POST("", sc.CreateSale)
	api.PUT("/:id", sc.EditSale)
	api.DELETE("/:id", sc.DeleteSale)

	items := api.Group("/:id/items")
	items.POST("", sc.AddItemToSale)
	items.DELETE("/:id/items/:itemId", sc.RemoveItemFromSale)
}
