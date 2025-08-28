package database

import (
	"database/sql"
	"fmt"
	"net/url"

	"criticalsys.net/dbchecker/config"
	_ "github.com/microsoft/go-mssqldb"
)

type SQLServer struct {
	SQLBase
}

func (s *SQLServer) Connect(cfg config.DatabaseConfig, decryptedPassword string) error {
	query := url.Values{}
	query.Add("database", cfg.Name)

	dsn := &url.URL{
		Scheme:   "sqlserver",
		User:     url.UserPassword(cfg.User, decryptedPassword),
		Host:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		RawQuery: query.Encode(),
	}

	db, err := sql.Open("sqlserver", dsn.String())
	if err != nil {
		return err
	}
	s.db = db
	return nil
}
