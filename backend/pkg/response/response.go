package response

import (
	"net/http"

	apperrors "erp-system/pkg/errors"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// Meta holds pagination metadata
type Meta struct {
	Page       int   `json:"page"`
	PerPage    int   `json:"per_page"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

// Response is the standard API envelope
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
	Error   *ErrorBody  `json:"error,omitempty"`
}

// ErrorBody holds structured error detail
type ErrorBody struct {
	Code      string `json:"code"`       // error_code dari CustomError
	MessageID string `json:"message_id"` // i18n message key
	Message   string `json:"message"`    // human-readable message
}

func Success(c echo.Context, data interface{}) error {
	return c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    data,
	})
}

func Created(c echo.Context, data interface{}) error {
	return c.JSON(http.StatusCreated, Response{
		Success: true,
		Data:    data,
	})
}

func SuccessWithMeta(c echo.Context, data interface{}, meta Meta) error {
	return c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    data,
		Meta:    &meta,
	})
}

func NoContent(c echo.Context) error {
	return c.NoContent(http.StatusNoContent)
}

// Error returns error response with proper status code and logging
func Error(c echo.Context, err error) error {
	// Extract logger from context if available
	logger, _ := c.Get("logger").(*zap.Logger)
	if logger == nil {
		logger = zap.L()
	}

	// Convert to CustomError if not already
	customErr := apperrors.NewCustomErrorFromError(err)

	// Log error with request context
	customErr.LogError(logger,
		zap.String("request_id", c.Response().Header().Get(echo.HeaderXRequestID)),
		zap.String("path", c.Request().URL.Path),
		zap.String("method", c.Request().Method),
	)

	// Build error response
	return c.JSON(customErr.HTTPCode, Response{
		Success: false,
		Error: &ErrorBody{
			Code:      string(customErr.ErrorCode),
			MessageID: customErr.MessageID,
			Message:   customErr.Message,
		},
	})
}
