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

	"blog/internal/db"
	"blog/internal/logger"
	"blog/internal/middleware"
	"blog/internal/models"
	"blog/internal/service"
)

type mockBlogRepository struct {
	db.BlogRepository
	GetAllBlogsFn    func(ctx context.Context, searchQuery, category string) ([]models.Blog, error)
	GetSingleBlogFn  func(ctx context.Context, id int) (*models.Blog, error)
	AddCommentFn     func(ctx context.Context, comment string, blogID string, userID, username string) error
	GetAllCommentsFn func(ctx context.Context, blogID string) ([]models.Comment, error)
	GetCommentByIDFn func(ctx context.Context, commentID int) (*models.Comment, error)
	DeleteCommentFn  func(ctx context.Context, commentID int) error
	GetSavedBlogFn   func(ctx context.Context, userID string, blogID string) (*models.SavedBlog, error)
	SaveBlogFn       func(ctx context.Context, userID string, blogID string) (bool, error)
	GetSavedBlogsFn  func(ctx context.Context, userID string) ([]models.SavedBlog, error)
}

func (m *mockBlogRepository) GetAllBlogs(ctx context.Context, searchQuery, category string) ([]models.Blog, error) {
	if m.GetAllBlogsFn != nil {
		return m.GetAllBlogsFn(ctx, searchQuery, category)
	}
	return []models.Blog{}, nil
}

func (m *mockBlogRepository) GetSingleBlog(ctx context.Context, id int) (*models.Blog, error) {
	if m.GetSingleBlogFn != nil {
		return m.GetSingleBlogFn(ctx, id)
	}
	return nil, nil
}

func (m *mockBlogRepository) AddComment(ctx context.Context, comment string, blogID string, userID, username string) error {
	if m.AddCommentFn != nil {
		return m.AddCommentFn(ctx, comment, blogID, userID, username)
	}
	return nil
}

func (m *mockBlogRepository) GetAllComments(ctx context.Context, blogID string) ([]models.Comment, error) {
	if m.GetAllCommentsFn != nil {
		return m.GetAllCommentsFn(ctx, blogID)
	}
	return []models.Comment{}, nil
}

func (m *mockBlogRepository) GetCommentByID(ctx context.Context, commentID int) (*models.Comment, error) {
	if m.GetCommentByIDFn != nil {
		return m.GetCommentByIDFn(ctx, commentID)
	}
	return nil, nil
}

func (m *mockBlogRepository) DeleteComment(ctx context.Context, commentID int) error {
	if m.DeleteCommentFn != nil {
		return m.DeleteCommentFn(ctx, commentID)
	}
	return nil
}

func (m *mockBlogRepository) GetSavedBlog(ctx context.Context, userID string, blogID string) (*models.SavedBlog, error) {
	if m.GetSavedBlogFn != nil {
		return m.GetSavedBlogFn(ctx, userID, blogID)
	}
	return nil, nil
}

func (m *mockBlogRepository) SaveBlog(ctx context.Context, userID string, blogID string) (bool, error) {
	if m.SaveBlogFn != nil {
		return m.SaveBlogFn(ctx, userID, blogID)
	}
	return false, nil
}

func (m *mockBlogRepository) GetSavedBlogs(ctx context.Context, userID string) ([]models.SavedBlog, error) {
	if m.GetSavedBlogsFn != nil {
		return m.GetSavedBlogsFn(ctx, userID)
	}
	return []models.SavedBlog{}, nil
}

func TestMain(m *testing.M) {
	logger.Logger = zap.NewNop()
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
		_, err := ioReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
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

		if w.Code == http.StatusOK {
			t.Errorf("expected status to fail or return error, got %d", w.Code)
		}
	})
}

func TestGetAllBlogs(t *testing.T) {
	mockRepo := &mockBlogRepository{
		GetAllBlogsFn: func(ctx context.Context, searchQuery, category string) ([]models.Blog, error) {
			return []models.Blog{
				{ID: 1, Title: "Test Blog 1", Author: "user1"},
				{ID: 2, Title: "Test Blog 2", Author: "user2"},
			}, nil
		},
	}

	svc := service.NewBlogService(mockRepo, "http://mock-user-service")
	h := NewBlogHandler(svc)

	req := httptest.NewRequest("GET", "/api/v1/blog/all?searchQuery=test", nil)
	w := httptest.NewRecorder()

	h.GetAllBlogs(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var resp map[string][]models.Blog
	err := json.NewDecoder(w.Body).Decode(&resp)
	if err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	blogs := resp["blogs"]
	if len(blogs) != 2 {
		t.Errorf("expected 2 blogs, got %d", len(blogs))
	}
	if blogs[0].Title != "Test Blog 1" {
		t.Errorf("expected title 'Test Blog 1', got '%s'", blogs[0].Title)
	}
}

func TestGetSingleBlog(t *testing.T) {
	mockRepo := &mockBlogRepository{
		GetSingleBlogFn: func(ctx context.Context, id int) (*models.Blog, error) {
			if id == 1 {
				return &models.Blog{ID: 1, Title: "Single Blog", Author: "user1"}, nil
			}
			return nil, nil
		},
	}

	mockUserServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"_id": "user1", "name": "Author Name"}`))
	}))
	defer mockUserServer.Close()

	svc := service.NewBlogService(mockRepo, mockUserServer.URL)
	h := NewBlogHandler(svc)

	t.Run("Blog Found", func(t *testing.T) {
		r := chi.NewRouter()
		r.Get("/blog/{id}", h.GetSingleBlog)

		req := httptest.NewRequest("GET", "/blog/1", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}

		var resp service.SingleBlogResponse
		err := json.NewDecoder(w.Body).Decode(&resp)
		if err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if resp.Blog.Title != "Single Blog" {
			t.Errorf("expected title 'Single Blog', got '%s'", resp.Blog.Title)
		}

		authorMap, ok := resp.Author.(map[string]interface{})
		if !ok || authorMap["name"] != "Author Name" {
			t.Errorf("expected author name 'Author Name', got %+v", resp.Author)
		}
	})

	t.Run("Blog Not Found", func(t *testing.T) {
		r := chi.NewRouter()
		r.Get("/blog/{id}", h.GetSingleBlog)

		req := httptest.NewRequest("GET", "/blog/999", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
		}
	})
}

func TestComments(t *testing.T) {
	mockRepo := &mockBlogRepository{
		AddCommentFn: func(ctx context.Context, comment string, blogID string, userID, username string) error {
			if comment == "my comment" && blogID == "1" && userID == "user1" {
				return nil
			}
			return errors.New("invalid comments params")
		},
		GetAllCommentsFn: func(ctx context.Context, blogID string) ([]models.Comment, error) {
			if blogID == "1" {
				return []models.Comment{
					{ID: 10, Comment: "comment 1", BlogID: "1"},
				}, nil
			}
			return []models.Comment{}, nil
		},
		GetCommentByIDFn: func(ctx context.Context, commentID int) (*models.Comment, error) {
			if commentID == 10 {
				return &models.Comment{ID: 10, UserID: "user1", Comment: "comment 1"}, nil
			}
			return nil, nil
		},
		DeleteCommentFn: func(ctx context.Context, commentID int) error {
			if commentID == 10 {
				return nil
			}
			return errors.New("db error")
		},
	}

	svc := service.NewBlogService(mockRepo, "")
	h := NewBlogHandler(svc)

	t.Run("Add Comment Success", func(t *testing.T) {
		r := chi.NewRouter()
		r.Post("/comment/{id}", h.AddComment)

		reqBody := `{"comment": "my comment"}`
		req := httptest.NewRequest("POST", "/comment/1", strings.NewReader(reqBody))
		ctx := middleware.WithUserContext(req.Context(), models.User{ID: "user1", Name: "User One"})
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req.WithContext(ctx))

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}

		var resp map[string]string
		_ = json.NewDecoder(w.Body).Decode(&resp)
		if resp["message"] != "Comment Added" {
			t.Errorf("expected 'Comment Added', got '%s'", resp["message"])
		}
	})

	t.Run("Get All Comments", func(t *testing.T) {
		r := chi.NewRouter()
		r.Get("/comment/{id}", h.GetAllComments)

		req := httptest.NewRequest("GET", "/comment/1", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}

		var comments []models.Comment
		_ = json.NewDecoder(w.Body).Decode(&comments)
		if len(comments) != 1 || comments[0].Comment != "comment 1" {
			t.Errorf("expected comment 'comment 1', got %+v", comments)
		}
	})

	t.Run("Delete Comment Success", func(t *testing.T) {
		r := chi.NewRouter()
		r.Delete("/comment/{commentid}", h.DeleteComment)

		req := httptest.NewRequest("DELETE", "/comment/10", nil)
		ctx := middleware.WithUserContext(req.Context(), models.User{ID: "user1"})
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req.WithContext(ctx))

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}

		var resp map[string]string
		_ = json.NewDecoder(w.Body).Decode(&resp)
		if resp["message"] != "comment deleted" {
			t.Errorf("expected 'comment deleted', got '%s'", resp["message"])
		}
	})

	t.Run("Delete Comment Forbidden (Not Owner)", func(t *testing.T) {
		r := chi.NewRouter()
		r.Delete("/comment/{commentid}", h.DeleteComment)

		req := httptest.NewRequest("DELETE", "/comment/10", nil)
		ctx := middleware.WithUserContext(req.Context(), models.User{ID: "other-user"})
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req.WithContext(ctx))

		if w.Code != http.StatusUnauthorized {
			t.Errorf("expected status %d, got %d", http.StatusUnauthorized, w.Code)
		}
	})
}

func TestSaveBlog(t *testing.T) {
	mockRepo := &mockBlogRepository{
		SaveBlogFn: func(ctx context.Context, userID string, blogID string) (bool, error) {
			if userID == "user1" && blogID == "5" {
				return true, nil
			}
			if userID == "user1" && blogID == "6" {
				return false, nil
			}
			return false, errors.New("invalid save params")
		},
		GetSavedBlogsFn: func(ctx context.Context, userID string) ([]models.SavedBlog, error) {
			if userID == "user1" {
				return []models.SavedBlog{
					{ID: 100, UserID: "user1", BlogID: "5"},
				}, nil
			}
			return []models.SavedBlog{}, nil
		},
	}

	svc := service.NewBlogService(mockRepo, "")
	h := NewBlogHandler(svc)

	t.Run("Save Blog - Created", func(t *testing.T) {
		r := chi.NewRouter()
		r.Post("/save/{blogid}", h.SaveBlog)

		req := httptest.NewRequest("POST", "/save/5", nil)
		ctx := middleware.WithUserContext(req.Context(), models.User{ID: "user1"})
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req.WithContext(ctx))

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}

		var resp map[string]string
		_ = json.NewDecoder(w.Body).Decode(&resp)
		if resp["message"] != "Blog saved" {
			t.Errorf("expected 'Blog saved', got '%s'", resp["message"])
		}
	})

	t.Run("Save Blog - Removed", func(t *testing.T) {
		r := chi.NewRouter()
		r.Post("/save/{blogid}", h.SaveBlog)

		req := httptest.NewRequest("POST", "/save/6", nil)
		ctx := middleware.WithUserContext(req.Context(), models.User{ID: "user1"})
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req.WithContext(ctx))

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}

		var resp map[string]string
		_ = json.NewDecoder(w.Body).Decode(&resp)
		if resp["message"] != "Blog removed from saved blogs" {
			t.Errorf("expected 'Blog removed from saved blogs', got '%s'", resp["message"])
		}
	})

	t.Run("Get Saved Blogs", func(t *testing.T) {
		r := chi.NewRouter()
		r.Get("/blogs/saved/all", h.GetSavedBlogs)

		req := httptest.NewRequest("GET", "/blogs/saved/all", nil)
		ctx := middleware.WithUserContext(req.Context(), models.User{ID: "user1"})
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req.WithContext(ctx))

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}

		var saved []models.SavedBlog
		_ = json.NewDecoder(w.Body).Decode(&saved)
		if len(saved) != 1 || saved[0].BlogID != "5" {
			t.Errorf("expected saved blog ID '5', got %+v", saved)
		}
	})
}

func ioReadAll(r io.Reader) ([]byte, error) {
	var buf bytes.Buffer
	_, err := buf.ReadFrom(r)
	return buf.Bytes(), err
}

