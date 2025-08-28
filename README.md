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

- `<database-id>`: A unique name for your database connection (e.g., `my_production_db`).
- `type`: The type of the database. Supported types are: `mysql`, `postgres`, `oracle`, `sqlserver`, `sqlite`, `mongodb`.
- `host`: The hostname or IP address of the database server.
- `port`: The port number for the connection.
- `user`: The username for authentication.
- `password`: The **encrypted** password for the user. See the security section for instructions on how to encrypt passwords.
- `name`: The name of the database or schema to connect to. For SQLite, this is the path to the database file.
- `tls`: A boolean (`true` or `false`) to enable or disable TLS/SSL for the connection.
- `health_query`: (Optional) A simple SQL query to execute to verify the connection is healthy (e.g., `"SELECT 1"`). This is not used for MongoDB.

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
    password: "ENCRYPTED_PASSWORD_HERE" # Replace with your encrypted password
    name: "mydatabase"
    tls: false
    health_query: "SELECT 1"

Run the check for this specific database:
* using a key file (recommended):
```bash
./dbchecker -key-file /path/to/your/secret.key -db my_postgres_db
```
* Alternatively, run the check using the environment variable:
```bash
export DB_SECRET_KEY="your-32-byte-secret-key"
./dbchecker -db my_postgres_db
```

### MongoDB Health Check

To check a MongoDB database, add an entry like this:

```yaml
databases:
  my_mongo_db:
    type: "mongodb"
    host: "localhost"
    port: 27017
    user: "mongouser"
    password: "ENCRYPTED_PASSWORD_HERE" # Replace with your encrypted password
    name: "admin"
    tls: false
```

Run the check for this specific database:
* using a key file (recommended):
```bash
./dbchecker -key-file /path/to/your/secret.key -db my_mongo_db
```
* Alternatively, run the check using the environment variable:
```bash
export DB_SECRET_KEY="your-32-byte-secret-key"
./dbchecker -db my_mongo_db
```

## Security and Credential Encryption

To avoid storing plaintext passwords in the configuration file, the application uses AES-GCM encryption to secure your credentials.
Providing the Secret Key

### How it Works

1.  **Secret Key**: The application uses a secret key for encryption and decryption. This key is provided at runtime via the `DB_SECRET_KEY` environment variable. **This key must be 32 bytes (256 bits) long.**
2.  **Encryption**: You use the `-encrypt` flag to encrypt your database password. The application uses the `DB_SECRET_KEY` to perform the encryption and gives you a base64-encoded string.
3.  **Configuration**: You paste this encrypted string into the `password` field in your `config.yaml`.
4.  **Decryption**: When the application runs, it reads the encrypted password from the config, reads the `DB_SECRET_KEY` from the environment, and decrypts the password in memory just before connecting to the database.

You can provide the secret key to the application in one of two ways.
* the `DB_SECRET_KEY` environment variable, passed through a XOR de-obfuscation step with a key compiled into the application.
  Key File (Recommended): You can store the key in a file and provide the path to it using the -key-file flag. This is the recommended approach for production environments.

* The -key-file flag always takes precedence over the DB_SECRET_KEY environment variable if both are present.


### How to Encrypt Your Credentials

**Step 1: Generate a Secure Secret Key**

First, you need a secure, 32-byte secret key. You can generate one using `openssl`.
 * To create a key file (recommended):

    On Linux or macOS:
```bash
    openssl rand 32 > /path/to/your/secret.key
```
    # Set strict, read-only permissions for the key file
```bash
    chmod 400 /path/to/your/secret.key
```
    The application will check for these exact permissions on non-Windows systems and will not run if the file is accessible by other users.

    On Windows: You should use the file system's ACL (Access Control List) features to restrict access to the current user. You can do this with the icacls command:

    # First, disable inheritance to remove other users
```bash
    icacls "C:\path\to\your\secret.key" /inheritance:r
```
    # Then, grant access only to your user
```bash
    icacls "C:\path\to\your\secret.key" /grant:r "$($env:USERNAME):(R)"
```

 * To generate a key for the environment variable:
```bash
    openssl rand -base64 32
```
** Step 2: Encrypt Your Password

Now, use the application's -encrypt flag to encrypt your database password.

  * Using the key file:
```bash
    ./dbchecker -key-file /path/to/your/secret.key -encrypt 'my-super-secret-password'
```
  * Using the environment variable:
```bash
    export DB_SECRET_KEY="your-generated-key"
    ./dbchecker -encrypt 'my-super-secret-password'
```
The application will output a long, base64-encoded string. This is your encrypted password.

** Step 3: Update Your Configuration

Copy the encrypted password from Step 2 and paste it into the password field in your config.yaml file.

Your application is now configured to securely connect to your databases!
