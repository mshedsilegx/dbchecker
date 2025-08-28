package database

import (
	"database/sql"

	"criticalsys.net/dbchecker/config"
	_ "github.com/mattn/go-sqlite3"
)

type SQLite struct {
	SQLBase
}

func (s *SQLite) Connect(cfg config.DatabaseConfig, decryptedPassword string) error {
	db, err := sql.Open("sqlite3", cfg.Name)
	if err != nil {
		return err
	}
	s.db = db
	return nil
}
