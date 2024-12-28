package router

import (
	"github.com/labstack/echo/v4"
	"github.com/ventry/internal/pkg/auth"
)

func AuthRoutes(e *echo.Echo, handler *auth.AuthHandler, authService *auth.AuthService) {
	api := e.Group("/auth")
	api.POST("/register", handler.Register)
	api.POST("/login", handler.Login)

	protected := e.Group("/auth")
	protected.Use(auth.AuthMiddleware(authService))
	protected.GET("/user-profile", handler.GetProfile)
	protected.GET("/check-token", handler.CheckTokenExpiration)

}
