package database

import (
	"fmt"
	"log"

	"perpustakaan/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// NewPostgresDB establishes a connection to PostgreSQL database
func NewPostgresDB(cfg *config.Config) *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	return db
}
