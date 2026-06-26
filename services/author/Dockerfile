FROM golang:1.26-alpine AS builder

WORKDIR /app

# Copy dependency definitions
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code layouts
COPY cmd/ ./cmd/
COPY internal/ ./internal/

# Build optimized, statically-linked binary from cmd entrypoint
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o author-service ./cmd/server/main.go

# Production runner image
FROM alpine:3.19

RUN apk add --no-cache ca-certificates

WORKDIR /app

# Add non-root application user for secure running
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser

# Copy built binary
COPY --from=builder /app/author-service .

EXPOSE 5000

CMD ["./author-service"]