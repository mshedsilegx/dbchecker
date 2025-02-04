package main

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/tls"
	"database/sql"
	"encoding/base64"
	"flag"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	_ "github.com/microsoft/go-mssqldb"
	_ "github.com/sijms/go-ora/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/yaml.v3" // Use v3 for better YAML support
)

var version string

type DatabaseConfig struct {
	Type        string `yaml:"type"`
	Host        string `yaml:"host"`
	Port        int    `yaml:"port"`
	User        string `yaml:"user"`
	Password    string `yaml:"password"`
	Name        string `yaml:"name"`
	HealthQuery string `yaml:"health_query"`
	TLS         bool   `yaml:"tls"`
}

type Config struct {
	Databases map[string]DatabaseConfig `yaml:"databases"`
}

func decrypt(data, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}

func xorDecrypt(data []byte, key []byte) []byte {
	decrypted := make([]byte, len(data))
	for i := 0; i < len(data); i++ {
		decrypted[i] = data[i] ^ key[i%len(key)]
	}
	return decrypted
}

func connectAndCheck(config DatabaseConfig, dbID string) error {
	var db *sql.DB
	var err error

	decryptedPassword, err := decrypt([]byte(config.Password), []byte(os.Getenv("DB_SECRET_KEY")))
	if err != nil {
		return fmt.Errorf("password decryption failed for %s: %w", dbID, err)
	}

	connectionString := ""

	switch config.Type {
	case "mysql":
		tlsConfig := ""
		if config.TLS {
			tlsConfig = "?tls=true"
		}
		connectionString = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s%s", config.User, string(decryptedPassword), config.Host, config.Port, config.Name, tlsConfig)
		db, err = sql.Open("mysql", connectionString)
	case "postgres":
		tlsConfig := ""
		if config.TLS {
			tlsConfig = " sslmode=require"
		}
		connectionString = fmt.Sprintf("postgres://%s:%s@%s:%d/%s%s", config.User, string(decryptedPassword), config.Host, config.Port, config.Name, tlsConfig)
		db, err = sql.Open("postgres", connectionString)
	case "oracle":
		connectionString = fmt.Sprintf("%s/%s@%s:%d/%s", config.User, string(decryptedPassword), config.Host, config.Port, config.Name) // Simplified Oracle connection
		db, err = sql.Open("oracle", connectionString)
	case "sqlserver":
		connectionString = fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=%s", config.User, string(decryptedPassword), config.Host, config.Port, config.Name)
		db, err = sql.Open("sqlserver", connectionString)
	case "sqlite":
		connectionString = config.Name // Database file path
		db, err = sql.Open("sqlite3", connectionString)
	case "mongodb":
		clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s:%d/%s", config.User, string(decryptedPassword), config.Host, config.Port, config.Name))
		if config.TLS {
			clientOptions = clientOptions.SetTLSConfig(&tls.Config{InsecureSkipVerify: true}) // Adjust TLS config as needed
		}
		client, err := mongo.Connect(context.TODO(), clientOptions)
		if err != nil {
			return fmt.Errorf("mongodb connection failed for %s: %w", dbID, err)
		}
		defer client.Disconnect(context.TODO())

		err = client.Ping(context.TODO(), nil)
		if err != nil {
			return fmt.Errorf("mongodb ping failed for %s: %w", dbID, err)
		}

		db := client.Database(config.Name)

		names, err := db.ListCollectionNames(context.TODO(), primitive.M{})
		if err != nil {
			return fmt.Errorf("mongodb list collections failed for %s: %w", dbID, err)
		}
		fmt.Println("Collections found:", names)

		return nil // MongoDB check successful
	default:
		return fmt.Errorf("unsupported database type: %s", config.Type)
	}

	if err != nil {
		return fmt.Errorf("connection failed for %s: %w", dbID, err)
	}
	defer db.Close()

	if config.Type != "mongodb" { // Skip health check for MongoDB (handled above)
		if _, err := db.Exec(config.HealthQuery); err != nil {
			return fmt.Errorf("health check query failed for %s: %w", dbID, err)
		}
	}

	fmt.Printf("Successfully connected and checked %s\n", dbID)
	return nil
}

func main() {
	configFile := flag.String("config", "config.yaml", "Path to the configuration file")
	dbID := flag.String("db", "", "Identifier of the database to check")
	versionFlag := flag.Bool("version", false, "Display version information")
	flag.Parse()

	// Version
	if *versionFlag {
		fmt.Printf("DB Connection Diags - Version: %s\n", version)
		os.Exit(0)
	}

	if os.Getenv("DB_SECRET_KEY") == "" {
		fmt.Println("DB_SECRET_KEY environment variable is not set.")
		os.Exit(1)
	}

	// XOR decrypt the secret key
	key := []byte(os.Getenv("DB_SECRET_KEY"))

	// Example obfuscated key (replace with your actual obfuscated key)
	obfuscatedKey, _ := base64.StdEncoding.DecodeString("your_obfuscated_key_here")

	deobfuscatedKey := xorDecrypt(obfuscatedKey, key)
	os.Setenv("DB_SECRET_KEY", string(deobfuscatedKey))

	data, err := os.ReadFile(*configFile)
	if err != nil {
		fmt.Printf("Error reading config file: %v\n", err)
		os.Exit(1)
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		fmt.Printf("Error unmarshalling config file: %v\n", err)
		os.Exit(1)
	}

	if *dbID != "" {
		dbConfig, ok := config.Databases[*dbID]
		if !ok {
			fmt.Printf("Database with ID '%s' not found in config\n", *dbID)
			os.Exit(1)
		}
		if err := connectAndCheck(dbConfig, *dbID); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else {
		for id, dbConfig := range config.Databases {
			if err := connectAndCheck(dbConfig, id); err != nil {
				fmt.Println(err)
				// Decide whether to continue on error or exit.  Currently continues.
				// os.Exit(1)  // Uncomment to stop on first error.
			}
		}
	}
}
