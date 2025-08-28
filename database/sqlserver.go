package database

import (
	"database/sql"
	"fmt"

	"criticalsys.net/dbchecker/config"
	_ "github.com/microsoft/go-mssqldb"
)

type SQLServer struct {
	db *sql.DB
}

func (s *SQLServer) Connect(cfg config.DatabaseConfig, decryptedPassword string) error {
	connectionString := fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=%s", cfg.User, decryptedPassword, cfg.Host, cfg.Port, cfg.Name)
	db, err := sql.Open("sqlserver", connectionString)
	if err != nil {
		return err
	}
	s.db = db
	return nil
}

func (s *SQLServer) Ping() error {
	return s.db.Ping()
}

func (s *SQLServer) HealthCheck(query string) error {
	_, err := s.db.Exec(query)
	if err != nil {
		return fmt.Errorf("health check query failed: %w", err)
	}
	return nil
}

func (s *SQLServer) Close() error {
	return s.db.Close()
}
