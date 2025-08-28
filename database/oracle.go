package database

import (
	"database/sql"
	"fmt"

	"criticalsys.net/dbchecker/config"
	_ "github.com/sijms/go-ora/v2"
)

type Oracle struct {
	SQLBase
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

func (o *Oracle) Close() error {
	return o.db.Close()
}
