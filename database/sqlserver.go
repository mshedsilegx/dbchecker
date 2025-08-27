package database

import (
	"database/sql"
	"fmt"

	"criticalsys.net/dbchecker/config"
	_ "github.com/microsoft/go-mssqldb"
)

type SQLServer struct{}

func (s *SQLServer) Connect(cfg config.DatabaseConfig, decryptedPassword string) error {
	connectionString := fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=%s", cfg.User, decryptedPassword, cfg.Host, cfg.Port, cfg.Name)
	db, err := sql.Open("sqlserver", connectionString)
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
