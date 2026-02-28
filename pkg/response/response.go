package response

import (
	"github.com/labstack/echo/v4"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Success is a helper for standard success responses
func Success(c echo.Context, status int, message string, data interface{}) error {
	return c.JSON(status, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// Error is a helper for standard error responses
func Error(c echo.Context, status int, message string) error {
	return c.JSON(status, Response{
		Success: false,
		Message: message,
	})
}
