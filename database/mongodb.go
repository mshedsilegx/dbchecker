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

type MongoDB struct{}

func (m *MongoDB) Connect(cfg config.DatabaseConfig, decryptedPassword string) error {
	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s:%d/%s", cfg.User, decryptedPassword, cfg.Host, cfg.Port, cfg.Name))
	if cfg.TLS {
		clientOptions = clientOptions.SetTLSConfig(&tls.Config{InsecureSkipVerify: true}) // Adjust TLS config as needed
	}
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return fmt.Errorf("mongodb connection failed: %w", err)
	}
	defer client.Disconnect(context.TODO())

	if err := client.Ping(context.TODO(), nil); err != nil {
		return fmt.Errorf("mongodb ping failed: %w", err)
	}

	db := client.Database(cfg.Name)
	names, err := db.ListCollectionNames(context.TODO(), primitive.M{})
	if err != nil {
		return fmt.Errorf("mongodb list collections failed: %w", err)
	}
	fmt.Println("Collections found:", names)

	return nil
}
