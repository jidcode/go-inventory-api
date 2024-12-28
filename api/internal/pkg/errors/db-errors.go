package errors

import "strings"

// MapPostgresError maps specific PostgreSQL errors to AppError
func MapPostgresError(err error, operation string) *AppError {
	if err == nil {
		return nil
	}

	details := err.Error()

	switch {
	case strings.Contains(details, "duplicate key value violates unique constraint"):
		return &AppError{
			Type:      DatabaseErr,
			Message:   "Duplicate entry: item already exists",
			Details:   details,
			Operation: operation,
			Code:      400,
		}
	case strings.Contains(details, "violates foreign key constraint"):
		return &AppError{
			Type:      DatabaseErr,
			Message:   "Invalid reference: related item not found",
			Details:   details,
			Operation: operation,
			Code:      400,
		}
	case strings.Contains(details, "value too long for type"):
		return &AppError{
			Type:      DatabaseErr,
			Message:   "Input exceeds allowed length",
			Details:   details,
			Operation: operation,
			Code:      400,
		}
	case strings.Contains(details, "invalid input syntax"):
		return &AppError{
			Type:      DatabaseErr,
			Message:   "Invalid format: check input syntax",
			Details:   details,
			Operation: operation,
			Code:      400,
		}
	case strings.Contains(details, "null value in column"):
		return &AppError{
			Type:      DatabaseErr,
			Message:   "Required field is missing",
			Details:   details,
			Operation: operation,
			Code:      400,
		}
	case strings.Contains(details, "division by zero"):
		return &AppError{
			Type:      DatabaseErr,
			Message:   "Math error: division by zero",
			Details:   details,
			Operation: operation,
			Code:      400,
		}
	case strings.Contains(details, "out of range"):
		return &AppError{
			Type:      DatabaseErr,
			Message:   "Value out of range",
			Details:   details,
			Operation: operation,
			Code:      400,
		}
	case strings.Contains(details, "could not serialize access due to concurrent update"):
		return &AppError{
			Type:      DatabaseErr,
			Message:   "Concurrent update conflict",
			Details:   details,
			Operation: operation,
			Code:      409, // Conflict
		}
	default:
		return &AppError{
			Type:      DatabaseErr,
			Message:   "Database error occurred",
			Details:   details,
			Operation: operation,
			Code:      500,
		}
	}
}
