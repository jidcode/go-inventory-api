package logger

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func RequestMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		start := time.Now()

		// Create a context with request ID
		requestID := c.Request().Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}
		ctx := context.WithValue(c.Request().Context(), RequestIDKey, requestID)

		// Add user ID to context if authenticated
		if user := c.Get("user"); user != nil {
			// Assuming user has an ID field
			ctx = context.WithValue(ctx, UserIDKey, user.(string))
		}

		// Store the enhanced context in Echo's context
		c.SetRequest(c.Request().WithContext(ctx))

		err := next(c)

		// Log the request after completion
		LogRequest(
			ctx,
			c.Request().Method,
			c.Request().URL.Path,
			time.Since(start),
			c.Response().Status,
		)

		return err
	}
}
