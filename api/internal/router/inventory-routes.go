package router

import (
	"github.com/labstack/echo/v4"
	"github.com/ventry/internal/features/inventories"
	"github.com/ventry/internal/pkg/auth"
)

func InventoryRoutes(e *echo.Echo, ic inventories.InventoryController, authService auth.AuthService) {
	api := e.Group("/api/inventories")
	api.Use(auth.AuthMiddleware(&authService), auth.RoleMiddleware("user"))

	api.GET("", ic.ListInventories)
	api.GET("/:id", ic.GetInventory)
	api.POST("", ic.CreateInventory)
	api.PUT("/:id", ic.EditInventory)
	api.DELETE("/:id", ic.DeleteInventory)
}
