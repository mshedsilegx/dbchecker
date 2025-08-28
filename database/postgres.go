package database

import (
	"database/sql"
	"fmt"

	"criticalsys.net/dbchecker/config"
	_ "github.com/lib/pq"
)

type Postgres struct {
	db *sql.DB
}

func (p *Postgres) Connect(cfg config.DatabaseConfig, decryptedPassword string) error {
	tlsConfig := "sslmode=disable"
	if cfg.TLS {
		tlsConfig = "sslmode=require"
	}
	connectionString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?%s", cfg.User, decryptedPassword, cfg.Host, cfg.Port, cfg.Name, tlsConfig)
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return err
	}
	p.db = db
	return nil
}

func (p *Postgres) Ping() error {
	return p.db.Ping()
}

func (p *Postgres) HealthCheck(query string) error {
	_, err := p.db.Exec(query)
	if err != nil {
		return fmt.Errorf("health check query failed: %w", err)
	}
	return nil
}

func (p *Postgres) Close() error {
	return p.db.Close()
}
