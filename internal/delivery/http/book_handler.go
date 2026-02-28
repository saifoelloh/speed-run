package http

import (
	"net/http"
	"strconv"

	"perpustakaan/internal/delivery/http/middleware"
	"perpustakaan/internal/domain"
	"perpustakaan/pkg/jwt"
	"perpustakaan/pkg/response"

	"github.com/labstack/echo/v4"
)

type BookHandler struct {
	BUsecase domain.BookUsecase
}

// NewBookHandler will initialize the book endpoints
func NewBookHandler(e *echo.Echo, us domain.BookUsecase, tokenMaker jwt.TokenMaker) {
	handler := &BookHandler{
		BUsecase: us,
	}

	// Public routes
	bookGrp := e.Group("/books")

	// Create Book
	// POST /books
	// Level 7: Invalid requests must return 400 or 422 - we use Bind which handles basic ones, but we also manually check
	bookGrp.POST("", handler.Create)

	// List Books (with pagination and search)
	// GET /books
	// GET /books?author=X
	// GET /books?page=1&limit=2
	// ONLY this endpoint is protected according to Level 5: "GET /books (protected)"
	bookGrp.GET("", handler.GetAll, middleware.AuthMiddleware(tokenMaker))

	// Get Book by ID
	// GET /books/:id
	bookGrp.GET("/:id", handler.GetByID)

	// Update Book
	// PUT /books/:id
	bookGrp.PUT("/:id", handler.Update)

	// Delete Book
	// DELETE /books/:id
	bookGrp.DELETE("/:id", handler.Delete)
}

func (h *BookHandler) Create(c echo.Context) error {
	var book domain.Book
	if err := c.Bind(&book); err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid request payload")
	}

	if book.Title == "" || book.Author == "" {
		return response.Error(c, http.StatusBadRequest, "Title and Author are required")
	}

	err := h.BUsecase.Create(c.Request().Context(), &book)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, err.Error())
	}

	return response.Success(c, http.StatusCreated, "Book created successfully", book)
}

func (h *BookHandler) GetAll(c echo.Context) error {
	author := c.QueryParam("author")

	page := 0
	limit := 0

	if pageStr := c.QueryParam("page"); pageStr != "" {
		p, err := strconv.Atoi(pageStr)
		if err == nil {
			page = p
		}
	}

	if limitStr := c.QueryParam("limit"); limitStr != "" {
		l, err := strconv.Atoi(limitStr)
		if err == nil {
			limit = l
		}
	}

	query := domain.BookQuery{
		Author: author,
		Page:   page,
		Limit:  limit,
	}

	books, err := h.BUsecase.GetAll(c.Request().Context(), query)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, err.Error())
	}

	return response.Success(c, http.StatusOK, "Books retrieved successfully", books)
}

func (h *BookHandler) GetByID(c echo.Context) error {
	id := c.Param("id")

	book, err := h.BUsecase.GetByID(c.Request().Context(), id)
	if err != nil {
		if err.Error() == "book not found" {
			return response.Error(c, http.StatusNotFound, "Book not found")
		}
		return response.Error(c, http.StatusInternalServerError, err.Error())
	}

	return response.Success(c, http.StatusOK, "Book retrieved successfully", book)
}

func (h *BookHandler) Update(c echo.Context) error {
	id := c.Param("id")

	var book domain.Book
	if err := c.Bind(&book); err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid request payload")
	}

	book.ID = id
	err := h.BUsecase.Update(c.Request().Context(), &book)
	if err != nil {
		if err.Error() == "book not found" {
			return response.Error(c, http.StatusNotFound, "Book not found")
		}
		return response.Error(c, http.StatusInternalServerError, err.Error())
	}

	updatedBook, _ := h.BUsecase.GetByID(c.Request().Context(), id)
	return response.Success(c, http.StatusOK, "Book updated successfully", updatedBook)
}

func (h *BookHandler) Delete(c echo.Context) error {
	id := c.Param("id")

	err := h.BUsecase.Delete(c.Request().Context(), id)
	if err != nil {
		if err.Error() == "book not found" {
			return response.Error(c, http.StatusNotFound, "Book not found")
		}
		return response.Error(c, http.StatusInternalServerError, err.Error())
	}

	return response.Success(c, http.StatusOK, "Book deleted successfully", nil)
}
