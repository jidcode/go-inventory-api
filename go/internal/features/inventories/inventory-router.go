package inventories

import (
	"github.com/labstack/echo/v4"
	"github.com/ventry/internal/features/auth"
)

func InventoryRouter(e *echo.Echo, handler InventoryHandler, authService auth.AuthService) {
	api := e.Group("/api/inventories")
	api.Use(auth.AuthMiddleware(&authService), auth.RoleMiddleware("user"))

	api.GET("", handler.ListInventories)
	api.GET("/:id", handler.GetInventory)
	api.POST("", handler.CreateInventory)
	api.PUT("/:id", handler.UpdateInventory)
	api.DELETE("/:id", handler.DeleteInventory)
}
