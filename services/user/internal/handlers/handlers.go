package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"user/internal/logger"
	"user/internal/middleware"
	"user/internal/service"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(s *service.UserService) *UserHandler {
	return &UserHandler{service: s}
}

// ValidateImageFile reads the first 512 bytes of a reader to verify it's an allowed image type
func ValidateImageFile(file io.Reader) (bool, error) {
	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return false, err
	}

	contentType := http.DetectContentType(buffer[:n])
	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/gif":  true,
		"image/webp": true,
	}

	if seeker, ok := file.(io.ReadSeeker); ok {
		_, _ = seeker.Seek(0, io.SeekStart)
	}

	return allowedTypes[contentType], nil
}

type loginRequest struct {
	Code string `json:"code"`
}

// LoginUser exchanges google auth code, registers/logs user, and issues JWT
func (h *UserHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Code == "" {
		middleware.JsonError(w, http.StatusBadRequest, "Authorization code is required")
		return
	}

	user, token, err := h.service.LoginGoogle(r.Context(), req.Code)
	if err != nil {
		logger.Logger.Error("Google login failed", zap.Error(err))
		middleware.JsonError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Login successful",
		"token":   token,
		"user":    user,
	})
}

// MyProfile returns current authenticated user details
func (h *UserHandler) MyProfile(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		middleware.JsonError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"user": user,
	})
}

// GetUserProfile loads any user details by ID
func (h *UserHandler) GetUserProfile(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		middleware.JsonError(w, http.StatusBadRequest, "Missing user ID")
		return
	}

	user, err := h.service.GetUserProfile(r.Context(), id)
	if err != nil {
		logger.Logger.Error("GetUserProfile failed", zap.String("user_id", id), zap.Error(err))
		middleware.JsonError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	if user == nil {
		middleware.JsonError(w, http.StatusNotFound, "User not found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"user": user,
	})
}

type updateRequest struct {
	Name      string `json:"name"`
	Instagram string `json:"instagram"`
	Linkedin  string `json:"linkedin"`
	Facebook  string `json:"facebook"`
	Bio       string `json:"bio"`
}

// UpdateUser modifies user profile details
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		middleware.JsonError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req updateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.JsonError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	updatedUser, token, err := h.service.UpdateUser(r.Context(), user.ID, req.Name, req.Instagram, req.Linkedin, req.Facebook, req.Bio)
	if err != nil {
		logger.Logger.Error("UpdateUser service call failed", zap.String("user_id", user.ID), zap.Error(err))
		middleware.JsonError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Profile updated successfully",
		"user":    updatedUser,
		"token":   token,
	})
}

// UpdateProfilePicture handles multipart upload and AWS S3 synchronization
func (h *UserHandler) UpdateProfilePicture(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		middleware.JsonError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Parse multipart form (max 10MB)
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		logger.Logger.Error("Failed to parse multipart form", zap.Error(err))
		middleware.JsonError(w, http.StatusBadRequest, "Failed to parse multipart form")
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		middleware.JsonError(w, http.StatusBadRequest, "No file uploaded")
		return
	}
	defer file.Close()

	// Content magic byte validation
	isValid, err := ValidateImageFile(file)
	if err != nil || !isValid {
		middleware.JsonError(w, http.StatusBadRequest, "Invalid image format. Allowed formats: JPEG, PNG, GIF, WEBP")
		return
	}

	// Reset seeker
	if seeker, ok := file.(io.ReadSeeker); ok {
		_, _ = seeker.Seek(0, io.SeekStart)
	}

	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "image/jpeg"
	}

	updatedUser, token, err := h.service.UpdateProfilePicture(r.Context(), user.ID, file, header.Filename, contentType)
	if err != nil {
		logger.Logger.Error("UpdateProfilePicture service call failed", zap.String("user_id", user.ID), zap.Error(err))
		middleware.JsonError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Profile picture updated successfully",
		"user":    updatedUser,
		"token":   token,
	})
}
