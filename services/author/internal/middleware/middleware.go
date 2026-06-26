package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"author/internal/models"
)

type contextKey string

const (
	userContextKey   contextKey = "user"
	GeminiContextKey contextKey = "GeminiAPIKey"
)

// GetUserFromContext retrieves the authenticated user from the context
func GetUserFromContext(ctx context.Context) (models.User, bool) {
	u, ok := ctx.Value(userContextKey).(models.User)
	return u, ok
}

// AuthMiddleware validates JWT Bearer tokens
func AuthMiddleware(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				_, _ = w.Write([]byte(`{"message": "Unauthorized"}`))
				return
			}

			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			claims := &models.Claims{}

			token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
				// Hardcode/enforce expected signing method (HS256) per secure web guidelines
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(jwtSecret), nil
			})

			if err != nil || !token.Valid || claims.User.ID == "" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				_, _ = w.Write([]byte(`{"message": "Unauthorized"}`))
				return
			}

			ctx := context.WithValue(r.Context(), userContextKey, claims.User)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// CORSMiddleware handles cross-origin requests
func CORSMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RateLimitMiddleware enforces IP-based rate limits (Token Bucket)
type clientLimiter struct {
	tokens     float64
	lastUpdate time.Time
}

func RateLimitMiddleware() func(http.Handler) http.Handler {
	var mu sync.Mutex
	clients := make(map[string]*clientLimiter)

	// Rate limit parameters: max 100 tokens, refills at 2 tokens per second
	const maxTokens = 100.0
	const refillRate = 2.0 // tokens/sec

	// Cleanup old entries periodically to avoid memory leaks
	go func() {
		for {
			time.Sleep(10 * time.Minute)
			mu.Lock()
			for ip, limiter := range clients {
				if time.Since(limiter.lastUpdate) > 1*time.Hour {
					delete(clients, ip)
				}
			}
			mu.Unlock()
		}
	}()

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				ip = r.RemoteAddr
			}

			mu.Lock()
			limiter, exists := existsInMap(clients, ip)
			now := time.Now()

			if !exists {
				limiter = &clientLimiter{
					tokens:     maxTokens,
					lastUpdate: now,
				}
				clients[ip] = limiter
			} else {
				// Refill tokens based on elapsed time
				elapsed := now.Sub(limiter.lastUpdate).Seconds()
				limiter.tokens += elapsed * refillRate
				if limiter.tokens > maxTokens {
					limiter.tokens = maxTokens
				}
				limiter.lastUpdate = now
			}

			if limiter.tokens < 1.0 {
				mu.Unlock()
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusTooManyRequests)
				_, _ = w.Write([]byte(`{"message": "Too many requests. Please try again later."}`))
				return
			}

			limiter.tokens -= 1.0
			mu.Unlock()

			next.ServeHTTP(w, r)
		})
	}
}

func existsInMap(m map[string]*clientLimiter, key string) (*clientLimiter, bool) {
	val, ok := m[key]
	return val, ok
}

// MaxBodySizeMiddleware limits request body size to prevent memory exhaustion / DoS
func MaxBodySizeMiddleware(maxBytes int64) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.Body = http.MaxBytesReader(w, r.Body, maxBytes)
			next.ServeHTTP(w, r)
		})
	}
}

// JsonError helper returns error response in JSON format
func JsonError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": message})
}
