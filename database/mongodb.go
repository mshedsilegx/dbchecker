package database

import (
	"context"
	"fmt"

	"criticalsys.net/dbchecker/config"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// MongoDB implements the DB interface for MongoDB using the v2 driver.
type MongoDB struct {
	client   *mongo.Client
	database string
}

// Connect establishes a connection to MongoDB using the provided configuration.
// It uses the v2 driver pattern where mongo.Connect takes only options.
func (m *MongoDB) Connect(ctx context.Context, cfg config.DatabaseConfig, decryptedPassword string) error {
	uri := fmt.Sprintf("mongodb://%s:%d", cfg.Host, cfg.Port)
	clientOptions := options.Client().ApplyURI(uri)

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

	client, err := mongo.Connect(clientOptions)
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
	_, err := db.ListCollectionNames(ctx, bson.D{})
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
