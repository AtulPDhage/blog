package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"

	"author/internal/gemini"
	"author/internal/logger"
	"author/internal/middleware"
	"author/internal/models"
	"author/internal/s3"
	"author/internal/service"
)

func TestMain(m *testing.M) {
	// Initialize logger to Nop to prevent nil pointer panics
	logger.Logger = zap.NewNop()

	// Mock S3 upload function globally for all tests
	s3.UploadImageFn = func(ctx context.Context, file io.Reader, originalFilename, contentType string) (string, error) {
		return "https://mock-bucket.s3.amazonaws.com/blogs/mock-file.png", nil
	}

	os.Exit(m.Run())
}

func TestAuthMiddleware(t *testing.T) {
	jwtSecret := "my-secret-key"
	authMiddleware := middleware.AuthMiddleware(jwtSecret)

	dummyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := middleware.GetUserFromContext(r.Context())
		if !ok {
			t.Errorf("expected user in context, got none")
		}
		if user.ID != "12345" {
			t.Errorf("expected user ID '12345', got '%s'", user.ID)
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	t.Run("Missing Auth Header", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()

		authMiddleware(dummyHandler).ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("expected status %d, got %d", http.StatusUnauthorized, w.Code)
		}
	})

	t.Run("Invalid Token", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer invalid-token")
		w := httptest.NewRecorder()

		authMiddleware(dummyHandler).ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("expected status %d, got %d", http.StatusUnauthorized, w.Code)
		}
	})

	t.Run("Valid Token", func(t *testing.T) {
		claims := &models.Claims{
			User: models.User{
				ID:   "12345",
				Name: "Test User",
			},
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte(jwtSecret))
		if err != nil {
			t.Fatalf("failed to generate token: %v", err)
		}

		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer "+tokenString)
		w := httptest.NewRecorder()

		authMiddleware(dummyHandler).ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}
	})
}

func TestMaxBodySizeMiddleware(t *testing.T) {
	authMiddleware := middleware.MaxBodySizeMiddleware(10) // 10 bytes limit

	dummyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := ioReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(body)
	})

	t.Run("Within Limit", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/test", strings.NewReader("hello"))
		w := httptest.NewRecorder()

		authMiddleware(dummyHandler).ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}
	})

	t.Run("Exceeds Limit", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/test", strings.NewReader("hello world, this is too long"))
		w := httptest.NewRecorder()

		authMiddleware(dummyHandler).ServeHTTP(w, req)

		// http.MaxBytesReader will trigger error during read
		if w.Code == http.StatusOK {
			t.Errorf("expected status to fail or return error, got %d", w.Code)
		}
	})
}

func ioReadAll(r ioReader) ([]byte, error) {
	var buf bytes.Buffer
	_, err := buf.ReadFrom(r)
	return buf.Bytes(), err
}

type ioReader interface {
	Read(p []byte) (n int, err error)
}

func TestCleanAISpanHelper(t *testing.T) {
	input := "**Hello** \r\n *World* `code` ~strikethrough~"
	expected := "Hello  World code strikethrough"
	output := CleanAISpanHelper(input)

	if output != expected {
		t.Errorf("expected '%s', got '%s'", expected, output)
	}
}

func TestAIHandlers(t *testing.T) {
	// 1. Mock Gemini API Server
	mockGeminiServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req gemini.GeminiRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		promptText := req.Contents[0].Parts[0].Text
		var replyText string

		if strings.Contains(promptText, "title") {
			replyText = "**Corrected Title**"
		} else if strings.Contains(promptText, "description") {
			replyText = "**Corrected Description**"
		} else if strings.Contains(promptText, "grammar correction engine") {
			replyText = "```html<p>Corrected HTML</p>```"
		} else {
			replyText = "Default Reply"
		}

		resp := gemini.GeminiResponse{
			Candidates: []gemini.Candidate{
				{
					Content: gemini.ResponseContent{
						Parts: []gemini.Part{
							{Text: replyText},
						},
					},
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer mockGeminiServer.Close()

	// Override Gemini base URL
	oldBaseURL := gemini.GeminiBaseURL
	gemini.GeminiBaseURL = mockGeminiServer.URL
	defer func() { gemini.GeminiBaseURL = oldBaseURL }()

	apiKey := "dummy-api-key"
	ctx := context.WithValue(context.Background(), middleware.GeminiContextKey, apiKey)

	h := NewBlogHandler(service.NewBlogService(nil))

	t.Run("AI Title Handler", func(t *testing.T) {
		reqPayload := models.AITitleRequest{Text: "bad title"}
		reqBody, _ := json.Marshal(reqPayload)

		req := httptest.NewRequest("POST", "/api/v1/ai/title", bytes.NewReader(reqBody))
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		h.AITitle(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}

		var resp string
		_ = json.NewDecoder(w.Body).Decode(&resp)

		if resp != "Corrected Title" {
			t.Errorf("expected 'Corrected Title', got '%s'", resp)
		}
	})

	t.Run("AI Description Handler", func(t *testing.T) {
		reqPayload := models.AIDescriptionRequest{Title: "Some Title", Description: "bad description"}
		reqBody, _ := json.Marshal(reqPayload)

		req := httptest.NewRequest("POST", "/api/v1/ai/description", bytes.NewReader(reqBody))
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		h.AIDescription(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}

		var resp string
		_ = json.NewDecoder(w.Body).Decode(&resp)

		if resp != "Corrected Description" {
			t.Errorf("expected 'Corrected Description', got '%s'", resp)
		}
	})

	t.Run("AI Blog Handler", func(t *testing.T) {
		reqPayload := models.AIBlogRequest{Blog: "<p>bad html</p>"}
		reqBody, _ := json.Marshal(reqPayload)

		req := httptest.NewRequest("POST", "/api/v1/ai/blog", bytes.NewReader(reqBody))
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		h.AIBlog(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}

		var resp models.AIBlogResponse
		_ = json.NewDecoder(w.Body).Decode(&resp)

		if resp.HTML != "<p>Corrected HTML</p>" {
			t.Errorf("expected '<p>Corrected HTML</p>', got '%s'", resp.HTML)
		}
	})
}
