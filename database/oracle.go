package database

import (
	"database/sql"
	"fmt"

	"criticalsys.net/dbchecker/config"
	go_ora "github.com/sijms/go-ora/v2"
)

type Oracle struct {
	SQLBase
}

func (o *Oracle) Connect(cfg config.DatabaseConfig, decryptedPassword string) error {
	urlOptions := make(map[string]string)

	switch cfg.TLSMode {
	case "disable", "":
		// No options needed for non-TLS connection
	case "require":
		urlOptions["ssl"] = "true"
		urlOptions["ssl verify"] = "false"
	case "verify-ca", "verify-full":
		// The go-ora driver requires a wallet for full certificate verification.
		// Since wallet_path is not a configuration option, we cannot support this mode.
		return fmt.Errorf("tls_mode '%s' is not supported for Oracle without a configured wallet path", cfg.TLSMode)
	default:
		return fmt.Errorf("invalid tls_mode for oracle: %s", cfg.TLSMode)
	}

	connectionString := go_ora.BuildUrl(cfg.Host, cfg.Port, cfg.Name, cfg.User, decryptedPassword, urlOptions)
	db, err := sql.Open("oracle", connectionString)
	if err != nil {
		return err
	}
	o.db = db
	return nil
}

func (o *Oracle) Close() error {
	return o.db.Close()
}
