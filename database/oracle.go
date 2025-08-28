package database

import (
	"database/sql"
	"fmt"

	"criticalsys.net/dbchecker/config"
	_ "github.com/sijms/go-ora/v2"
)

type Oracle struct{}

func (o *Oracle) Connect(cfg config.DatabaseConfig, decryptedPassword string) error {
	connectionString := fmt.Sprintf("%s/%s@%s:%d/%s", cfg.User, decryptedPassword, cfg.Host, cfg.Port, cfg.Name)
	db, err := sql.Open("oracle", connectionString)
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
