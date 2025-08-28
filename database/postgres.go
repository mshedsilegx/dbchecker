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
	if cfg.TLS {
		query.Set("sslmode", "require")
	} else {
		query.Set("sslmode", "disable")
	}
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
