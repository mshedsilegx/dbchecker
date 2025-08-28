package database

import "criticalsys.net/dbchecker/config"

type DB interface {
	Connect(cfg config.DatabaseConfig, decryptedPassword string) error
	Ping() error
	HealthCheck(query string) error
	Close() error
}
