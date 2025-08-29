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

func (m *MongoDB) Connect(ctx context.Context, cfg config.DatabaseConfig, decryptedPassword string) error {
	clientOptions := options.Client()
	clientOptions.SetHosts([]string{fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)})

	if cfg.User != "" {
		creds := options.Credential{
			Username: cfg.User,
			Password: decryptedPassword,
		}
		clientOptions.SetAuth(creds)
	}

	tlsConfig, err := buildTLSConfig(cfg.TLSMode, cfg.Host, cfg.RootCertPath, cfg.ClientCertPath, cfg.ClientKeyPath)
	if err != nil {
		return err
	}
	if tlsConfig != nil {
		clientOptions.SetTLSConfig(tlsConfig)
	}

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return fmt.Errorf("mongodb connection failed: %w", err)
	}
	m.client = client
	m.database = cfg.Name
	return nil
}

func (m *MongoDB) Ping(ctx context.Context) error {
	return m.client.Ping(ctx, nil)
}

// HealthCheck for MongoDB lists collection names as a basic check. The query parameter is ignored.
func (m *MongoDB) HealthCheck(ctx context.Context, query string) error {
	db := m.client.Database(m.database)
	_, err := db.ListCollectionNames(ctx, primitive.M{})
	if err != nil {
		return fmt.Errorf("mongodb list collections failed: %w", err)
	}
	return nil
}

func (m *MongoDB) Close() error {
	// Close does not need a context, but we might want one for graceful shutdown in the future.
	// For now, we'll use a background context to satisfy the disconnect method.
	return m.client.Disconnect(context.Background())
}
