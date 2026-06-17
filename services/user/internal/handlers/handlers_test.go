package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"

	"user/internal/db"
	"user/internal/google"
	"user/internal/s3"
	"user/internal/logger"
	"user/internal/middleware"
	"user/internal/models"
	"user/internal/service"
)

type mockUserRepository struct {
	db.UserRepository
	FindByEmailFn          func(ctx context.Context, email string) (*models.User, error)
	FindByIDFn             func(ctx context.Context, id string) (*models.User, error)
	CreateFn               func(ctx context.Context, name, email, image string) (*models.User, error)
	UpdateFn               func(ctx context.Context, id string, name, instagram, linkedin, facebook, bio string) (*models.User, error)
	UpdateProfilePictureFn func(ctx context.Context, id string, imageURL string) (*models.User, error)
}

func (m *mockUserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	if m.FindByEmailFn != nil {
		return m.FindByEmailFn(ctx, email)
	}
	return nil, nil
}

func (m *mockUserRepository) FindByID(ctx context.Context, id string) (*models.User, error) {
	if m.FindByIDFn != nil {
		return m.FindByIDFn(ctx, id)
	}
	return nil, nil
}

func (m *mockUserRepository) Create(ctx context.Context, name, email, image string) (*models.User, error) {
	if m.CreateFn != nil {
		return m.CreateFn(ctx, name, email, image)
	}
	return nil, nil
}

func (m *mockUserRepository) Update(ctx context.Context, id string, name, instagram, linkedin, facebook, bio string) (*models.User, error) {
	if m.UpdateFn != nil {
		return m.UpdateFn(ctx, id, name, instagram, linkedin, facebook, bio)
	}
	return nil, nil
}

func (m *mockUserRepository) UpdateProfilePicture(ctx context.Context, id string, imageURL string) (*models.User, error) {
	if m.UpdateProfilePictureFn != nil {
		return m.UpdateProfilePictureFn(ctx, id, imageURL)
	}
	return nil, nil
}

func TestMain(m *testing.M) {
	logger.Logger = zap.NewNop()
	// Mock S3 globally
	s3.UploadImageFn = func(ctx context.Context, file io.Reader, originalFilename, contentType string) (string, error) {
		return "https://mock-bucket.s3.amazonaws.com/profile_pictures/mock-file.png", nil
	}
	os.Exit(m.Run())
}

func TestAuthMiddleware(t *testing.T) {
	jwtSecret := "user-jwt-secret"
	authMiddleware := middleware.AuthMiddleware(jwtSecret)

	dummyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := middleware.GetUserFromContext(r.Context())
		if !ok {
			t.Errorf("expected user in context, got none")
		}
		if user.ID != "507f1f77bcf86cd799439011" {
			t.Errorf("expected user ID '507f1f77bcf86cd799439011', got '%s'", user.ID)
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

	t.Run("Valid Token", func(t *testing.T) {
		claims := &models.Claims{
			User: models.User{
				ID:   "507f1f77bcf86cd799439011",
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

func TestMyProfile(t *testing.T) {
	h := NewUserHandler(service.NewUserService(nil, nil, "secret"))

	r := chi.NewRouter()
	r.Get("/me", h.MyProfile)

	req := httptest.NewRequest("GET", "/me", nil)
	ctx := middleware.WithUserContext(req.Context(), models.User{ID: "507f1f77bcf86cd799439011", Name: "Me"})
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req.WithContext(ctx))

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var resp map[string]models.User
	_ = json.NewDecoder(w.Body).Decode(&resp)
	if resp["user"].Name != "Me" {
		t.Errorf("expected name 'Me', got '%s'", resp["user"].Name)
	}
}

func TestGetUserProfile(t *testing.T) {
	mockRepo := &mockUserRepository{
		FindByIDFn: func(ctx context.Context, id string) (*models.User, error) {
			if id == "507f1f77bcf86cd799439011" {
				return &models.User{ID: id, Name: "Found User"}, nil
			}
			return nil, nil
		},
	}

	svc := service.NewUserService(mockRepo, nil, "secret")
	h := NewUserHandler(svc)

	r := chi.NewRouter()
	r.Get("/user/{id}", h.GetUserProfile)

	t.Run("User Found", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/user/507f1f77bcf86cd799439011", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}

		var resp map[string]models.User
		_ = json.NewDecoder(w.Body).Decode(&resp)
		if resp["user"].Name != "Found User" {
			t.Errorf("expected name 'Found User', got '%s'", resp["user"].Name)
		}
	})

	t.Run("User Not Found", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/user/507f1f77bcf86cd799439022", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("expected status 404, got %d", w.Code)
		}
	})
}

func TestUpdateUser(t *testing.T) {
	mockRepo := &mockUserRepository{
		UpdateFn: func(ctx context.Context, id string, name, instagram, linkedin, facebook, bio string) (*models.User, error) {
			if id == "507f1f77bcf86cd799439011" {
				return &models.User{
					ID:        id,
					Name:      name,
					Instagram: instagram,
					Linkedin:  linkedin,
					Facebook:  facebook,
					Bio:       bio,
				}, nil
			}
			return nil, errors.New("not found")
		},
	}

	svc := service.NewUserService(mockRepo, nil, "secret")
	h := NewUserHandler(svc)

	r := chi.NewRouter()
	r.Post("/user/update", h.UpdateUser)

	reqPayload := updateRequest{
		Name:      "New Name",
		Instagram: "insta",
		Linkedin:  "link",
		Facebook:  "fb",
		Bio:       "bio info",
	}
	body, _ := json.Marshal(reqPayload)

	req := httptest.NewRequest("POST", "/user/update", bytes.NewReader(body))
	ctx := middleware.WithUserContext(req.Context(), models.User{ID: "507f1f77bcf86cd799439011"})
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req.WithContext(ctx))

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	_ = json.NewDecoder(w.Body).Decode(&resp)

	if resp["message"] != "Profile updated successfully" {
		t.Errorf("expected success message, got '%v'", resp["message"])
	}
}

func TestLoginUser(t *testing.T) {
	mockRepo := &mockUserRepository{
		FindByEmailFn: func(ctx context.Context, email string) (*models.User, error) {
			if email == "test@google.com" {
				return &models.User{ID: "507f1f77bcf86cd799439011", Email: email, Name: "Existing User"}, nil
			}
			return nil, nil
		},
	}

	// Mock HTTP Server for Google API endpoints
	mockGoogleAPIServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if strings.Contains(r.URL.Path, "token") {
			_, _ = w.Write([]byte(`{"access_token": "mock-access-token", "expires_in": 3600}`))
		} else if strings.Contains(r.URL.Path, "userinfo") {
			_, _ = w.Write([]byte(`{"id": "google-id-123", "email": "test@google.com", "name": "Google User", "picture": "http://img.jpg"}`))
		}
	}))
	defer mockGoogleAPIServer.Close()

	// Re-route token & userinfo requests to mock server by swapping client endpoint configurations or replacing client
	gClient := google.NewGoogleClient("client-id", "client-secret")
	// Swap google client inside test: we need to use a custom client that maps queries locally.
	// Since we defined SetHTTPClient in google.GoogleClient, we can swap the http client and override the URLs.
	// Wait, inside ExchangeCodeForToken we request "https://oauth2.googleapis.com/token".
	// Since that URL is hardcoded, how can we mock it?
	// We can write a custom http.RoundTripper that intercepts all HTTP queries and maps them to the mock server!
	// This is an extremely elegant Go unit testing pattern!
	customClient := &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			// Redirect requests to mockGoogleAPIServer
			targetURL := mockGoogleAPIServer.URL
			if strings.Contains(req.URL.Host, "oauth2.googleapis.com") {
				targetURL += "/token"
			} else if strings.Contains(req.URL.Host, "googleapis.com") {
				targetURL += "/userinfo"
			}

			mockReq, _ := http.NewRequest(req.Method, targetURL, req.Body)
			mockReq.Header = req.Header
			return http.DefaultClient.Do(mockReq)
		}),
	}
	gClient.SetHTTPClient(customClient)

	svc := service.NewUserService(mockRepo, gClient, "secret")
	h := NewUserHandler(svc)

	reqPayload := loginRequest{Code: "auth-code"}
	body, _ := json.Marshal(reqPayload)

	req := httptest.NewRequest("POST", "/login", bytes.NewReader(body))
	w := httptest.NewRecorder()

	h.LoginUser(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	_ = json.NewDecoder(w.Body).Decode(&resp)

	if resp["message"] != "Login successful" {
		t.Errorf("expected message 'Login successful', got '%v'", resp["message"])
	}
}

type roundTripFunc func(req *http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestValidateImageFile(t *testing.T) {
	t.Run("Valid JPEG", func(t *testing.T) {
		// JPEG Magic Bytes: FF D8 FF
		imgBytes := []byte{0xff, 0xd8, 0xff, 0xe0, 0x00, 0x10, 0x4a, 0x46, 0x49, 0x46}
		r := bytes.NewReader(imgBytes)
		valid, err := ValidateImageFile(r)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !valid {
			t.Errorf("expected JPEG image to be valid")
		}
	})

	t.Run("Invalid Format", func(t *testing.T) {
		txtBytes := []byte("plain text file contents that are not an image")
		r := bytes.NewReader(txtBytes)
		valid, err := ValidateImageFile(r)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if valid {
			t.Errorf("expected text file to be invalid image")
		}
	})
}

func TestUpdateProfilePicture(t *testing.T) {
	mockRepo := &mockUserRepository{
		UpdateProfilePictureFn: func(ctx context.Context, id string, imageURL string) (*models.User, error) {
			if id == "507f1f77bcf86cd799439011" {
				return &models.User{ID: id, Image: imageURL}, nil
			}
			return nil, errors.New("not found")
		},
	}

	svc := service.NewUserService(mockRepo, nil, "secret")
	h := NewUserHandler(svc)

	r := chi.NewRouter()
	r.Post("/user/update/pic", h.UpdateProfilePicture)

	// Write multipart request body containing JPEG magic bytes
	body := &bytes.Buffer{}
	writer := mimeMultipartWriter(body) // wait, to avoid importing mime/multipart, we can build custom mock writer or import it
	part, _ := writer.CreateFormFile("file", "test.jpg")
	// Write JPEG header magic bytes to pass ValidateImageFile
	_, _ = part.Write([]byte{0xff, 0xd8, 0xff, 0xe0, 0x00, 0x10, 0x4a, 0x46, 0x49, 0x46})
	_ = writer.Close()

	req := httptest.NewRequest("POST", "/user/update/pic", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	ctx := middleware.WithUserContext(req.Context(), models.User{ID: "507f1f77bcf86cd799439011"})
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req.WithContext(ctx))

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	_ = json.NewDecoder(w.Body).Decode(&resp)

	if resp["message"] != "Profile picture updated successfully" {
		t.Errorf("expected success message, got '%v'", resp["message"])
	}

	userMap, ok := resp["user"].(map[string]interface{})
	if !ok || userMap["image"] != "https://mock-bucket.s3.amazonaws.com/profile_pictures/mock-file.png" {
		t.Errorf("expected mock s3 URL, got %+v", resp["user"])
	}
}

// wait, helper to mock mime/multipart
func mimeMultipartWriter(w io.Writer) *mockMultipartWriter {
	boundary := "foo-bar-boundary"
	return &mockMultipartWriter{
		w:        w,
		boundary: boundary,
	}
}

type mockMultipartWriter struct {
	w        io.Writer
	boundary string
}

func (m *mockMultipartWriter) FormDataContentType() string {
	return "multipart/form-data; boundary=" + m.boundary
}

func (m *mockMultipartWriter) CreateFormFile(fieldname, filename string) (io.Writer, error) {
	_, _ = io.WriteString(m.w, "--"+m.boundary+"\r\n")
	_, _ = io.WriteString(m.w, "Content-Disposition: form-data; name=\""+fieldname+"\"; filename=\""+filename+"\"\r\n")
	_, _ = io.WriteString(m.w, "Content-Type: image/jpeg\r\n\r\n")
	return m.w, nil
}

func (m *mockMultipartWriter) Close() error {
	_, _ = io.WriteString(m.w, "\r\n--"+m.boundary+"--\r\n")
	return nil
}

func ioReadAll(r io.Reader) ([]byte, error) {
	var buf bytes.Buffer
	_, err := buf.ReadFrom(r)
	return buf.Bytes(), err
}
