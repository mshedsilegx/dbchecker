package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"os"
	"runtime"

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

	// Example obfuscated key (replace with your actual obfuscated key)
	obfuscatedKey, _ := base64.StdEncoding.DecodeString("your_obfuscated_key_here")
	deobfuscatedKey := crypto.XORDecrypt(obfuscatedKey, secretKeyBytes)

	if *encryptFlag != "" {
		encryptedPassword, err := crypto.Encrypt([]byte(*encryptFlag), deobfuscatedKey)
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

	if *dbID != "" {
		dbConfig, ok := cfg.Databases[*dbID]
		if !ok {
			fmt.Printf("Database with ID '%s' not found in config\n", *dbID)
			os.Exit(1)
		}
		if err := checkDatabase(dbConfig, *dbID, deobfuscatedKey); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else {
		for id, dbConfig := range cfg.Databases {
			if err := checkDatabase(dbConfig, id, deobfuscatedKey); err != nil {
				fmt.Println(err)
			}
		}
	}
}

func checkDatabase(dbConfig config.DatabaseConfig, dbID string, secretKey []byte) error {
	decryptedPassword, err := crypto.Decrypt([]byte(dbConfig.Password), secretKey)
	if err != nil {
		return fmt.Errorf("password decryption failed for %s: %w", dbID, err)
	}

	var db database.DB
	switch dbConfig.Type {
	case "mysql":
		db = &database.MySQL{}
	case "postgres":
		db = &database.Postgres{}
	case "oracle":
		db = &database.Oracle{}
	case "sqlserver":
		db = &database.SQLServer{}
	case "sqlite":
		db = &database.SQLite{}
	case "mongodb":
		db = &database.MongoDB{}
	default:
		return fmt.Errorf("unsupported database type: %s", dbConfig.Type)
	}

	if err := db.Connect(dbConfig, string(decryptedPassword)); err != nil {
		return fmt.Errorf("connection failed for %s: %w", dbID, err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		return fmt.Errorf("ping failed for %s: %w", dbID, err)
	}

	if dbConfig.HealthQuery != "" {
		if err := db.HealthCheck(dbConfig.HealthQuery); err != nil {
			return fmt.Errorf("health check failed for %s: %w", dbID, err)
		}
	}

	fmt.Printf("Successfully connected and checked %s\n", dbID)
	return nil
}
