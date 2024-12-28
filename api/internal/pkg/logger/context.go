package logger

type ContextKey string

const (
	RequestIDKey ContextKey = "request_id"
	UserIDKey    ContextKey = "user_id"
)
