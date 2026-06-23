# Author Service (Golang)

This is the production-grade **Author Service** for the blog application, fully migrated from Express.js (TypeScript) to Golang. It is built using the standard Go directory layout and implements a decoupled **Service-Repository** pattern.

## Directory Structure 

The repository is organized according to clean architecture guidelines to support long-term extensibility:

```
author/
  ├── cmd/
  │   └── server/
  │       └── main.go           # Application entry point ( package main)
  ├── internal/
  │   ├── config/
  │   │   └── config.go         # Configuration parsing & fail-close check
  │   ├── logger/
  │   │   └── logger.go         # Zap structured JSON logging setup
  │   ├── models/
  │   │   └── models.go         # Shared model schemas (prevents circular imports)
  │   ├── db/
  │   │   ├── db.go             # PostgreSQL connection pool setup
  │   │   ├── queries.go        # Isolated SQL templates & AI prompts constants
  │   │   └── repository.go     # Low-level DB query maps
  │   ├── s3/
  │   │   └── s3.go             # AWS S3 file upload wrappers
  │   ├── rabbitmq/
  │   │   └── rabbitmq.go       # RabbitMQ publisher for cache invalidations
  │   ├── gemini/
  │   │   └── gemini.go         # REST client to Google Gemini API
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
2.  **AWS S3 File Storage**: Standard multipart uploads are streamed to S3 buckets, renaming files to unpredictable UUID formats to prevent name collision and directory traversal exploits.
3.  **RabbitMQ Event Publishing**: Automatically publishes events to the `cache-invalidation` queue to keep caching services up to date.
4.  **Google Gemini AI API Integration**: Grammar checking and title/description generation assistants powered by Gemini `gemini-2.5-flash` model requests.
5.  **Strict Security Configurations**:
    *   **Fail-Close Config**: Refuses to boot if critical configuration or secret credentials are missing.
    *   **IP-Based Rate Limiting**: Uses a thread-safe Token Bucket algorithm per client IP.
    *   **Body Size Limiting**: Restricts maximum body size requests to 10MB to prevent memory exhaustion / Denial-of-Service attacks.
    *   **Magic Byte Content Detection**: Reads file header magic bytes to ensure only valid images (`JPEG`, `PNG`, `GIF`, `WEBP`) can be uploaded.
    *   **Error Masking**: Diagnostic log logs database details to stdout via Zap while returning generic error messages (e.g., "Database error occurred") to clients.
    *   **Non-Root Docker Execution**: Docker container drops permissions to a non-root application user (`appuser`).

---

## Configuration

The service loads configurations from environment variables or a local `.env` file:

| Environment Variable | Default Value | Description |
| :--- | :--- | :--- |
| `PORT` | `5000` | HTTP Port to listen on |
| `DB_URL` | *Required* | PostgreSQL (Neon) Connection URL |
| `DB_MAX_CONNS` | `25` | Maximum database connection pool size |
| `DB_MIN_CONNS` | `5` | Minimum database connection pool size |
| `DB_MAX_CONN_LIFETIME` | `30m` | Maximum lifetime duration of a database connection |
| `DB_MAX_CONN_IDLE_TIME` | `15m` | Maximum idle time duration of a database connection |
| `Rabbitmq_Host` | *Required* | RabbitMQ broker host address |
| `Rabbitmq_Username` | *Required* | RabbitMQ broker username |
| `Rabbitmq_Password` | *Required* | RabbitMQ broker password |
| `AWS_ACCESS_KEY_ID` | *Required* | AWS Credentials - Access Key ID |
| `AWS_SECRET_ACCESS_KEY` | *Required* | AWS Credentials - Secret Access Key |
| `AWS_REGION` | *Required* | AWS S3 Bucket Region |
| `AWS_S3_BUCKET` | *Required* | AWS S3 Bucket Name |
| `JWT_SECRET` | *Required* | JWT Signature key |
| `GEMINI_API_KEY` | *Required* | Google Gemini API Key |
| `LOG_LEVEL` | `info` | Zap logging verbosity (`debug`, `info`, `warn`, `error`) |

---

## Running Locally

### Prerequisites
*   Go version 1.22+ (tested under Go `1.26.1`)
*   PostgreSQL & RabbitMQ instances

### Steps

1.  Create a `.env` file in the service root folder with your environment configuration values.
2.  Install dependencies:
    ```bash
    go mod download
    ```
3.  Run the application:
    ```bash
    go run ./cmd/server/main.go
    ```

## Running Tests

Automated tests verify middlewares (JWT, Rate Limiting, Size Limiting) and AI Gemini handlers (using local mock servers):

```bash
go test -v ./...
```

## Docker Deployment

Build and run the service using the optimized multi-stage secure Dockerfile:

```bash
# Build
docker build -t author-service .

# Run
docker run -p 5000:5000 --env-file .env author-service
```

---

## API Endpoints

All routes are served under `/api/v1` through the Chi Router:

*   `POST /api/v1/blog/new` (JWT Authenticated)
    *   Creates a new blog. Accepts `multipart/form-data` with fields `title`, `description`, `blogcontent`, `category`, and a file field `file`.
*   `POST /api/v1/blog/{id}` (JWT Authenticated)
    *   Updates an existing blog. Checks if the authenticated user is the blog author. Optional `file` upload replaces the header image.
*   `DELETE /api/v1/blog/{id}` (JWT Authenticated)
    *   Deletes a blog post and cascade-deletes comment/saved references in a SQL transaction.
*   `POST /api/v1/ai/title`
    *   Cleans and corrects grammar of titles. Body: `{"text": "title"}`
*   `POST /api/v1/ai/description`
    *   Grammar check or generates description. Body: `{"title": "title", "description": ""}`
*   `POST /api/v1/ai/blog`
    *   Cleans grammar of Jodit Editor HTML templates. Body: `{"blog": "<html>...</html>"}`
