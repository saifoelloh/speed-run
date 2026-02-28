package usecase

import (
	"context"
	"crypto/rand"
	"fmt"
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

	if book.Title != "" {
		existing.Title = book.Title
	}
	if book.Author != "" {
		existing.Author = book.Author
	}
	if book.Year != 0 {
		existing.Year = book.Year
	}
	existing.UpdatedAt = time.Now()

	return u.bookRepo.Update(ctx, existing)
}

func (u *bookUsecase) Delete(c context.Context, id string) error {
	ctx, cancel := context.WithTimeout(c, u.ctxTimeout)
	defer cancel()
	return u.bookRepo.Delete(ctx, id)
}

// generateSimpleUUID generates a pseudo-UUID v4 string without external dependencies
func generateSimpleUUID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	// Set standard UUID v4 bits
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}
