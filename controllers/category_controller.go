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

// GetAllCategories
func GetAllCategories(c *gin.Context) {
	var categories []models.Category
	rows, err := config.DB.Query("SELECT id, name, created_at, created_by, modified_at, modified_by FROM categories WHERE deleted_at IS NULL ORDER BY id ASC")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch categories"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var category models.Category
		if err := rows.Scan(&category.ID, &category.Name, &category.CreatedAt, &category.CreatedBy, &category.ModifiedAt, &category.ModifiedBy); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan category"})
			return
		}
		categories = append(categories, category)
	}

	c.JSON(http.StatusOK, categories)
}

// GetCategoryByID
func GetCategoryByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	var category models.Category
	err = config.DB.QueryRow("SELECT id, name, created_at, created_by, modified_at, modified_by FROM categories WHERE id = $1 AND deleted_at IS NULL", id).
		Scan(&category.ID, &category.Name, &category.CreatedAt, &category.CreatedBy, &category.ModifiedAt, &category.ModifiedBy)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch category"})
		return
	}

	c.JSON(http.StatusOK, category)
}

// CreateCategory
func CreateCategory(c *gin.Context) {
	var category models.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Mengambil username dari JWT
	username, _ := c.MustGet("username").(string)

	category.CreatedBy = sql.NullString{String: username, Valid: true}
	category.CreatedAt = sql.NullTime{Time: time.Now(), Valid: true}
	category.ModifiedAt = sql.NullTime{Time: time.Time{}, Valid: false}
	category.ModifiedBy = sql.NullString{String: "", Valid: false}

	err := config.DB.QueryRow("INSERT INTO categories (name, created_by) VALUES ($1, $2) RETURNING id", category.Name, category.CreatedBy).Scan(&category.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create category"})
		return
	}

	c.JSON(http.StatusCreated, category)
}

// UpdateCategory memperbarui data kategori berdasarkan ID
func UpdateCategory(c *gin.Context) {
	// 1. Mengambil ID dari URL
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	// 2. Mengikat JSON input ke struct Category
	var category models.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 3. Mengambil username dari JWT
	username, _ := c.MustGet("username").(string)

	// 4. Menjalankan query UPDATE
	result, err := config.DB.Exec("UPDATE categories SET name = $1, modified_at = $2, modified_by = $3 WHERE id = $4 AND deleted_at IS NULL",
		category.Name, time.Now(), username, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update category"})
		return
	}

	// 5. Memeriksa apakah ada data yang terpengaruh
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get rows affected"})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found or already deleted"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Category updated successfully"})
}

// DeleteCategory
func DeleteCategory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	// Mengambil username dari JWT
	username, _ := c.MustGet("username").(string)

	// Menggunakan UPDATE untuk soft delete
	result, err := config.DB.Exec("UPDATE categories SET deleted_at = $1, deleted_by = $2 WHERE id = $3 AND deleted_at IS NULL", time.Now(), username, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to soft delete category"})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get rows affected"})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found or already deleted"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Category soft deleted successfully"})
}

// GetBooksByCategory
func GetBooksByCategory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	var books []models.Book
	rows, err := config.DB.Query("SELECT id, title, description, image_url, release_year, price, total_page, thickness FROM books WHERE category_id = $1 AND deleted_at IS NULL", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch books"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var book models.Book
		if err := rows.Scan(&book.ID, &book.Title, &book.Description, &book.ImageURL, &book.ReleaseYear, &book.Price, &book.TotalPage, &book.Thickness); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan book"})
			return
		}
		books = append(books, book)
	}

	c.JSON(http.StatusOK, books)
}
