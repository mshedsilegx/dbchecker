package database

import (
	"context"
	"crypto/tls"
	"fmt"

	"criticalsys.net/dbchecker/config"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	client   *mongo.Client
	database string
}

func (m *MongoDB) Connect(cfg config.DatabaseConfig, decryptedPassword string) error {
	clientOptions := options.Client()
	clientOptions.SetHosts([]string{fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)})

	if cfg.User != "" {
		creds := options.Credential{
			Username: cfg.User,
			Password: decryptedPassword,
		}
		clientOptions.SetAuth(creds)
	}

	if cfg.TLS {
		clientOptions.SetTLSConfig(&tls.Config{InsecureSkipVerify: true}) // TODO: This is insecure, should be made configurable
	}

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return fmt.Errorf("mongodb connection failed: %w", err)
	}
	m.client = client
	m.database = cfg.Name
	return nil
}

func (m *MongoDB) Ping() error {
	return m.client.Ping(context.TODO(), nil)
}

// HealthCheck for MongoDB lists collection names as a basic check. The query parameter is ignored.
func (m *MongoDB) HealthCheck(query string) error {
	db := m.client.Database(m.database)
	_, err := db.ListCollectionNames(context.TODO(), primitive.M{})
	if err != nil {
		return fmt.Errorf("mongodb list collections failed: %w", err)
	}
	return nil
}

func (m *MongoDB) Close() error {
	return m.client.Disconnect(context.TODO())
}
