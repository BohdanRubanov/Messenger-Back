package models

import "time"

type Product struct {
	// name for json and sqlx
	ID          int       `json:"id" db:"id"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
	Price       int       `json:"price" db:"price"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

type CreateProduct struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Price       int    `json:"price"`
}

type UpdateProduct struct {
	// pointer to a string, can be nil if the field is not provided
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Price       *int    `json:"price"`
}
