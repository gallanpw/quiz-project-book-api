package controllers

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"quiz-project-book-api/config"
	"quiz-project-book-api/models"

	"github.com/gin-gonic/gin"
)

// GetAllBooks
func GetAllBooks(c *gin.Context) {
	var books []models.Book
	rows, err := config.DB.Query("SELECT id, title, description, image_url, release_year, price, total_page, thickness, category_id, created_at, created_by, modified_at, modified_by FROM books WHERE deleted_at IS NULL ORDER BY id ASC")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch books"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var book models.Book
		if err := rows.Scan(&book.ID, &book.Title, &book.Description, &book.ImageURL, &book.ReleaseYear, &book.Price, &book.TotalPage, &book.Thickness, &book.CategoryID, &book.CreatedAt, &book.CreatedBy, &book.ModifiedAt, &book.ModifiedBy); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan book"})
			return
		}
		books = append(books, book)
	}

	c.JSON(http.StatusOK, books)
}

// GetBookByID
func GetBookByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	var book models.Book
	err = config.DB.QueryRow("SELECT id, title, description, image_url, release_year, price, total_page, thickness, category_id, created_at, created_by, modified_at, modified_by FROM books WHERE id = $1 AND deleted_at IS NULL", id).
		Scan(&book.ID, &book.Title, &book.Description, &book.ImageURL, &book.ReleaseYear, &book.Price, &book.TotalPage, &book.Thickness, &book.CategoryID, &book.CreatedAt, &book.CreatedBy, &book.ModifiedAt, &book.ModifiedBy)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch book"})
		return
	}

	c.JSON(http.StatusOK, book)
}

// CreateBook
func CreateBook(c *gin.Context) {
	var book models.Book
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Mengambil username dari JWT
	username, _ := c.MustGet("username").(string)

	// Logic konversi untuk thickness
	if book.TotalPage > 100 {
		book.Thickness = "tebal"
	} else {
		book.Thickness = "tipis"
	}

	book.CreatedBy = sql.NullString{String: username, Valid: true}
	book.CreatedAt = sql.NullTime{Time: time.Now(), Valid: true}
	book.ModifiedAt = sql.NullTime{Time: time.Time{}, Valid: false}
	book.ModifiedBy = sql.NullString{String: "", Valid: false}

	err := config.DB.QueryRow("INSERT INTO books (title, description, image_url, release_year, price, total_page, thickness, category_id, created_by) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id",
		book.Title, book.Description, book.ImageURL, book.ReleaseYear, book.Price, book.TotalPage, book.Thickness, book.CategoryID, book.CreatedBy).
		Scan(&book.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create book"})
		return
	}

	c.JSON(http.StatusCreated, book)
}

// DeleteBook
func DeleteBook(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	// Mengambil username dari JWT
	username, _ := c.MustGet("username").(string)

	// Menggunakan UPDATE untuk soft delete
	result, err := config.DB.Exec("UPDATE books SET deleted_at = $1, deleted_by = $2 WHERE id = $3 AND deleted_at IS NULL", time.Now(), username, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to soft delete book"})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get rows affected"})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found or already deleted"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Book soft deleted successfully"})
}

// UpdateBook
func UpdateBook(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	var book models.Book
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	username, _ := c.MustGet("username").(string)

	// Logic konversi thickness
	if book.TotalPage > 100 {
		book.Thickness = "tebal"
	} else {
		book.Thickness = "tipis"
	}

	// Menambahkan validasi `deleted_at IS NULL` pada query UPDATE
	result, err := config.DB.Exec("UPDATE books SET title = $1, description = $2, image_url = $3, release_year = $4, price = $5, total_page = $6, thickness = $7, category_id = $8, modified_at = $9, modified_by = $10 WHERE id = $11 AND deleted_at IS NULL",
		book.Title, book.Description, book.ImageURL, book.ReleaseYear, book.Price, book.TotalPage, book.Thickness, book.CategoryID, time.Now(), username, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update book"})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get rows affected"})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found or already deleted"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Book updated successfully"})
}
