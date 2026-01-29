package database

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Connect initializes and returns a PostgreSQL connection pool
func Connect(databaseURL string) (*pgxpool.Pool, error) {
	ctx := context.Background()

	// Parse pool configuration from the connection string
	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, err
	}

	// Set pool size limits
	config.MinConns = 5  // minimum number of connections kept alive
	config.MaxConns = 20 // maximum number of connections in the pool

	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = 30 * time.Minute

	// Create the pool with the custom config
	dbpool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	return dbpool, nil
}
