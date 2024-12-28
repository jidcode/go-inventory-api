package errors

import (
	"github.com/labstack/echo/v4"
)

// HTTPError handles converting AppError to HTTP responses
func Send(c echo.Context, err error) error {
	if err == nil {
		return nil
	}

	appErr, ok := err.(*AppError)
	if !ok {
		// If it's not our custom error, wrap it as an internal error
		appErr = InternalError(err, "An unexpected error occurred")
	}

	return c.JSON(appErr.Code, appErr)
}
