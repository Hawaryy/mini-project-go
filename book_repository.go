package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/yourusername/book-management-api/models"
)

type BookRepository struct {
	db *sql.DB
}

func NewBookRepository(db *sql.DB) *BookRepository {
	return &BookRepository{db: db}
}

func (r *BookRepository) GetAll() ([]models.Book, error) {
	rows, err := r.db.Query(`
		SELECT id, title, description, image_url, release_year, price, total_page, thickness, category_id, 
			   created_at, created_by, modified_at, modified_by 
		FROM books
		ORDER BY id ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []models.Book
	for rows.Next() {
		var book models.Book
		err := rows.Scan(&book.ID, &book.Title, &book.Description, &book.ImageURL, &book.ReleaseYear,
			&book.Price, &book.TotalPage, &book.Thickness, &book.CategoryID,
			&book.CreatedAt, &book.CreatedBy, &book.ModifiedAt, &book.ModifiedBy)
		if err != nil {
			return nil, err
		}
		books = append(books, book)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return books, nil
}

func (r *BookRepository) GetByID(id int) (*models.Book, error) {
	book := &models.Book{}
	err := r.db.QueryRow(`
		SELECT id, title, description, image_url, release_year, price, total_page, thickness, category_id, 
			   created_at, created_by, modified_at, modified_by 
		FROM books
		WHERE id = $1
	`, id).Scan(&book.ID, &book.Title, &book.Description, &book.ImageURL, &book.ReleaseYear,
		&book.Price, &book.TotalPage, &book.Thickness, &book.CategoryID,
		&book.CreatedAt, &book.CreatedBy, &book.ModifiedAt, &book.ModifiedBy)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("book not found")
	}
	if err != nil {
		return nil, err
	}

	return book, nil
}

func (r *BookRepository) Create(title, description, imageURL string, releaseYear, price, totalPage, categoryID int, createdBy string) (*models.Book, error) {
	// Calculate thickness based on total_page
	thickness := "tipis"
	if totalPage >= 100 {
		thickness = "tebal"
	}

	book := &models.Book{
		Title:       title,
		Description: description,
		ImageURL:    imageURL,
		ReleaseYear: releaseYear,
		Price:       price,
		TotalPage:   totalPage,
		Thickness:   thickness,
		CategoryID:  categoryID,
		CreatedAt:   time.Now(),
		CreatedBy:   createdBy,
		ModifiedAt:  time.Now(),
		ModifiedBy:  createdBy,
	}

	err := r.db.QueryRow(`
		INSERT INTO books (title, description, image_url, release_year, price, total_page, thickness, category_id, created_by, modified_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, created_at, modified_at
	`, title, description, imageURL, releaseYear, price, totalPage, thickness, categoryID, createdBy, createdBy).
		Scan(&book.ID, &book.CreatedAt, &book.ModifiedAt)

	if err != nil {
		return nil, err
	}

	return book, nil
}

func (r *BookRepository) Delete(id int) error {
	// Check if book exists
	_, err := r.GetByID(id)
	if err != nil {
		return err
	}

	result, err := r.db.Exec("DELETE FROM books WHERE id = $1", id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("book not found")
	}

	return nil
}

func (r *BookRepository) GetByCategory(categoryID int) ([]models.Book, error) {
	rows, err := r.db.Query(`
		SELECT id, title, description, image_url, release_year, price, total_page, thickness, category_id, 
			   created_at, created_by, modified_at, modified_by 
		FROM books
		WHERE category_id = $1
		ORDER BY id ASC
	`, categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []models.Book
	for rows.Next() {
		var book models.Book
		err := rows.Scan(&book.ID, &book.Title, &book.Description, &book.ImageURL, &book.ReleaseYear,
			&book.Price, &book.TotalPage, &book.Thickness, &book.CategoryID,
			&book.CreatedAt, &book.CreatedBy, &book.ModifiedAt, &book.ModifiedBy)
		if err != nil {
			return nil, err
		}
		books = append(books, book)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return books, nil
}
