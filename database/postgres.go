package database

import (
	"database/sql"
	"fmt"
	"net/url"

	"criticalsys.net/dbchecker/config"
	_ "github.com/lib/pq"
)

type Postgres struct {
	SQLBase
}

func (p *Postgres) Connect(cfg config.DatabaseConfig, decryptedPassword string) error {
	dsn := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(cfg.User, decryptedPassword),
		Host:   fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Path:   cfg.Name,
	}

	query := dsn.Query()
	sslMode := "disable" // Default to disable
	switch cfg.TLSMode {
	case "require":
		sslMode = "require"
	case "verify-ca":
		sslMode = "verify-ca"
	case "verify-full":
		sslMode = "verify-full"
	case "disable", "": // Treat empty as disable
		sslMode = "disable"
	default:
		return fmt.Errorf("invalid tls_mode for postgres: %s", cfg.TLSMode)
	}
	query.Set("sslmode", sslMode)
	dsn.RawQuery = query.Encode()

	db, err := sql.Open("postgres", dsn.String())
	if err != nil {
		return err
	}
	p.db = db
	return nil
}

func (p *Postgres) Close() error {
	return p.db.Close()
}
