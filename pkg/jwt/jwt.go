package jwt

import (
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type TokenMaker interface {
	CreateToken(userID string) (string, error)
	ExtractUserID(c echo.Context) (string, error)
}

type JWTTokenMaker struct {
	secretKey   string
	expireHours int
}

// NewJWTTokenMaker creates a new instance of JWTTokenMaker
func NewJWTTokenMaker(secretKey string, expireHours int) TokenMaker {
	return &JWTTokenMaker{
		secretKey:   secretKey,
		expireHours: expireHours,
	}
}

func (maker *JWTTokenMaker) CreateToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * time.Duration(maker.expireHours)).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(maker.secretKey))
}

func (maker *JWTTokenMaker) ExtractUserID(c echo.Context) (string, error) {
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("authorization header missing")
	}

	tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(maker.secretKey), nil
	})

	if err != nil || !token.Valid {
		return "", errors.New("invalid or expired token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid token claims")
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", errors.New("invalid user_id claim in token")
	}

	return userID, nil
}
