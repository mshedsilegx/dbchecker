package database

import (
	"database/sql"
	"fmt"

	"criticalsys.net/dbchecker/config"
	_ "github.com/mattn/go-sqlite3"
)

type SQLite struct {
	db *sql.DB
}

func (s *SQLite) Connect(cfg config.DatabaseConfig, decryptedPassword string) error {
	db, err := sql.Open("sqlite3", cfg.Name)
	if err != nil {
		return err
	}
	s.db = db
	return nil
}

func (s *SQLite) Ping() error {
	return s.db.Ping()
}

func (s *SQLite) HealthCheck(query string) error {
	_, err := s.db.Exec(query)
	if err != nil {
		return fmt.Errorf("health check query failed: %w", err)
	}
	return nil
}

func (s *SQLite) Close() error {
	return s.db.Close()
}
