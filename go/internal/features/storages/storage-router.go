package storages

import (
	"github.com/labstack/echo/v4"
	"github.com/ventry/internal/features/auth"
)

func StorageRouter(e *echo.Echo, handler *StorageHandler, UnitsHandler *UnitsHandler, authService auth.AuthService) {
	api := e.Group("/api/storages")
	api.Use(auth.AuthMiddleware(&authService), auth.RoleMiddleware("user"))

	api.GET("/inventory/:inventoryId", handler.ListStorages)
	api.GET("/:id", handler.GetStorage)
	api.POST("", handler.CreateStorage)
	api.PUT("/:id", handler.UpdateStorage)
	api.DELETE("/:id", handler.DeleteStorage)

	units := api.Group("/:storageId/units")
	units.GET("", UnitsHandler.ListUnits)
	units.GET("/:id", UnitsHandler.GetUnit)
	units.POST("", UnitsHandler.CreateUnit)
	units.PUT("/:id", UnitsHandler.UpdateUnit)
	units.DELETE("/:id", UnitsHandler.DeleteUnit)
}
