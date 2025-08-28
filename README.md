# DB Connection Diags

## Overview and Objectives

DB Connection Diags is a command-line utility designed to help developers and system administrators diagnose and verify connections to multiple types of databases. It reads a configuration file, connects to the specified databases, and runs a health check to ensure they are accessible and operational.

The primary objectives of this tool are:
- To provide a simple, configuration-driven way to check database health.
- To support a wide range of popular databases (SQL and NoSQL).
- To provide a secure way to handle database credentials by using encryption.

## Command-Line Arguments

The application accepts the following command-line arguments:

- `-config <path>`: Specifies the path to the YAML configuration file. Defaults to `config.yaml`.
- `-db <id>`: The identifier of a specific database to check (as defined in the config file). If this flag is omitted, the tool will check all databases listed in the configuration.
- `-version`: Displays the application's version information and exits.
- `-encrypt <password>`: Encrypts the given password using the secret key, prints the base64-encoded result to the console, and exits.
- `-key-file <path>`: Path to the secret key file. If provided, it overrides the `DB_SECRET_KEY` environment variable.

## Configuration File

The application uses a YAML file (e.g., `config.yaml`) to configure the database connections.

### Structure

```yaml
databases:
  <database-id>:
    type: "string"
    host: "string"
    port: int
    user: "string"
    password: "string"
    name: "string"
    tls: bool
    health_query: "string"
```

- `password`: The **encrypted** password for the user. See the security section for instructions on how to encrypt passwords.

## Examples

### PostgreSQL Health Check

To check a PostgreSQL database, add an entry like this to your `config.yaml`:

```yaml
databases:
  my_postgres_db:
    type: "postgres"
    host: "localhost"
    port: 5432
    user: "pguser"
    password: "ENCRYPTED_PASSWORD_HERE"
    name: "mydatabase"
    tls: false
    health_query: "SELECT 1"
```

Run the check using a key file (recommended):
```bash
./dbchecker -key-file /path/to/your/secret.key -db my_postgres_db
```

Alternatively, run the check using the environment variable:
```bash
export DB_SECRET_KEY="your-32-byte-secret-key"
./dbchecker -db my_postgres_db
```

## Security and Credential Encryption

To avoid storing plaintext passwords in the configuration file, the application uses AES-GCM encryption to secure your credentials.

### Providing the Secret Key

You can provide the secret key to the application in one of two ways.

1.  **Key File (Recommended)**: You can store the key in a file and provide the path to it using the `-key-file` flag. This is the recommended approach for production environments.
2.  **Environment Variable**: You can provide the key via the `DB_SECRET_KEY` environment variable.

**The `-key-file` flag always takes precedence over the `DB_SECRET_KEY` environment variable if both are present.**

### How it Works

1.  **Secret Key**: The application uses a secret key for encryption and decryption. This key **must be 32 bytes (256 bits) long.**
2.  **Encryption**: You use the `-encrypt` flag to encrypt your database password. The application uses the secret key (from the key file or environment variable) to perform the encryption and gives you a base64-encoded string.
3.  **Configuration**: You paste this encrypted string into the `password` field in your `config.yaml`.
4.  **Decryption**: When the application runs, it reads the encrypted password from the config, gets the secret key, and decrypts the password in memory just before connecting to the database.

As a legacy feature, the secret key itself is passed through a weak XOR de-obfuscation step with a key compiled into the application. For new deployments, it is recommended to manage the secret key securely using a proper secrets management system.

### How to Encrypt Your Credentials

**Step 1: Generate a Secure Secret Key**

You need a secure, 32-byte secret key.

*   **To create a key file (recommended):**

    On Linux or macOS:
    ```bash
    openssl rand 32 > /path/to/your/secret.key
    # Set strict, read-only permissions for the key file
    chmod 400 /path/to/your/secret.key
    ```
    The application will check for these exact permissions on non-Windows systems and will not run if the file is accessible by other users.

    On Windows:
    You should use the file system's ACL (Access Control List) features to restrict access to the current user. You can do this with the `icacls` command:
    ```powershell
    # First, disable inheritance to remove other users
    icacls "C:\path\to\your\secret.key" /inheritance:r
    # Then, grant access only to your user
    icacls "C:\path\to\your\secret.key" /grant:r "$($env:USERNAME):(R)"
    ```

*   **To generate a key for the environment variable:**
    ```bash
    openssl rand -base64 32
    ```

**Step 2: Encrypt Your Password**

Now, use the application's `-encrypt` flag to encrypt your database password.

*   **Using the key file:**
    ```bash
    ./dbchecker -key-file /path/to/your/secret.key -encrypt 'my-super-secret-password'
    ```

*   **Using the environment variable:**
    ```bash
    export DB_SECRET_KEY="your-generated-key"
    ./dbchecker -encrypt 'my-super-secret-password'
    ```

The application will output a long, base64-encoded string. This is your encrypted password.

**Step 3: Update Your Configuration**

Copy the encrypted password from Step 2 and paste it into the `password` field in your `config.yaml` file.

Your application is now configured to securely connect to your databases!
