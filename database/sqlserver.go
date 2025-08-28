package database

import (
	"database/sql"
	"fmt"

	"criticalsys.net/dbchecker/config"
	_ "github.com/microsoft/go-mssqldb"
)

type SQLServer struct {
	SQLBase
}

func (s *SQLServer) Connect(cfg config.DatabaseConfig, decryptedPassword string) error {
	connectionString := fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=%s", cfg.User, decryptedPassword, cfg.Host, cfg.Port, cfg.Name)
	db, err := sql.Open("sqlserver", connectionString)
	if err != nil {
		return err
	}
	s.db = db
	return nil
}
