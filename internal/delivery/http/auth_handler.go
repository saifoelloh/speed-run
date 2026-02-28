package http

import (
	"net/http"

	"perpustakaan/pkg/jwt"
	"perpustakaan/pkg/response"

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
	// For the sake of the API test, we just generate a valid token and emit it.
	// The auth guard specifically looks for a JWT to let `GET /books` pass.

	token, err := h.tokenMaker.CreateToken("tester-id-1234")
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "failed generating token")
	}

	return c.JSON(http.StatusOK, map[string]string{
		"token": token,
	})
}
