package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"

	"blog/internal/db"
	"blog/internal/logger"
	"blog/internal/models"
	"blog/internal/redis"
)

type BlogService struct {
	repo           db.BlogRepository
	userServiceURL string
	httpClient     *http.Client
}

// NewBlogService returns a new instance of BlogService
func NewBlogService(repo db.BlogRepository, userServiceURL string) *BlogService {
	return &BlogService{
		repo:           repo,
		userServiceURL: userServiceURL,
		httpClient:     &http.Client{Timeout: 10 * time.Second},
	}
}

// SetHTTPClient allows overriding the default HTTP client (useful for unit tests/mocks)
func (s *BlogService) SetHTTPClient(client *http.Client) {
	s.httpClient = client
}

// GetAllBlogs fetches blogs with cache check and fallback to DB
func (s *BlogService) GetAllBlogs(ctx context.Context, searchQuery, category string, limit, offset int) ([]models.Blog, error) {
	cacheKey := fmt.Sprintf("blogs:%s:%s:%d:%d", searchQuery, category, limit, offset)

	// Check Redis cache first
	cached, err := redis.Get(ctx, cacheKey)
	if err == nil && cached != "" {
		var blogs []models.Blog
		if json.Unmarshal([]byte(cached), &blogs) == nil {
			logger.Logger.Info("Serving blogs list from Redis cache", zap.String("key", cacheKey))
			return blogs, nil
		}
	}

	// Fallback to Database
	blogs, err := s.repo.GetAllBlogs(ctx, searchQuery, category, limit, offset)
	if err != nil {
		return nil, err
	}

	logger.Logger.Info("Serving blogs list from Database", zap.String("key", cacheKey))

	// Async update to cache
	go func() {
		bJSON, err := json.Marshal(blogs)
		if err == nil {
			cacheCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			_ = redis.Set(cacheCtx, cacheKey, string(bJSON), 3600*time.Second)
		}
	}()

	return blogs, nil
}

// SingleBlogResponse holds the blog and its author details
type SingleBlogResponse struct {
	Blog   models.Blog `json:"blog"`
	Author interface{} `json:"author"`
}

// GetSingleBlog fetches a blog by ID, checks cache, retrieves author details from User microservice, and caches result
func (s *BlogService) GetSingleBlog(ctx context.Context, id int) (*SingleBlogResponse, error) {
	cacheKey := fmt.Sprintf("blog:%d", id)

	// Check Redis cache
	cached, err := redis.Get(ctx, cacheKey)
	if err == nil && cached != "" {
		var resp SingleBlogResponse
		if json.Unmarshal([]byte(cached), &resp) == nil {
			logger.Logger.Info("Serving single blog from Redis cache", zap.String("key", cacheKey))
			return &resp, nil
		}
	}

	// Fallback to DB
	blog, err := s.repo.GetSingleBlog(ctx, id)
	if err != nil {
		return nil, err
	}
	if blog == nil {
		return nil, nil // Not found
	}

	logger.Logger.Info("Serving single blog from Database", zap.String("key", cacheKey))

	// Fetch author details from external User service
	var authorDetails interface{}
	urlStr := fmt.Sprintf("%s/api/v1/user/%s", s.userServiceURL, blog.Author)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlStr, nil)
	if err == nil {
		resp, err := s.httpClient.Do(req)
		if err == nil {
			defer resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				var parsed interface{}
				if json.NewDecoder(resp.Body).Decode(&parsed) == nil {
					authorDetails = parsed
				}
			} else {
				logger.Logger.Warn("User service returned non-200 status", zap.Int("status", resp.StatusCode), zap.String("url", urlStr))
			}
		} else {
			logger.Logger.Error("Failed to request user service", zap.Error(err), zap.String("url", urlStr))
		}
	}

	if authorDetails == nil {
		authorDetails = map[string]interface{}{} // fallback empty map
	}

	finalResp := &SingleBlogResponse{
		Blog:   *blog,
		Author: authorDetails,
	}

	// Cache result
	go func() {
		respJSON, err := json.Marshal(finalResp)
		if err == nil {
			cacheCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			_ = redis.Set(cacheCtx, cacheKey, string(respJSON), 3600*time.Second)
		}
	}()

	return finalResp, nil
}

// AddComment inserts a comment
func (s *BlogService) AddComment(ctx context.Context, comment string, blogID string, userID, username string) error {
	if comment == "" || blogID == "" || userID == "" || username == "" {
		return errors.New("missing required comment fields")
	}
	return s.repo.AddComment(ctx, comment, blogID, userID, username)
}

// GetAllComments retrieves all comments for a blog
func (s *BlogService) GetAllComments(ctx context.Context, blogID string) ([]models.Comment, error) {
	if blogID == "" {
		return nil, errors.New("missing blog ID")
	}
	return s.repo.GetAllComments(ctx, blogID)
}

// DeleteComment checks ownership and deletes a comment
func (s *BlogService) DeleteComment(ctx context.Context, commentID int, userID string) error {
	existing, err := s.repo.GetCommentByID(ctx, commentID)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("comment not found")
	}
	if existing.UserID != userID {
		return errors.New("unauthorized to delete this comment")
	}

	return s.repo.DeleteComment(ctx, commentID)
}

// SaveBlog toggles the saved state of a blog
func (s *BlogService) SaveBlog(ctx context.Context, userID string, blogID string) (bool, error) {
	if userID == "" || blogID == "" {
		return false, errors.New("missing user ID or blog ID")
	}
	return s.repo.SaveBlog(ctx, userID, blogID)
}

// GetSavedBlogs retrieves all saved blogs
func (s *BlogService) GetSavedBlogs(ctx context.Context, userID string) ([]models.Blog, error) {
	if userID == "" {
		return nil, errors.New("missing user ID")
	}
	return s.repo.GetSavedBlogs(ctx, userID)
}
