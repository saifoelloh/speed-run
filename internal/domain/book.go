package domain

import (
	"context"
	"time"
)

// Book represents the core entity
type Book struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Author    string    `json:"author"`
	Year      int       `json:"year"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// BookQuery contains search and pagination filters
type BookQuery struct {
	Author string
	Page   int
	Limit  int
}

type BookRepository interface {
	Create(ctx context.Context, book *Book) error
	GetAll(ctx context.Context, query BookQuery) ([]*Book, error)
	GetByID(ctx context.Context, id string) (*Book, error)
	Update(ctx context.Context, book *Book) error
	Delete(ctx context.Context, id string) error
}

type BookUsecase interface {
	Create(ctx context.Context, book *Book) error
	GetAll(ctx context.Context, query BookQuery) ([]*Book, error)
	GetByID(ctx context.Context, id string) (*Book, error)
	Update(ctx context.Context, book *Book) error
	Delete(ctx context.Context, id string) error
}
