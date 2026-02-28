package main

import (
	"log"
	"time" // Added for time.Duration

	"perpustakaan/internal/config"
	httpDelivery "perpustakaan/internal/delivery/http" // Fixed import case
	"perpustakaan/internal/repository/inmemory"
	"perpustakaan/internal/usecase"
	"perpustakaan/pkg/jwt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Very basic pseudo UUID mock to keep things simple
func generateSimpleUUID() string {
	// fallback simple generation
	return time.Now().Format("20060102150405.000000000") // rough unique stand-in
}

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	cfg := config.LoadConfig()

	// Initialize basic routes
	httpDelivery.NewBasicHandler(e)

	tokenMaker := jwt.NewJWTTokenMaker(cfg.JWTSecret, cfg.JWTExpireHours)
	httpDelivery.NewAuthHandler(e, tokenMaker)

	// Memory Repos
	bookRepo := inmemory.NewBookRepository()
	bookUsecase := usecase.NewBookUsecase(bookRepo, time.Duration(2)*time.Second)

	httpDelivery.NewBookHandler(e, bookUsecase, tokenMaker)

	log.Printf("Server configured on port %s", cfg.AppPort)
	log.Fatal(e.Start(":" + cfg.AppPort))
}
