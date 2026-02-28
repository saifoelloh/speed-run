package inmemory

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"perpustakaan/internal/domain"
)

type bookRepository struct {
	mu       sync.RWMutex
	filePath string
	books    map[string]*domain.Book
}

// NewBookRepository creates a new instance of an in-memory BookRepository
// and attempts to load existing data from data.json
func NewBookRepository() domain.BookRepository {
	repo := &bookRepository{
		filePath: "data.json",
		books:    make(map[string]*domain.Book),
	}

	repo.loadFromFile()

	// If it's completely empty after load, we can seed the dummy data
	if len(repo.books) == 0 {
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
		repo.saveToFile()
	}

	return repo
}

// saveToFile writes the current map to a JSON file
func (r *bookRepository) saveToFile() {
	data, err := json.MarshalIndent(r.books, "", "  ")
	if err == nil {
		_ = ioutil.WriteFile(r.filePath, data, 0644)
	}
}

// loadFromFile reads the map from a JSON file if it exists
func (r *bookRepository) loadFromFile() {
	if _, err := os.Stat(r.filePath); os.IsNotExist(err) {
		return
	}

	data, err := ioutil.ReadFile(r.filePath)
	if err == nil {
		_ = json.Unmarshal(data, &r.books)
	}
}

func (r *bookRepository) Create(ctx context.Context, book *domain.Book) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.books[book.ID] = book
	r.saveToFile()
	return nil
}

func (r *bookRepository) GetAll(ctx context.Context, query domain.BookQuery) ([]*domain.Book, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*domain.Book
	for _, book := range r.books {
		if query.Author != "" {
			if !strings.Contains(strings.ToLower(book.Author), strings.ToLower(query.Author)) {
				continue
			}
		}
		result = append(result, book)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].CreatedAt.Before(result[j].CreatedAt)
	})

	if query.Limit > 0 {
		page := query.Page
		if page < 1 {
			page = 1
		}
		start := (page - 1) * query.Limit
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
	r.saveToFile()
	return nil
}

func (r *bookRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.books[id]; !exists {
		return errors.New("book not found")
	}

	delete(r.books, id)
	r.saveToFile()
	return nil
}
