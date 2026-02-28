package http

import (
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
)

type BasicHandler struct{}

func NewBasicHandler(e *echo.Echo) {
	handler := &BasicHandler{}

	e.GET("/ping", handler.Ping)
	e.POST("/echo", handler.Echo)
}

func (h *BasicHandler) Ping(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
	})
}

func (h *BasicHandler) Echo(c echo.Context) error {
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.String(http.StatusBadRequest, "Failed to read body")
	}
	defer c.Request().Body.Close()

	if len(body) == 0 {
		return c.JSON(http.StatusOK, map[string]interface{}{})
	}

	return c.Blob(http.StatusOK, "application/json", body)
}
