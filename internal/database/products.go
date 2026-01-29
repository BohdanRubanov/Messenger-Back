package database

import (
	"context"
	"errors"
	"fmt"
	"lesson-proj/internal/models"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ProductRepository — structure (class-like),
// responsible for working with products in the database
type ProductRepository struct {
	// db — pointer to the database connection (sqlx.DB).
	// A pointer is stored to:
	// 1 avoid copying a heavy object
	// 2 use a single connection pool across the entire application
	db *pgxpool.Pool
}

// newProductRepository — factory function (constructor).
// It creates a new ProductRepository object.
//
// Accepts:
// db *pgxpool.Pool — an already initialized database connection
//
// Returns:
// - *ProductRepository — a pointer to the newly created repository
func NewProductRepository(db *pgxpool.Pool) *ProductRepository {

	// Create a ProductRepository struct,
	// put the database reference inside it,
	// and return a pointer to it
	return &ProductRepository{
		db: db,
	}
}

func (productRepository *ProductRepository) GetAllProducts(ctx context.Context) ([]models.Product, error) {
	// Declare a slice to store products fetched from the database.
	// At this point it is nil and has length 0.
	var products []models.Product

	// SQL query to select all products.
	// Backticks are used to allow a multi-line string.
	query := `
		SELECT id, title, description, price, created_at 
		FROM products 
		ORDER BY created_at;`

	// rows is products from db
	// Query - for multiple rows
	rows, err := productRepository.db.Query(ctx, query)

	if err != nil {
		return nil, err
	}
	// defer will close rows when function ends
	defer rows.Close()

	for rows.Next() {
		var product models.Product
		// Scan the current row into the product struct fields.
		err := rows.Scan(
			&product.ID,
			&product.Title,
			&product.Description,
			&product.Price,
			&product.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		products = append(products, product)
	}
	// Check for any errors encountered during iteration.
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// If everything is successful, return the filled slice
	// and nil error.
	return products, nil
}

func (productRepository *ProductRepository) GetProductByID(ctx context.Context, id int) (*models.Product, error) {
	// Declare a variable to store the product.
	// This struct will be filled with data from the database.
	var product models.Product

	// SQL query to fetch a single product by its ID.
	// $1 is a positional placeholder for the id parameter (PostgreSQL syntax).
	// Using placeholders prevents SQL injection.
	query := `
		SELECT id, title, description, price, created_at
		FROM products
		WHERE id = $1;
	`

	// QueryRow is used for a single row result.
	// Scan fills the product struct fields with data from the database.
	// The id variable is passed as the parameter for $1.
	err := productRepository.db.QueryRow(ctx, query, id).Scan(
		&product.ID,
		&product.Title,
		&product.Description,
		&product.Price,
		&product.CreatedAt,
	)
	// Check if the error indicates that no rows were found.
	// erros.Is checks if the error is of type pgx.ErrNoRows
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("product with id %d not found", id)
	}

	if err != nil {
		return nil, err
	}

	// Return a pointer to the filled product struct
	// and nil error to indicate success.
	return &product, nil
}

func (productRepository *ProductRepository) CreateProduct(ctx context.Context, inputProduct models.CreateProduct) (*models.Product, error) {
	var product models.Product

	query := `
		INSERT INTO products (title, description, price, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id, title, description, price, created_at;`
	timeNow := time.Now()
	err := productRepository.db.QueryRow(ctx, query,
		inputProduct.Title,
		inputProduct.Description,
		inputProduct.Price,
		timeNow,
	).Scan(
		&product.ID,
		&product.Title,
		&product.Description,
		&product.Price,
		&product.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (productRepository *ProductRepository) UpdateProduct(ctx context.Context, id int, inputProduct models.UpdateProduct) (*models.Product, error) {
	query := `
		UPDATE products
		SET
			title = COALESCE($1, title),
			description  = COALESCE($2, description),
			price = COALESCE($3, price)
		WHERE id = $4
		RETURNING id, title, description, price, created_at;
	`
	var updatedProduct models.Product
	err := productRepository.db.QueryRow(ctx, query,
		inputProduct.Title,
		inputProduct.Description,
		inputProduct.Price,
		id,
	).Scan(
		&updatedProduct.ID,
		&updatedProduct.Title,
		&updatedProduct.Description,
		&updatedProduct.Price,
		&updatedProduct.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &updatedProduct, nil
}

func (productRepository *ProductRepository) DeleteProduct(ctx context.Context, id int) error {
	query := `
		DELETE FROM products
		WHERE id = $1;`
	// Exec is used for queries that do not return rows
	result, err := productRepository.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	//how many rows were affected
	rows := result.RowsAffected()
	// if product with given id not found
	if rows == 0 {
		return fmt.Errorf("product with id %d not found", id)
	}
	return nil
}
