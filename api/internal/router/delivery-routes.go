package router

import (
	"github.com/labstack/echo/v4"
	"github.com/ventry/internal/features/deliveries"
	"github.com/ventry/internal/pkg/auth"
)

func DeliveryRoutes(e *echo.Echo, dc deliveries.DeliveryController, authService auth.AuthService) {
	api := e.Group("/api/deliveries")
	api.Use(auth.AuthMiddleware(&authService), auth.RoleMiddleware("user"))

	api.GET("/inventory/:inventoryId", dc.ListDeliveries)
	api.GET("/:id", dc.GetDelivery)
	api.POST("", dc.CreateDelivery)
	api.PUT("/:id", dc.UpdateDelivery)
	api.DELETE("/:id", dc.DeleteDelivery)
}
