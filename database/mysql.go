package database

import (
	"database/sql"
	"fmt"

	"criticalsys.net/dbchecker/config"
	_ "github.com/go-sql-driver/mysql"
)

type MySQL struct {
	SQLBase
}

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
	m.db = db
	return nil
}
