# Architecture - DB Connection Diags

## Overview
DB Connection Diags is a modular Go application designed for multi-database connectivity verification. It emphasizes security, extensibility, and ease of use through a configuration-driven approach.

## Component Diagram
```
[ Main Entry ]
      |
      |--- [ Config Package ] (YAML Loading & Validation)
      |
      |--- [ Crypto Package ] (AES-GCM Encryption/Decryption)
      |
      |--- [ Database Package ] (Abstract Interface & Drivers)
                |
                |--- MySQL Implementation
                |--- PostgreSQL Implementation
                |--- MongoDB Implementation (v2)
                |--- Oracle Implementation
                |--- SQL Server Implementation
                |--- SQLite Implementation
```

## Core Packages

### 1. Main Package (`main.go`)
- **Responsibility**: Application lifecycle management.
- **Workflow**:
    - Parses command-line flags.
    - Retrieves the secret key from a file or environment variable.
    - Handles password encryption requests (`-encrypt`).
    - Orchestrates database checks concurrently using goroutines if multiple IDs are specified.

### 2. Config Package (`config/`)
- **Responsibility**: Typed configuration management.
- **Key Features**:
    - YAML unmarshaling into structured Go types.
    - Strict validation of supported database types and TLS modes.
    - **Security**: Uses `os.OpenRoot` (Go 1.24+) to prevent directory traversal when loading configuration files.

### 3. Crypto Package (`crypto/`)
- **Responsibility**: Secure credential handling.
- **Mechanism**:
    - Uses **AES-256-GCM** for authenticated encryption.
    - Every encryption operation generates a unique nonce, which is prepended to the ciphertext.
    - Decryption extracts the nonce and verifies the integrity of the ciphertext before returning the plaintext.

### 4. Database Package (`database/`)
- **Responsibility**: Unified abstraction for heterogeneous database engines.
- **Architecture**:
    - **Interface-based**: All drivers implement the `DB` interface (`Connect`, `Ping`, `HealthCheck`, `Close`).
    - **Factory Pattern**: The `New()` function returns the appropriate implementation based on the `type` string.
    - **TLS Support**: Centralized TLS configuration logic in `tls.go` using `os.OpenRoot` for secure certificate loading.
    - **v2 Driver**: The MongoDB implementation specifically utilizes `go.mongodb.org/mongo-driver/v2`.

## Security Model
- **Credential Protection**: Passwords are never stored in plaintext. They are decrypted in memory only at the moment of connection.
- **File System Security**: 
    - Key files require strict permissions (0400 on Linux/macOS).
    - Scoped file access using `os.Root` prevents unauthorized directory traversal.
- **TLS/SSL**: Supports various levels of certificate verification and mutual TLS (mTLS).
