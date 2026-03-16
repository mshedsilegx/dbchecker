# Design - DB Connection Diags

## Design Principles

### 1. Extensibility
The application is designed to support new database types with minimal changes to the core logic. 
- To add a new database:
    1. Implement the `DB` interface in a new file within the `database/` package.
    2. Register the new type in the `New()` factory function.
    3. Update the `supportedDBTypes` map in `config/config.go`.

### 2. Security by Default
- **AES-GCM**: Chosen for its performance and built-in integrity checking.
- **Lazy Connection (MongoDB v2)**: The MongoDB implementation follows the v2 driver design where `Connect` initializes the client, but actual connection health is verified via `Ping`.
- **Scoped IO**: All file operations (Config, TLS certificates) are scoped to their respective parent directories using Go 1.24 `os.Root` to mitigate CWE-22 (Path Traversal).

### 3. Concurrency
- When checking multiple databases, the application uses a `sync.WaitGroup` to run checks in parallel. This significantly improves performance when dealing with high-latency network connections or multiple database clusters.

### 4. Error Handling
- Errors are wrapped with context (e.g., `fmt.Errorf("connection failed for %s: %w", ...)`) to ensure that logs provide clear information about which database failed and why, while preserving the underlying error for programmatic inspection if needed.

## Data Flow
1. **Input**: User provides `-config`, `-db` (optional), and a secret key.
2. **Process**:
    - Load and validate YAML.
    - For each target database:
        - Decrypt password using the provided key.
        - Instantiate driver via factory.
        - Call `Connect()`.
        - Call `Ping()` to verify basic connectivity.
        - Call `HealthCheck()` (if query provided) to verify functional status.
3. **Output**: Success message or detailed error per database.
