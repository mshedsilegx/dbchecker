package database

import (
	"database/sql"
	"fmt"

	"criticalsys.net/dbchecker/config"
	_ "github.com/go-sql-driver/mysql"
)

type MySQL struct{}

func (m *MySQL) Connect(cfg config.DatabaseConfig, decryptedPassword string) error {
	tlsConfig := ""
	if cfg.TLS {
		tlsConfig = "?tls=true"
	}
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s%s", cfg.User, decryptedPassword, cfg.Host, cfg.Port, cfg.Name, tlsConfig)
	db, err := sql.Open("mysql", connectionString)
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
