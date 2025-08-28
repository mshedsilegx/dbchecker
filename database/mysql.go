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
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	mysqlConfig := mysql.Config{
		User:                 cfg.User,
		Passwd:               decryptedPassword,
		Net:                  "tcp",
		Addr:                 addr,
		DBName:               cfg.Name,
		AllowNativePasswords: true,
	}

	tlsConfig, err := buildTLSConfig(cfg.TLSMode, cfg.Host)
	if err != nil {
		return err
	}

	if tlsConfig != nil {
		// The mysql driver requires registering the config and using a key.
		// We use the address as a unique key for this connection.
		tlsKey := fmt.Sprintf("dbchecker-tls-%s", addr)
		// RegisterTLSConfig is not thread-safe, but in our concurrent model,
		// each goroutine works on a different db config. If multiple configs
		// point to the same server, this could be an issue, but it's a rare edge case.
		// The check for ErrTLSConfigAlreadyRegistered mitigates races.
		if err := mysql.RegisterTLSConfig(tlsKey, tlsConfig); err != nil && err != mysql.ErrTLSConfigAlreadyRegistered {
			return fmt.Errorf("could not register mysql tls config: %w", err)
		}
		mysqlConfig.TLSConfig = tlsKey
	} else {
		// Explicitly set TLS to false when disabled
		mysqlConfig.TLSConfig = "false"
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
