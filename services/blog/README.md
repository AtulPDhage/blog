# Blog Service (Golang)

This is the production-grade **Blog Service** for the blog application, fully migrated from Express.js (TypeScript) to Golang. It is built using the standard Go directory layout and implements a decoupled **Service-Repository** pattern.

## Directory Structure

The repository is organized according to clean architecture guidelines to support long-term extensibility:

```
blog/
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
  │   │   ├── db.go             # PostgreSQL connection pool setup
  │   │   ├── queries.go        # Isolated SQL templates
  │   │   └── repository.go     # Database Repository Layer
  │   ├── redis/
  │   │   └── redis.go          # Redis connection pool & cache operations
  │   ├── rabbitmq/
  │   │   ├── rabbitmq.go       # RabbitMQ connection setup
  │   │   └── consumer.go       # Cache-invalidation consumer event listener
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
2.  **Redis Caching & Pre-warming**: Leverages Redis (`go-redis/v9`) to cache blog listings and details. The service listens on RabbitMQ for cache-invalidation events, invalidates matching keys, and proactively pre-warms database listing queries in Redis for snappy UI load times.
3.  **RabbitMQ Consumer**: Subscribes to the `cache-invalidation` queue to listen to cache events published by the `author` microservice.
4.  **Downstream User Service Integration**: Queries the downstream User microservice via HTTP to fetch author profile details and embeds them inside single blog response payloads.
5.  **Strict Security Configurations**:
    *   **Fail-Close Config**: Refuses to boot if critical configuration or secret credentials are missing.
    *   **IP-Based Rate Limiting**: Uses a thread-safe Token Bucket algorithm per client IP.
    *   **Body Size Limiting**: Restricts maximum body size requests to 10MB to prevent memory exhaustion / Denial-of-Service attacks.
    *   **JWT Token Verification**: Validates JWT signature keys using HS256 algorithm and unpacks user context details safely.
    *   **Error Masking**: Diagnostic logs record database details to stdout via Zap while returning generic error messages (e.g., "Internal Server Error") to clients.
    *   **Non-Root Docker Execution**: Docker container drops permissions to a non-root application user (`appuser`).

---

## Configuration

The service loads configurations from environment variables or a local `.env` file:

| Environment Variable | Default Value | Description |
| :--- | :--- | :--- |
| `PORT` | `5001` | HTTP Port to listen on |
| `DB_URL` | *Required* | PostgreSQL Connection URL |
| `DB_MAX_CONNS` | `25` | Maximum database connection pool size |
| `DB_MIN_CONNS` | `5` | Minimum database connection pool size |
| `DB_MAX_CONN_LIFETIME` | `30m` | Maximum lifetime duration of a database connection |
| `DB_MAX_CONN_IDLE_TIME` | `15m` | Maximum idle time duration of a database connection |
| `Rabbitmq_Host` | *Required* | RabbitMQ broker host address |
| `Rabbitmq_Username` | *Required* | RabbitMQ broker username |
| `Rabbitmq_Password` | *Required* | RabbitMQ broker password |
| `REDIS_URL` | *Required* | Redis connection pool URL |
| `USER_SERVICE` | *Required* | User service base URL |
| `JWT_SECRET` | *Required* | JWT Signature key |
| `LOG_LEVEL` | `info` | Zap logging verbosity (`debug`, `info`, `warn`, `error`) |

---

## Running Locally

### Prerequisites
*   Go version 1.22+ (tested under Go `1.26.1`)
*   PostgreSQL, Redis, and RabbitMQ instances

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
