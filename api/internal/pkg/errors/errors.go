package errors

import (
	"fmt"
	"runtime"
)

type ErrorType string

const (
	NotFound      ErrorType = "NOT_FOUND"
	ValidationErr ErrorType = "VALIDATION_ERROR"
	InternalErr   ErrorType = "INTERNAL_ERROR"
	Unauthorized  ErrorType = "UNAUTHORIZED"
	BadRequest    ErrorType = "BAD_REQUEST"
	Forbidden     ErrorType = "FORBIDDEN"
	DatabaseErr   ErrorType = "DATABASE_ERROR"
)

// AppError represents a custom application error
type AppError struct {
	Type      ErrorType `json:"type"`
	Message   string    `json:"message"`
	Details   string    `json:"details,omitempty"`
	Stack     string    `json:"-"`
	Operation string    `json:"operation,omitempty"`
	Code      int       `json:"code"`
}

func (e *AppError) Error() string {
	return fmt.Sprintf("[%s] %s: %s", e.Type, e.Message, e.Details)
}

// New creates a new AppError with stack trace
func New(errType ErrorType, message string, code int) *AppError {
	stack := make([]byte, 4096)
	runtime.Stack(stack, false)

	return &AppError{
		Type:    errType,
		Message: message,
		Stack:   string(stack),
		Code:    code,
	}
}

// Wrap wraps an existing error with additional context
func Wrap(err error, errType ErrorType, message string, code int) *AppError {
	if err == nil {
		return nil
	}

	stack := make([]byte, 4096)
	runtime.Stack(stack, false)

	return &AppError{
		Type:    errType,
		Message: message,
		Details: err.Error(),
		Stack:   string(stack),
		Code:    code,
	}
}

// Helper functions for common errors
func NotFoundError(message string) *AppError {
	return New(NotFound, message, 404)
}

func ValidationError(message string) *AppError {
	return New(ValidationErr, message, 400)
}

func InternalError(err error, message string) *AppError {
	return Wrap(err, InternalErr, message, 500)
}

func DatabaseError(err error, operation string) *AppError {
	if pgErr := MapPostgresError(err, operation); pgErr != nil {
		return pgErr
	}
	appErr := Wrap(err, DatabaseErr, "Database operation failed", 500)
	appErr.Operation = operation
	return appErr
}
