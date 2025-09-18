package main

import (
	"context"
	"encoding/base64"
	"flag"
	"fmt"
	"os"
	"runtime"
	"sync"

	"criticalsys.net/dbchecker/config"
	"criticalsys.net/dbchecker/crypto"
	"criticalsys.net/dbchecker/database"
)

var version string

func main() {
	configFile := flag.String("config", "config.yaml", "Path to the configuration file")
	dbID := flag.String("db", "", "Identifier of the database to check")
	versionFlag := flag.Bool("version", false, "Display version information")
	encryptFlag := flag.String("encrypt", "", "Encrypt a password and exit")
	keyFileFlag := flag.String("key-file", "", "Path to the secret key file (overrides DB_SECRET_KEY)")
	flag.Parse()

	if *versionFlag {
		fmt.Printf("DB Connection Diags - Version: %s\n", version)
		os.Exit(0)
	}

	var secretKeyBytes []byte
	if *keyFileFlag != "" {
		// Use key file
		fileInfo, err := os.Stat(*keyFileFlag)
		if err != nil {
			fmt.Printf("Error accessing key file: %v\n", err)
			os.Exit(1)
		}
		if runtime.GOOS != "windows" && fileInfo.Mode().Perm() != 0400 {
			fmt.Printf("Error: Key file permissions must be 400 (read-only for owner).")
			os.Exit(1)
		}
		secretKeyBytes, err = os.ReadFile(*keyFileFlag)
		if err != nil {
			fmt.Printf("Error reading key file: %v\n", err)
			os.Exit(1)
		}
	} else {
		// Fallback to environment variable
		secretKey := os.Getenv("DB_SECRET_KEY")
		if secretKey == "" {
			fmt.Println("Error: No secret key provided. Use -key-file flag or set DB_SECRET_KEY environment variable.")
			os.Exit(1)
		}
		secretKeyBytes = []byte(secretKey)
	}

	// The user-provided secret key is used directly.
	// The unnecessary XOR obfuscation layer has been removed.
	if *encryptFlag != "" {
		encryptedPassword, err := crypto.Encrypt([]byte(*encryptFlag), secretKeyBytes)
		if err != nil {
			fmt.Printf("Error encrypting password: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(base64.StdEncoding.EncodeToString(encryptedPassword))
		os.Exit(0)
	}

	cfg, err := config.LoadConfig(*configFile)
	if err != nil {
		fmt.Printf("Error loading config file: %v\n", err)
		os.Exit(1)
	}

	ctx := context.Background()

	if *dbID != "" {
		dbConfig, ok := cfg.Databases[*dbID]
		if !ok {
			fmt.Printf("Database with ID '%s' not found in config\n", *dbID)
			os.Exit(1)
		}
		if err := checkDatabase(ctx, dbConfig, *dbID, secretKeyBytes); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else {
		var wg sync.WaitGroup
		for id, dbConfig := range cfg.Databases {
			wg.Add(1)
			go func(id string, dbConfig config.DatabaseConfig) {
				defer wg.Done()
				if err := checkDatabase(ctx, dbConfig, id, secretKeyBytes); err != nil {
					fmt.Println(err)
				}
			}(id, dbConfig)
		}
		wg.Wait()
	}
}

func checkDatabase(ctx context.Context, dbConfig config.DatabaseConfig, dbID string, secretKey []byte) error {
	decryptedPassword, err := crypto.Decrypt([]byte(dbConfig.Password), secretKey)
	if err != nil {
		return fmt.Errorf("password decryption failed for %s: %w", dbID, err)
	}

	db, err := database.New(dbConfig.Type)
	if err != nil {
		return err
	}

	if err := db.Connect(ctx, dbConfig, string(decryptedPassword)); err != nil {
		return fmt.Errorf("connection failed for %s: %w", dbID, err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			fmt.Printf("Error closing database connection for %s: %v\n", dbID, err)
		}
	}()

	if err := db.Ping(ctx); err != nil {
		return fmt.Errorf("ping failed for %s: %w", dbID, err)
	}

	if dbConfig.HealthQuery != "" {
		if err := db.HealthCheck(ctx, dbConfig.HealthQuery); err != nil {
			return fmt.Errorf("health check failed for %s: %w", dbID, err)
		}
	}

	fmt.Printf("Successfully connected and checked %s\n", dbID)
	return nil
}
