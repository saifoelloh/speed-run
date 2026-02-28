package http

import (
	"net/http"

	"perpustakaan/pkg/jwt"

	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	tokenMaker jwt.TokenMaker
}

// NewAuthHandler initializes the auth endpoint for Level 5
func NewAuthHandler(e *echo.Echo, tm jwt.TokenMaker) {
	handler := &AuthHandler{
		tokenMaker: tm,
	}

	// Used simply to bypass Level 5 restrictions quickly.
	e.POST("/auth/token", handler.GenerateToken)
}

func (h *AuthHandler) GenerateToken(c echo.Context) error {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "invalid request"})
	}

	if req.Username != "admin" || req.Password != "password" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "invalid credentials"})
	}

	// Generate token for the valid user
	token, err := h.tokenMaker.CreateToken("admin-uuid")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "failed generating token"})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"token": token,
	})
}
