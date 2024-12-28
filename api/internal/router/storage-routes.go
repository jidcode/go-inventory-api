package router

import (
	"github.com/labstack/echo/v4"
	"github.com/ventry/internal/features/storages"
	"github.com/ventry/internal/pkg/auth"
)

func StorageRoutes(e *echo.Echo, sc storages.StorageController, authService auth.AuthService) {
	api := e.Group("/api/storages")
	api.Use(auth.AuthMiddleware(&authService), auth.RoleMiddleware("user"))

	api.GET("/inventory/:inventoryId", sc.ListStorages)
	api.GET("/:id", sc.GetStorage)
	api.POST("", sc.CreateStorage)
	api.PUT("/:id", sc.EditStorage)
	api.DELETE("/:id", sc.DeleteStorage)

	units := api.Group("/:storageId/units")
	units.GET("", sc.ListUnits)
	units.GET("/:id", sc.GetUnit)
	units.POST("", sc.CreateUnit)
	units.PUT("/:id", sc.EditUnit)
	units.DELETE("/:id", sc.DeleteUnit)
}
