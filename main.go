package main

import (
        "context"
        "crypto/aes"
        "crypto/cipher"
        "crypto/rand"
        "database/sql"
        "encoding/base64"
        "fmt"
        "log"
        "os"
        "time"

        "github.com/go-yaml/yaml"
        _ "github.com/go-sql-driver/mysql"
        _ "github.com/lib/pq"
        _ "github.com/mattn/go-sqlite3"
        "github.com/microsoft/go-mssqldb"
        "github.com/sijms/go-ora/v2"
        "go.mongodb.org/mongo-driver/mongo"
        "go.mongodb.org/mongo-driver/mongo/options"
        "go.mongodb.org/mongo-driver/bson/primitive"
)

type DBConfig struct {
        Type     string `yaml:"type"`
        Hostname string `yaml:"hostname"`
        Port     int    `yaml:"port"`
        Username string `yaml:"username"`
        Password string `yaml:"password"`
        Database string `yaml:"database"`
        HealthQuery string `yaml:"health_query"`
        TLS      bool   `yaml:"tls"`
}

func main() {
        // Load configuration from YAML file
        configFile, err := os.ReadFile("dbconfig.yaml")
        if err != nil {
                log.Fatal("Error reading config file:", err)
        }

        var configs map[string]DBConfig
        err = yaml.Unmarshal(configFile, &configs)
        if err != nil {
                log.Fatal("Error unmarshalling config:", err)
        }

        // Get and decrypt the secret key
        secretKeyObfuscated := os.Getenv("DB_SECRET_KEY")
    if secretKeyObfuscated == "" {
        log.Fatal("DB_SECRET_KEY environment variable not set")
    }

    secretKey, err := decryptSecretKey(secretKeyObfuscated)
    if err != nil {
        log.Fatal("Error decrypting secret key:", err)
    }

    // Use the deobfuscated key for decryption later on
    // ...

        for name, config := range configs {
                fmt.Printf("Checking connection for %s...\n", name)

                // Decrypt password if needed.
                if config.Password != "" {
                    config.Password, err = decrypt(config.Password, secretKey)
                    if err != nil {
                        log.Fatalf("Error decrypting password for %s: %v", name, err)
                    }
                }

                err = checkConnection(config)
                if err != nil {
                        fmt.Printf("Connection check for %s failed: %v\n", name, err)
                } else {
                        fmt.Printf("Connection check for %s successful.\n", name)
                }
        }
}

func decryptSecretKey(obfuscated string) ([]byte, error) {
    // XOR deobfuscation (example, adapt as needed)
    key := []byte("your_xor_key") // Replace with your XOR key
    decrypted := make([]byte, len(obfuscated))
    for i := 0; i < len(obfuscated); i++ {
        decrypted[i] = obfuscated[i] ^ key[i%len(key)]
    }
    return decrypted, nil
}

func decrypt(encrypted string, key []byte) (string, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return "", err
    }
    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return "", err
    }
    nonceSize := gcm.NonceSize()
    nonce, ciphertext := encrypted[:nonceSize], encrypted[nonceSize:]
    plaintext, err := gcm.Open(nil, []byte(nonce), []byte(ciphertext), nil)
    if err != nil {
        return "", err
    }
    return string(plaintext), nil
}



func checkConnection(config DBConfig) error {
        var db *sql.DB
        var err error

        switch config.Type {
        case "mysql":
                        var tlsConfig string
                        if config.TLS {
                                tlsConfig = "&tls=true"
                        }
                        connectionString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true%s", config.Username, config.Password, config.Hostname, config.Port, config.Database, tlsConfig)
                        db, err = sql.Open("mysql", connectionString)
        case "postgres":
                        var tlsConfig string
                        if config.TLS {
                                tlsConfig = "sslmode=require" // Or verify-full if you have certs
                        }
                        connectionString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?%s", config.Username, config.Password, config.Hostname, config.Port, config.Database, tlsConfig)
                        db, err = sql.Open("postgres", connectionString)
        case "oracle":
                        var tlsConfig string
            if config.TLS {
                tlsConfig = "?connection_string=TRUE&server_cert_mode=1"
            }
                        connectionString := fmt.Sprintf("%s/%s@%s:%d/%s%s", config.Username, config.Password, config.Hostname, config.Port, config.Database, tlsConfig)
                        db, err = sql.Open("oracle", connectionString)
        case "sqlserver":
                        var tlsConfig string
            if config.TLS {
                tlsConfig = "encrypt=true"
            }
                        connectionString := fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=%s&%s", config.Username, config.Password, config.Hostname, config.Port, config.Database, tlsConfig)
                        db, err = sql.Open("sqlserver", connectionString)
        case "sqlite":
                        connectionString := config.Database // SQLite database file path
                        db, err = sql.Open("sqlite3", connectionString)
        case "mongodb":
                        var tlsConfig string
                        if config.TLS {
                                tlsConfig = "&tls=true"
                        }
                        connectionURI := fmt.Sprintf("mongodb://%s:%s@%s:%d/%s%s", config.Username, config.Password, config.Hostname, config.Port, config.Database, tlsConfig)
                        clientOptions := options.Client().ApplyURI(connectionURI)

                        client, err := mongo.Connect(context.Background(), clientOptions)
                        if err != nil {
                                return err
                        }

                        err = client.Ping(context.Background(), nil)
                        if err != nil {
                                return err
                        }

                        // List collections (mongodb healthcheck)
                        names, err := client.Database(config.Database).ListCollectionNames(context.Background(), primitive.M{})
                        if err != nil {
                                return err
                        }

                        fmt.Println("Collections in", config.Database, ":", names)
            return nil // MongoDB connection and healthcheck successful
        default:
                return fmt.Errorf("unsupported database type: %s", config.Type)
        }

        if err != nil {
                return err
        }
        defer db.Close()


        if config.Type != "mongodb" && config.HealthQuery != "" { // Skip health query for MongoDB
                _, err = db.Exec(config.HealthQuery)
                if err != nil {
                        return fmt.Errorf("health query failed: %w", err)
                }
        }

        return nil
}
