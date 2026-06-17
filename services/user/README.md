# User Service (Golang)

This is the production-grade **User Service** for the blog application, fully migrated from Express.js (TypeScript) to Golang. It is built using the standard Go directory layout and implements a decoupled **Service-Repository** pattern.

## Directory Structure

The repository is organized according to clean architecture guidelines to support long-term extensibility:

```
user/
  ├── cmd/
  │   └── server/
  │       └── main.go           # Application entry point (package main)
  ├── internal/
  │   ├── config/
  │   │   └── config.go         # Configuration parsing & fail-close check
  │   ├── logger/
  │   │   └── logger.go         # Zap structured JSON logging setup
  │   ├── models/
  │   │   └── models.go         # Shared model schemas (prevents circular imports)
  │   ├── db/
  │   │   ├── db.go             # MongoDB driver connection pool setup
  │   │   └── repository.go     # Database Repository Layer
  │   ├── s3/
  │   │   └── s3.go             # AWS S3 media uploader client
  │   ├── google/
  │   │   └── google.go         # Google OAuth2 API exchange client
  │   ├── middleware/
  │   │   └── middleware.go     # CORS, Auth filters, Rate Limits, and Request sizes
  │   ├── service/
  │   │   └── service.go        # Decoupled business and domain orchestrations
  │   └── handlers/
  │       ├── handlers.go       # Chi-based HTTP controllers
  │       └── handlers_test.go  # Unit & integration tests
  ├── go.mod
  ├── go.sum
  └── Dockerfile
```

---

## Key Features & Security Design

1.  **Decoupled Architecture**: Separation of concerns ensures that Handlers only parse requests/responses, Services contain business rules, and Repositories run database queries.
2.  **MongoDB Database**: Handles document storage using the official `go.mongodb.org/mongo-driver` MongoDB client, providing thread-safe connection pooling, ping verification, and structured BSON parsing.
3.  **AWS S3 Storage**: Provides clean interfaces to stream memory form file buffers to AWS S3, generating random UUID filename hashes to prevent namespace collisions.
4.  **Google OAuth2 Exchange**: Implements lightweight exchange routines using HTTP client requests to handshake authorization codes and download user details.
5.  **Strict Security Configurations**:
    *   **Fail-Close Config**: Refuses to boot if critical configuration or secret credentials are missing.
    *   **IP-Based Rate Limiting**: Uses a thread-safe Token Bucket algorithm per client IP.
    *   **Body Size Limiting**: Restricts maximum body size requests to 10MB to prevent memory exhaustion / Denial-of-Service attacks.
    *   **Magic Byte Image Validation**: Verifies file content magic bytes before uploading to AWS S3 to ensure only valid files (`JPEG`, `PNG`, `GIF`, `WEBP`) can be loaded.
    *   **JWT Token Verification**: Validates JWT signature keys using HS256 algorithm and unpacks user context details safely.
    *   **Error Masking**: Diagnostic logs record database details to stdout via Zap while returning generic error messages (e.g., "Internal Server Error") to clients.
    *   **Non-Root Docker Execution**: Docker container drops permissions to a non-root application user (`appuser`).

---

## Configuration

The service loads configurations from environment variables or a local `.env` file:

| Environment Variable | Default Value | Description |
| :--- | :--- | :--- |
| `PORT` | `5002` | HTTP Port to listen on |
| `MONGO_URI` | *Required* | MongoDB Connection URL |
| `DB_NAME` | `MasterDB` | MongoDB Database Name |
| `AWS_ACCESS_KEY_ID` | *Required* | AWS Credentials - Access Key ID |
| `AWS_SECRET_ACCESS_KEY` | *Required* | AWS Credentials - Secret Access Key |
| `AWS_REGION` | *Required* | AWS S3 Bucket Region |
| `AWS_S3_BUCKET` | *Required* | AWS S3 Bucket Name |
| `Google_Client_Id` | *Required* | Google OAuth Client ID |
| `Google_client_Secret` | *Required* | Google OAuth Client Secret |
| `JWT_SECRET` | *Required* | JWT Signature key |
| `LOG_LEVEL` | `info` | Zap logging verbosity (`debug`, `info`, `warn`, `error`) |

---

## Running Locally

### Prerequisites
*   Go version 1.22+ (tested under Go `1.26.1`)
*   MongoDB, AWS S3 bucket, and Google client credentials

### Steps

1.  Create a `.env` file in the service root folder with your environment configuration values.
2.  Install dependencies:
    ```bash
    go mod download
    ```
3.  Run the application:
    ```bash
    go run cmd/server/main.go
    ```

---

## Running Verification Tests

Unit and integration tests verify authentication middlewares, request limits, mock repositories, and HTTP clients calling external user service endpoints.

Run all tests:
```bash
go test -v ./...
```
