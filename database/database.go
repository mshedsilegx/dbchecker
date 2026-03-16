/*
Package database provides a common interface and specific implementations for various database drivers.
It abstracts connection management, pinging, and health check execution across different SQL and NoSQL engines.
*/
package database

import (
	"context"
	"criticalsys.net/dbchecker/config"
	"fmt"
)

// DB defines the standard operations required for any supported database type.
type DB interface {
	Connect(ctx context.Context, cfg config.DatabaseConfig, decryptedPassword string) error
	Ping(ctx context.Context) error
	HealthCheck(ctx context.Context, query string) error
	Close() error
}

// New is a factory function that returns an implementation of the DB interface
// based on the provided database type string.
func New(dbType string) (DB, error) {
	switch dbType {
	case "mysql":
		return &MySQL{}, nil
	case "postgres":
		return &Postgres{}, nil
	case "oracle":
		return &Oracle{}, nil
	case "sqlserver":
		return &SQLServer{}, nil
	case "sqlite":
		return &SQLite{}, nil
	case "mongodb":
		return &MongoDB{}, nil
	default:
		return nil, fmt.Errorf("unsupported database type: %s", dbType)
	}
}
