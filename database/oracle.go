package database

import (
	"database/sql"
	"fmt"
	"net/url"

	"criticalsys.net/dbchecker/config"
	_ "github.com/sijms/go-ora/v2"
)

type Oracle struct {
	SQLBase
}

func (o *Oracle) Connect(cfg config.DatabaseConfig, decryptedPassword string) error {
	dsn := url.URL{
		Scheme: "oracle",
		User:   url.UserPassword(cfg.User, decryptedPassword),
		Host:   fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Path:   cfg.Name,
	}

	db, err := sql.Open("oracle", dsn.String())
	if err != nil {
		return err
	}
	o.db = db
	return nil
}

func (o *Oracle) Close() error {
	return o.db.Close()
}
