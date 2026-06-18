package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"blog/internal/logger"
	"blog/internal/middleware"
	"blog/internal/service"
)

type BlogHandler struct {
	service *service.BlogService
}

func NewBlogHandler(s *service.BlogService) *BlogHandler {
	return &BlogHandler{service: s}
}

// GetAllBlogs handler retrieves all blogs with optional query filtering and pagination
func (h *BlogHandler) GetAllBlogs(w http.ResponseWriter, r *http.Request) {
	searchQuery := r.URL.Query().Get("searchQuery")
	category := r.URL.Query().Get("category")
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 12
	if limitStr != "" {
		if val, err := strconv.Atoi(limitStr); err == nil && val > 0 {
			limit = val
		}
	}

	offset := 0
	if offsetStr != "" {
		if val, err := strconv.Atoi(offsetStr); err == nil && val >= 0 {
			offset = val
		}
	}

	blogs, err := h.service.GetAllBlogs(r.Context(), searchQuery, category, limit, offset)
	if err != nil {
		logger.Logger.Error("GetAllBlogs service call failed", zap.Error(err))
		middleware.JsonError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"blogs": blogs,
	})
}

// GetSingleBlog handler retrieves a single blog and details about its author
func (h *BlogHandler) GetSingleBlog(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		middleware.JsonError(w, http.StatusBadRequest, "Invalid blog ID")
		return
	}

	var userID string
	if user, ok := middleware.GetUserFromContext(r.Context()); ok {
		userID = user.ID
	}

	resp, err := h.service.GetSingleBlog(r.Context(), id, userID)
	if err != nil {
		logger.Logger.Error("GetSingleBlog service call failed", zap.Int("blog_id", id), zap.Error(err))
		middleware.JsonError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	if resp == nil {
		middleware.JsonError(w, http.StatusNotFound, "Blog not found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

type commentRequest struct {
	Comment string `json:"comment"`
}

// AddComment handler inserts a comment for a blog post
func (h *BlogHandler) AddComment(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		middleware.JsonError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	blogID := chi.URLParam(r, "id")
	if blogID == "" {
		middleware.JsonError(w, http.StatusBadRequest, "Missing blog ID")
		return
	}

	var req commentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Comment == "" {
		middleware.JsonError(w, http.StatusBadRequest, "Missing or invalid comment body")
		return
	}

	err := h.service.AddComment(r.Context(), req.Comment, blogID, user.ID, user.Name)
	if err != nil {
		logger.Logger.Error("AddComment service call failed", zap.Error(err))
		middleware.JsonError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"message": "Comment Added",
	})
}

// GetAllComments handler retrieves comments for a blog ID
func (h *BlogHandler) GetAllComments(w http.ResponseWriter, r *http.Request) {
	blogID := chi.URLParam(r, "id")
	if blogID == "" {
		middleware.JsonError(w, http.StatusBadRequest, "Missing blog ID")
		return
	}

	comments, err := h.service.GetAllComments(r.Context(), blogID)
	if err != nil {
		logger.Logger.Error("GetAllComments service call failed", zap.Error(err))
		middleware.JsonError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(comments)
}

// DeleteComment handler checks comment ownership and deletes it
func (h *BlogHandler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		middleware.JsonError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	commentIDStr := chi.URLParam(r, "commentid")
	commentID, err := strconv.Atoi(commentIDStr)
	if err != nil {
		middleware.JsonError(w, http.StatusBadRequest, "Invalid comment ID")
		return
	}

	err = h.service.DeleteComment(r.Context(), commentID, user.ID)
	if err != nil {
		logger.Logger.Error("DeleteComment service call failed", zap.Int("comment_id", commentID), zap.Error(err))
		if strings.Contains(err.Error(), "unauthorized") {
			middleware.JsonError(w, http.StatusUnauthorized, "You are not owner of this comment")
		} else if strings.Contains(err.Error(), "not found") {
			middleware.JsonError(w, http.StatusNotFound, "Comment not found")
		} else {
			middleware.JsonError(w, http.StatusInternalServerError, "Internal Server Error")
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"message": "comment deleted",
	})
}

// SaveBlog handler toggles saving/unsaving a blog post
func (h *BlogHandler) SaveBlog(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		middleware.JsonError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	blogID := chi.URLParam(r, "blogid")
	if blogID == "" {
		middleware.JsonError(w, http.StatusBadRequest, "Missing blog ID")
		return
	}

	saved, err := h.service.SaveBlog(r.Context(), user.ID, blogID)
	if err != nil {
		logger.Logger.Error("SaveBlog service call failed", zap.Error(err))
		middleware.JsonError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	var msg string
	if saved {
		msg = "Blog saved"
	} else {
		msg = "Blog removed from saved blogs"
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"message": msg,
	})
}

// GetSavedBlogs handler lists all blogs saved by user
func (h *BlogHandler) GetSavedBlogs(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		middleware.JsonError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	blogs, err := h.service.GetSavedBlogs(r.Context(), user.ID)
	if err != nil {
		logger.Logger.Error("GetSavedBlogs service call failed", zap.Error(err))
		middleware.JsonError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(blogs)
}

// LikeBlog handler toggles liking/unliking a blog post
func (h *BlogHandler) LikeBlog(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		middleware.JsonError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	blogID := chi.URLParam(r, "blogid")
	if blogID == "" {
		middleware.JsonError(w, http.StatusBadRequest, "Missing blog ID")
		return
	}

	liked, err := h.service.LikeBlog(r.Context(), user.ID, blogID)
	if err != nil {
		logger.Logger.Error("LikeBlog service call failed", zap.Error(err))
		middleware.JsonError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	var msg string
	if liked {
		msg = "Blog liked"
	} else {
		msg = "Blog unliked"
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"message": msg,
		"liked":   liked,
	})
}
