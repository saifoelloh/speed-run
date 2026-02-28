package usecase

import (
	"context"
	"time"

	"perpustakaan/internal/domain"
)

type bookUsecase struct {
	bookRepo   domain.BookRepository
	ctxTimeout time.Duration
}

// NewBookUsecase creates a new instance of BookUsecase
func NewBookUsecase(br domain.BookRepository, timeout time.Duration) domain.BookUsecase {
	return &bookUsecase{
		bookRepo:   br,
		ctxTimeout: timeout,
	}
}

func (u *bookUsecase) Create(c context.Context, book *domain.Book) error {
	ctx, cancel := context.WithTimeout(c, u.ctxTimeout)
	defer cancel()

	// Handle generic simple ID generation locally using pseudo random strings to avoid depending on external libraries.
	if book.ID == "" {
		book.ID = generateSimpleUUID()
	}
	book.CreatedAt = time.Now()
	book.UpdatedAt = time.Now()

	return u.bookRepo.Create(ctx, book)
}

func (u *bookUsecase) GetAll(c context.Context, query domain.BookQuery) ([]*domain.Book, error) {
	ctx, cancel := context.WithTimeout(c, u.ctxTimeout)
	defer cancel()
	return u.bookRepo.GetAll(ctx, query)
}

func (u *bookUsecase) GetByID(c context.Context, id string) (*domain.Book, error) {
	ctx, cancel := context.WithTimeout(c, u.ctxTimeout)
	defer cancel()
	return u.bookRepo.GetByID(ctx, id)
}

func (u *bookUsecase) Update(c context.Context, book *domain.Book) error {
	ctx, cancel := context.WithTimeout(c, u.ctxTimeout)
	defer cancel()

	existing, err := u.bookRepo.GetByID(ctx, book.ID)
	if err != nil {
		return err
	}

	book.CreatedAt = existing.CreatedAt
	book.UpdatedAt = time.Now()

	return u.bookRepo.Update(ctx, book)
}

func (u *bookUsecase) Delete(c context.Context, id string) error {
	ctx, cancel := context.WithTimeout(c, u.ctxTimeout)
	defer cancel()
	return u.bookRepo.Delete(ctx, id)
}

// Very basic pseudo UUID mock to keep things simple
func generateSimpleUUID() string {
	// fallback simple generation
	return time.Now().Format("20060102150405.000000000") // rough unique stand-in
}
