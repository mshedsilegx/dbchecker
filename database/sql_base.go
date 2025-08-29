package database

import (
	"context"
	"database/sql"
	"fmt"
)

// SQLBase is a helper struct that can be embedded by SQL-based database
// providers to share common functionality.
type SQLBase struct {
	db *sql.DB
}

// Ping sends a ping to the database to verify the connection is alive.
func (b *SQLBase) Ping(ctx context.Context) error {
	return b.db.PingContext(ctx)
}

// HealthCheck executes a simple query to verify the database is operational.
func (b *SQLBase) HealthCheck(ctx context.Context, query string) error {
	_, err := b.db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("health check query failed: %w", err)
	}
	return nil
}

// Close closes the database connection.
func (b *SQLBase) Close() error {
	return b.db.Close()
}
