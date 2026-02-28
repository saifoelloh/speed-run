package middleware

import (
	"net/http"

	"perpustakaan/pkg/jwt"
	"perpustakaan/pkg/response"

	"github.com/labstack/echo/v4"
)

// AuthMiddleware is the middleware function
func AuthMiddleware(tokenMaker jwt.TokenMaker) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userID, err := tokenMaker.ExtractUserID(c)
			if err != nil {
				return response.Error(c, http.StatusUnauthorized, err.Error())
			}

			// Set user_id in context for subsequent handlers
			c.Set("user_id", userID)
			return next(c)
		}
	}
}
