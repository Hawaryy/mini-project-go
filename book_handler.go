package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/book-management-api/models"
	"github.com/yourusername/book-management-api/repository"
)

type BookHandler struct {
	bookRepo     *repository.BookRepository
	categoryRepo *repository.CategoryRepository
}

func NewBookHandler(bookRepo *repository.BookRepository, categoryRepo *repository.CategoryRepository) *BookHandler {
	return &BookHandler{
		bookRepo:     bookRepo,
		categoryRepo: categoryRepo,
	}
}

func (h *BookHandler) GetAllBooks(c *gin.Context) {
	books, err := h.bookRepo.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: "Failed to fetch books: " + err.Error(),
		})
		return
	}

	if books == nil {
		books = []models.Book{}
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Status:  http.StatusOK,
		Message: "Books retrieved successfully",
		Data:    books,
	})
}

func (h *BookHandler) GetBookByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid book ID",
		})
		return
	}

	book, err := h.bookRepo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Status:  http.StatusNotFound,
			Message: "Book not found",
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Status:  http.StatusOK,
		Message: "Book retrieved successfully",
		Data:    book,
	})
}

func (h *BookHandler) CreateBook(c *gin.Context) {
	var req models.CreateBookRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid request: " + err.Error(),
		})
		return
	}

	// Validate release year
	if req.ReleaseYear < 1980 || req.ReleaseYear > 2024 {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Status:  http.StatusBadRequest,
			Message: "Release year must be between 1980 and 2024",
		})
		return
	}

	// Check if category exists
	_, err := h.categoryRepo.GetByID(req.CategoryID)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Status:  http.StatusNotFound,
			Message: "Category not found",
		})
		return
	}

	username, exists := c.Get("username")
	if !exists {
		username = "unknown"
	}

	book, err := h.bookRepo.Create(
		req.Title,
		req.Description,
		req.ImageURL,
		req.ReleaseYear,
		req.Price,
		req.TotalPage,
		req.CategoryID,
		username.(string),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: "Failed to create book: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, models.SuccessResponse{
		Status:  http.StatusCreated,
		Message: "Book created successfully",
		Data:    book,
	})
}

func (h *BookHandler) DeleteBook(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid book ID",
		})
		return
	}

	err = h.bookRepo.Delete(id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Status:  http.StatusNotFound,
			Message: "Book not found",
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Status:  http.StatusOK,
		Message: "Book deleted successfully",
	})
}

// Handler untuk endpoint /api/categories/:id/books
func (h *BookHandler) GetBooksByCategory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid category ID",
		})
		return
	}

	// Check if category exists
	_, err = h.categoryRepo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Status:  http.StatusNotFound,
			Message: "Category not found",
		})
		return
	}

	books, err := h.bookRepo.GetByCategory(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: "Failed to fetch books: " + err.Error(),
		})
		return
	}

	if books == nil {
		books = []models.Book{}
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Status:  http.StatusOK,
		Message: "Books retrieved successfully",
		Data:    books,
	})
}
