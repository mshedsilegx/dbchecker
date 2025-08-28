package database

import (
	"database/sql"
	"fmt"

	"criticalsys.net/dbchecker/config"
	_ "github.com/sijms/go-ora/v2"
)

type Oracle struct {
	db *sql.DB
}

func (o *Oracle) Connect(cfg config.DatabaseConfig, decryptedPassword string) error {
	connectionString := fmt.Sprintf("%s/%s@%s:%d/%s", cfg.User, decryptedPassword, cfg.Host, cfg.Port, cfg.Name)
	db, err := sql.Open("oracle", connectionString)
	if err != nil {
		return err
	}
	o.db = db
	return nil
}

func (o *Oracle) Ping() error {
	return o.db.Ping()
}

func (o *Oracle) HealthCheck(query string) error {
	_, err := o.db.Exec(query)
	if err != nil {
		return fmt.Errorf("health check query failed: %w", err)
	}
	return nil
}

func (o *Oracle) Close() error {
	return o.db.Close()
}
