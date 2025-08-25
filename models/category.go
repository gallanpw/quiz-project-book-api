package models

import "database/sql"

type Category struct {
	ID         int            `json:"id"`
	Name       string         `json:"name" binding:"required"`
	CreatedAt  sql.NullTime   `json:"created_at"`
	CreatedBy  sql.NullString `json:"created_by"`
	ModifiedAt sql.NullTime   `json:"modified_at"`
	ModifiedBy sql.NullString `json:"modified_by"`
	DeletedAt  sql.NullTime   `json:"deleted_at"`
	DeletedBy  sql.NullString `json:"deleted_by"`
}
