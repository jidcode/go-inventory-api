package logger

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/ventry/internal/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log *zap.Logger

// Field represents a log field
type Field struct {
	Key   string
	Value interface{}
}

// Init initializes the logger
func Init(env string) error {
	var config zap.Config

	if env == "production" {
		config = zap.NewProductionConfig()
		config.EncoderConfig.TimeKey = "timestamp"
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	} else {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	var err error
	log, err = config.Build()
	if err != nil {
		return fmt.Errorf("failed to initialize logger: %v", err)
	}

	return nil
}

// WithContext adds context fields to the log entry
func WithContext(ctx context.Context, fields ...Field) []zap.Field {
	zapFields := make([]zap.Field, 0)

	// Add request ID if present
	if reqID := ctx.Value(RequestIDKey); reqID != nil {
		zapFields = append(zapFields, zap.String("request_id", reqID.(string)))
	}

	// Add user ID if present
	if userID := ctx.Value(UserIDKey); userID != nil {
		zapFields = append(zapFields, zap.String("user_id", userID.(string)))
	}

	// Add custom fields
	for _, field := range fields {
		zapFields = append(zapFields, zap.Any(field.Key, field.Value))
	}

	return zapFields
}

// Error logs an error with context
func Error(ctx context.Context, err error, msg string, fields ...Field) {
	if appErr, ok := err.(*errors.AppError); ok {
		fields = append(fields, Field{Key: "error_type", Value: appErr.Type})
		fields = append(fields, Field{Key: "error_code", Value: appErr.Code})
		fields = append(fields, Field{Key: "stack_trace", Value: appErr.Stack})
		if appErr.Operation != "" {
			fields = append(fields, Field{Key: "operation", Value: appErr.Operation})
		}
	}

	log.Error(msg, WithContext(ctx, fields...)...)
}

// Info logs an info message with context
func Info(ctx context.Context, msg string, fields ...Field) {
	log.Info(msg, WithContext(ctx, fields...)...)
}

// Debug logs a debug message with context
func Debug(ctx context.Context, msg string, fields ...Field) {
	log.Debug(msg, WithContext(ctx, fields...)...)
}

// Fatal logs a fatal message and exits
func Fatal(ctx context.Context, msg string, fields ...Field) {
	log.Fatal(msg, WithContext(ctx, fields...)...)
	os.Exit(1)
}

// LogRequest logs incoming HTTP requests
func LogRequest(ctx context.Context, method, path string, duration time.Duration, statusCode int) {
	Info(ctx, "HTTP Request",
		Field{Key: "method", Value: method},
		Field{Key: "path", Value: path},
		Field{Key: "duration_ms", Value: duration.Milliseconds()},
		Field{Key: "status", Value: statusCode},
	)
}

// Sync flushes any buffered log entries
func Sync() error {
	return log.Sync()
}
