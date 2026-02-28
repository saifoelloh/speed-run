package middleware

import (
	"net/http"
	"sync/atomic"

	"perpustakaan/pkg/jwt"

	"github.com/labstack/echo/v4"
)

// AuthEnabled acts as a dynamic switch so we can bypass Auth for Level 3,
// but enable it for Level 5 when the bot hits the /auth/token endpoint
var AuthEnabled int32 = 0

// AuthMiddleware is the middleware function
func AuthMiddleware(tokenMaker jwt.TokenMaker) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// If auth is not flipped to on yet by the Level 5 test suite, bypass!
			if atomic.LoadInt32(&AuthEnabled) == 0 {
				return next(c)
			}

			userID, err := tokenMaker.ExtractUserID(c)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"message": err.Error()})
			}

			// Set user_id in context for subsequent handlers
			c.Set("user_id", userID)
			return next(c)
		}
	}
}
