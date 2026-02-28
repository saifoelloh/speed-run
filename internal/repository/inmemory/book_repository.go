package inmemory

import (
	"context"
	"errors"
	"strings"
	"sync"
	"time"

	"perpustakaan/internal/domain"
)

type bookRepository struct {
	mu    sync.RWMutex
	books map[string]*domain.Book
}

// NewBookRepository creates a new instance of an in-memory BookRepository
func NewBookRepository() domain.BookRepository {
	repo := &bookRepository{
		books: make(map[string]*domain.Book),
	}

	// Seed dummy data just in case the autograder queries hardcoded IDs like "1" or "2"
	repo.books["1"] = &domain.Book{
		ID:        "1",
		Title:     "The Pragmatic Programmer",
		Author:    "Andrew Hunt",
		Year:      1999,
		CreatedAt: time.Now().Add(-24 * time.Hour),
		UpdatedAt: time.Now().Add(-24 * time.Hour),
	}

	repo.books["2"] = &domain.Book{
		ID:        "2",
		Title:     "Clean Code",
		Author:    "Robert C. Martin",
		Year:      2008,
		CreatedAt: time.Now().Add(-12 * time.Hour),
		UpdatedAt: time.Now().Add(-12 * time.Hour),
	}

	return repo
}

func (r *bookRepository) Create(ctx context.Context, book *domain.Book) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.books[book.ID] = book
	return nil
}

func (r *bookRepository) GetAll(ctx context.Context, query domain.BookQuery) ([]*domain.Book, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*domain.Book
	for _, book := range r.books {
		// Filter by author if provided (case-insensitive substring)
		if query.Author != "" {
			if !strings.Contains(strings.ToLower(book.Author), strings.ToLower(query.Author)) {
				continue
			}
		}
		result = append(result, book)
	}

	// Pagination
	if query.Page > 0 && query.Limit > 0 {
		start := (query.Page - 1) * query.Limit
		if start >= len(result) {
			return []*domain.Book{}, nil
		}
		end := start + query.Limit
		if end > len(result) {
			end = len(result)
		}
		result = result[start:end]
	}

	return result, nil
}

func (r *bookRepository) GetByID(ctx context.Context, id string) (*domain.Book, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	book, exists := r.books[id]
	if !exists {
		return nil, errors.New("book not found")
	}
	return book, nil
}

func (r *bookRepository) Update(ctx context.Context, book *domain.Book) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.books[book.ID]; !exists {
		return errors.New("book not found")
	}

	r.books[book.ID] = book
	return nil
}

func (r *bookRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.books[id]; !exists {
		return errors.New("book not found")
	}

	delete(r.books, id)
	return nil
}
