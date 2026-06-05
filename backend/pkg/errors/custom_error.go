package errors

import (
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

// CustomError is a structured application error with HTTP status, error code, and message ID.
type CustomError struct {
	HTTPCode  int       `json:"http_code"`
	ErrorCode ErrorCode `json:"error_code"`
	MessageID string    `json:"message_id"`
	Message   string    `json:"message"`
	Cause     error     `json:"-"` // underlying error, not exposed to client
}

// Error implements the error interface
func (e *CustomError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("[%s] %s: %v", e.ErrorCode, e.Message, e.Cause)
	}
	return fmt.Sprintf("[%s] %s", e.ErrorCode, e.Message)
}

// Unwrap returns the underlying cause for errors.Is/As compatibility
func (e *CustomError) Unwrap() error {
	return e.Cause
}

// ═══════════════════════════════════════════════════════════════
// CONSTRUCTORS
// ═══════════════════════════════════════════════════════════════

// NewCustomError creates a new custom error
func NewCustomError(message string) *CustomError {
	return &CustomError{
		HTTPCode:  http.StatusInternalServerError,
		ErrorCode: ErrCodeInternalServer,
		MessageID: "error_internal_server",
		Message:   message,
	}
}

// NewCustomErrorFromError wraps an existing error into CustomError
func NewCustomErrorFromError(err error) *CustomError {
	if err == nil {
		return nil
	}

	// If it's already a CustomError, return as-is
	if customErr, ok := err.(*CustomError); ok {
		return customErr
	}

	return &CustomError{
		HTTPCode:  http.StatusInternalServerError,
		ErrorCode: ErrCodeInternalServer,
		MessageID: "error_internal_server",
		Message:   err.Error(),
		Cause:     err,
	}
}

// ═══════════════════════════════════════════════════════════════
// BUILDER METHODS (fluent API)
// ═══════════════════════════════════════════════════════════════

// WithHTTPCode sets the HTTP status code
func (e *CustomError) WithHTTPCode(code int) *CustomError {
	e.HTTPCode = code
	return e
}

// WithErrorCode sets the application error code
func (e *CustomError) WithErrorCode(code ErrorCode) *CustomError {
	e.ErrorCode = code
	return e
}

// WithMessageID sets the i18n message ID
func (e *CustomError) WithMessageID(id string) *CustomError {
	e.MessageID = id
	return e
}

// WithMessage sets the error message
func (e *CustomError) WithMessage(msg string) *CustomError {
	e.Message = msg
	return e
}

// WithCause sets the underlying error
func (e *CustomError) WithCause(err error) *CustomError {
	e.Cause = err
	return e
}

// ═══════════════════════════════════════════════════════════════
// LOGGING HELPER
// ═══════════════════════════════════════════════════════════════

// LogError logs the error with zap logger including all context
func (e *CustomError) LogError(logger *zap.Logger, additionalFields ...zap.Field) {
	fields := []zap.Field{
		zap.String("error_code", string(e.ErrorCode)),
		zap.String("message_id", e.MessageID),
		zap.Int("http_code", e.HTTPCode),
		zap.String("message", e.Message),
	}

	if e.Cause != nil {
		fields = append(fields, zap.Error(e.Cause))
	}

	fields = append(fields, additionalFields...)

	// Log level based on HTTP status
	switch {
	case e.HTTPCode >= 500:
		logger.Error("server error", fields...)
	case e.HTTPCode >= 400:
		logger.Warn("client error", fields...)
	default:
		logger.Info("error", fields...)
	}
}

// ═══════════════════════════════════════════════════════════════
// COMMON ERROR FACTORIES (pre-configured)
// ═══════════════════════════════════════════════════════════════

// NotFound creates a 404 error
func NotFound(message string) *CustomError {
	return NewCustomError(message).
		WithHTTPCode(http.StatusNotFound).
		WithErrorCode(ErrCodeResourceNotFound).
		WithMessageID("error_not_found")
}

// BadRequest creates a 400 error
func BadRequest(message string) *CustomError {
	return NewCustomError(message).
		WithHTTPCode(http.StatusBadRequest).
		WithErrorCode(ErrCodeBadRequest).
		WithMessageID("error_bad_request")
}

// Unauthorized creates a 401 error
func Unauthorized(message string) *CustomError {
	return NewCustomError(message).
		WithHTTPCode(http.StatusUnauthorized).
		WithErrorCode(ErrCodeUnauthorized).
		WithMessageID("error_unauthorized")
}

// Forbidden creates a 403 error
func Forbidden(message string) *CustomError {
	return NewCustomError(message).
		WithHTTPCode(http.StatusForbidden).
		WithErrorCode(ErrCodeForbidden).
		WithMessageID("error_forbidden")
}

// Conflict creates a 409 error
func Conflict(message string) *CustomError {
	return NewCustomError(message).
		WithHTTPCode(http.StatusConflict).
		WithErrorCode(ErrCodeResourceConflict).
		WithMessageID("error_conflict")
}

// Internal creates a 500 error
func Internal(message string) *CustomError {
	return NewCustomError(message).
		WithHTTPCode(http.StatusInternalServerError).
		WithErrorCode(ErrCodeInternalServer).
		WithMessageID("error_internal_server")
}

// ═══════════════════════════════════════════════════════════════
// HTTP STATUS HELPER
// ═══════════════════════════════════════════════════════════════

// HTTPStatusFromError extracts HTTP status from error, defaulting to 500
func HTTPStatusFromError(err error) int {
	if err == nil {
		return http.StatusOK
	}

	if customErr, ok := err.(*CustomError); ok {
		return customErr.HTTPCode
	}

	return http.StatusInternalServerError
}
