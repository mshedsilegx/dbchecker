package database

import (
	"database/sql"
	"fmt"

	"criticalsys.net/dbchecker/config"
	_ "github.com/mattn/go-sqlite3"
)

type SQLite struct{}

func (s *SQLite) Connect(cfg config.DatabaseConfig, decryptedPassword string) error {
	db, err := sql.Open("sqlite3", cfg.Name)
	if err != nil {
		return err
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		return err
	}

	if cfg.HealthQuery != "" {
		_, err := db.Exec(cfg.HealthQuery)
		if err != nil {
			return fmt.Errorf("health check query failed: %w", err)
		}
	}

	return nil
}
