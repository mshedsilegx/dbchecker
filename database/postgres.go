package database

import (
	"database/sql"
	"fmt"

	"criticalsys.net/dbchecker/config"
	_ "github.com/lib/pq"
)

type Postgres struct{}

func (p *Postgres) Connect(cfg config.DatabaseConfig, decryptedPassword string) error {
	tlsConfig := ""
	if cfg.TLS {
		tlsConfig = " sslmode=require"
	}
	connectionString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s%s", cfg.User, decryptedPassword, cfg.Host, cfg.Port, cfg.Name, tlsConfig)
	db, err := sql.Open("postgres", connectionString)
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
