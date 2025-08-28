package database

import (
	"criticalsys.net/dbchecker/config"
	"fmt"
)

type DB interface {
	Connect(cfg config.DatabaseConfig, decryptedPassword string) error
	Ping() error
	HealthCheck(query string) error
	Close() error
}

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
