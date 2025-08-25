package models

import "database/sql"

type Book struct {
	ID          int            `json:"id"`
	Title       string         `json:"title" binding:"required"`
	Description string         `json:"description"`
	ImageURL    string         `json:"image_url"`
	ReleaseYear int            `json:"release_year" binding:"required,gte=1980,lte=2024"`
	Price       int            `json:"price" binding:"required"`
	TotalPage   int            `json:"total_page" binding:"required"`
	Thickness   string         `json:"thickness"`
	CategoryID  int            `json:"category_id" binding:"required"`
	CreatedAt   sql.NullTime   `json:"created_at"`
	CreatedBy   sql.NullString `json:"created_by"`
	ModifiedAt  sql.NullTime   `json:"modified_at"`
	ModifiedBy  sql.NullString `json:"modified_by"`
	DeletedAt   sql.NullTime   `json:"deleted_at"`
	DeletedBy   sql.NullString `json:"deleted_by"`
}
