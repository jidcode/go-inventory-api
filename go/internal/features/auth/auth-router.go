package auth

import (
	"github.com/labstack/echo/v4"
)

func AuthRouter(e *echo.Echo, handler *AuthHandler, authService *AuthService) {
	auth := e.Group("/auth")
	auth.POST("/register", handler.Register)
	auth.POST("/login", handler.Login)
	auth.GET("/token/check", handler.CheckTokenExpiration)

	user := e.Group("/user")
	user.Use(AuthMiddleware(authService))
	user.GET("/profile", handler.GetProfile)
}
