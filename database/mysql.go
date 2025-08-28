package database

import (
	"database/sql"
	"fmt"

	"criticalsys.net/dbchecker/config"
	"github.com/go-sql-driver/mysql"
)

type MySQL struct {
	SQLBase
}

func (m *MySQL) Connect(cfg config.DatabaseConfig, decryptedPassword string) error {
	mysqlConfig := mysql.Config{
		User:                 cfg.User,
		Passwd:               decryptedPassword,
		Net:                  "tcp",
		Addr:                 fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		DBName:               cfg.Name,
		AllowNativePasswords: true,
	}
	if cfg.TLS {
		mysqlConfig.TLSConfig = "true"
	}

	connectionString := mysqlConfig.FormatDSN()
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		return err
	}
	m.db = db
	return nil
}

func (m *MySQL) Close() error {
	return m.db.Close()
}
