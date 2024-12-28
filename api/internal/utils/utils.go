package utils

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/ventry/internal/pkg/errors"
	"github.com/ventry/internal/pkg/logger"
)

func BindAndValidateInput(ctx echo.Context, input interface{}) error {
	if err := ctx.Bind(input); err != nil {
		errors.Send(ctx, err)
		logger.Error(ctx.Request().Context(), err, "Failed to bind request input",
			logger.Field{Key: "path", Value: ctx.Path()},
			logger.Field{Key: "method", Value: ctx.Request().Method},
		)
		return errors.Wrap(err, errors.BadRequest, "Failed to parse request body", 400)
	}

	if err := ctx.Validate(input); err != nil {
		errors.Send(ctx, err)
		logger.Error(ctx.Request().Context(), err, "Input validation failed",
			logger.Field{Key: "path", Value: ctx.Path()},
			logger.Field{Key: "method", Value: ctx.Request().Method},
			logger.Field{Key: "input_type", Value: fmt.Sprintf("%T", input)},
		)
		return errors.Wrap(err, errors.ValidationErr, "Input validation failed", 400)
	}

	logger.Debug(ctx.Request().Context(), "Successfully validated input",
		logger.Field{Key: "path", Value: ctx.Path()},
		logger.Field{Key: "method", Value: ctx.Request().Method},
		logger.Field{Key: "input_type", Value: fmt.Sprintf("%T", input)},
	)

	return nil
}
