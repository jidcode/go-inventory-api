package auth

import (
	"log"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/ventry/internal/domain"
)

func AuthMiddleware(service *AuthService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			authHeader := ctx.Request().Header.Get("Authorization")
			if authHeader == "" {
				return ctx.JSON(http.StatusUnauthorized, "missing auth token")
			}

			tokenParts := strings.Split(authHeader, " ")
			if len(tokenParts) != 2 || strings.ToLower(tokenParts[0]) != "bearer" {
				return ctx.JSON(http.StatusUnauthorized, "invalid auth token format")
			}

			token := tokenParts[1]
			user, err := service.GetUserFromToken(token)
			if err != nil {
				log.Printf("Error retrieving user from token: %v", err)
				return ctx.JSON(http.StatusUnauthorized, "invalid auth token")
			}

			ctx.Set("user", user)
			return next(ctx)
		}
	}
}

func RoleMiddleware(roles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			user, ok := ctx.Get("user").(*domain.User)
			if !ok {
				return ctx.JSON(http.StatusUnauthorized, "user not found in context")
			}

			for _, role := range roles {
				if user.Role == role {
					return next(ctx)
				}
			}

			return ctx.JSON(http.StatusForbidden, "Access forbidden: insufficient permissions")
		}
	}
}
