package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"author/internal/db"
	"author/internal/gemini"
	"author/internal/logger"
	"author/internal/middleware"
	"author/internal/models"
	"author/internal/service"
)

type BlogHandler struct {
	service *service.BlogService
}

func NewBlogHandler(s *service.BlogService) *BlogHandler {
	return &BlogHandler{service: s}
}

// ValidateImageFile reads the first 512 bytes of a reader to check if the file is an allowed image type
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

	// If the file implements io.Seeker, reset the read offset so we can upload the full file
	if seeker, ok := file.(io.ReadSeeker); ok {
		_, _ = seeker.Seek(0, io.SeekStart)
	}

	return allowedTypes[contentType], nil
}

// CreateBlog creates a new blog post
func (h *BlogHandler) CreateBlog(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		middleware.JsonError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Parse multipart form (max 10MB in memory)
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		logger.Logger.Error("Failed to parse multipart form", zap.Error(err))
		middleware.JsonError(w, http.StatusBadRequest, "Failed to parse multipart form")
		return
	}

	title := r.FormValue("title")
	description := r.FormValue("description")
	blogcontent := r.FormValue("blogcontent")
	category := r.FormValue("category")

	file, header, err := r.FormFile("file")
	if err != nil {
		middleware.JsonError(w, http.StatusBadRequest, "No file uploaded")
		return
	}
	defer file.Close()

	// Validate content type of uploaded file to prevent arbitrary file upload
	isValid, err := ValidateImageFile(file)
	if err != nil || !isValid {
		middleware.JsonError(w, http.StatusBadRequest, "Invalid image format. Allowed formats: JPEG, PNG, GIF, WEBP")
		return
	}

	// Reset seeker before uploading if it's seekable
	if seeker, ok := file.(io.ReadSeeker); ok {
		_, _ = seeker.Seek(0, io.SeekStart)
	}

	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "image/jpeg"
	}

	// Delegate logic to Service Layer
	blog, err := h.service.CreateBlog(r.Context(), user.ID, title, description, blogcontent, category, file, header.Filename, contentType)
	if err != nil {
		logger.Logger.Error("Create blog service call failed", zap.Error(err))
		middleware.JsonError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Blog created successfully",
		"blog":    blog,
	})
}

// UpdateBlog updates an existing blog post
func (h *BlogHandler) UpdateBlog(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		middleware.JsonError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Using go-chi parameters parsing
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		middleware.JsonError(w, http.StatusBadRequest, "Invalid blog ID")
		return
	}

	// Parse multipart form
	err = r.ParseMultipartForm(10 << 20)
	if err != nil {
		logger.Logger.Error("Failed to parse multipart form for update", zap.Error(err))
		middleware.JsonError(w, http.StatusBadRequest, "Failed to parse multipart form")
		return
	}

	title := r.FormValue("title")
	description := r.FormValue("description")
	blogcontent := r.FormValue("blogcontent")
	category := r.FormValue("category")

	var file io.Reader
	var filename string
	var contentType string

	// Handle optional new file upload
	uploadedFile, header, err := r.FormFile("file")
	if err == nil {
		defer uploadedFile.Close()
		// Validate content type of uploaded file
		isValid, err := ValidateImageFile(uploadedFile)
		if err != nil || !isValid {
			middleware.JsonError(w, http.StatusBadRequest, "Invalid image format. Allowed formats: JPEG, PNG, GIF, WEBP")
			return
		}

		if seeker, ok := uploadedFile.(io.ReadSeeker); ok {
			_, _ = seeker.Seek(0, io.SeekStart)
		}

		file = uploadedFile
		filename = header.Filename
		contentType = header.Header.Get("Content-Type")
		if contentType == "" {
			contentType = "image/jpeg"
		}
	}

	// Delegate logic to Service Layer
	updated, err := h.service.UpdateBlog(r.Context(), user.ID, id, title, description, blogcontent, category, file, filename, contentType)
	if err != nil {
		logger.Logger.Error("Update blog service call failed", zap.Int("blog_id", id), zap.Error(err))
		if strings.Contains(err.Error(), "unauthorized") {
			middleware.JsonError(w, http.StatusForbidden, "You are not authorized to update this blog")
		} else if strings.Contains(err.Error(), "not found") {
			middleware.JsonError(w, http.StatusNotFound, "Blog not found")
		} else {
			middleware.JsonError(w, http.StatusInternalServerError, "Internal Server Error")
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Blog updated successfully",
		"blog":    updated,
	})
}

// DeleteBlog deletes an existing blog post
func (h *BlogHandler) DeleteBlog(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		middleware.JsonError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		middleware.JsonError(w, http.StatusBadRequest, "Invalid blog ID")
		return
	}

	// Delegate logic to Service Layer
	err = h.service.DeleteBlog(r.Context(), user.ID, id)
	if err != nil {
		logger.Logger.Error("Delete blog service call failed", zap.Int("blog_id", id), zap.Error(err))
		if strings.Contains(err.Error(), "unauthorized") {
			middleware.JsonError(w, http.StatusForbidden, "You are not authorized to delete this blog")
		} else if strings.Contains(err.Error(), "not found") {
			middleware.JsonError(w, http.StatusNotFound, "Blog not found")
		} else {
			middleware.JsonError(w, http.StatusInternalServerError, "Internal Server Error")
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"message": "Blog deleted successfully",
	})
}

// CleanAISpanHelper cleans Markdown formatting symbols from Gemini responses
func CleanAISpanHelper(text string) string {
	// Replaces "**", "\r", "\n", "*", "_", "`", "~"
	re := regexp.MustCompile(`\*\*|[\r\n]+|[*_\x60~]`)
	cleaned := re.ReplaceAllString(text, "")
	return strings.TrimSpace(cleaned)
}

// AITitle corrects the grammar of a blog title
func (h *BlogHandler) AITitle(w http.ResponseWriter, r *http.Request) {
	var req models.AITitleRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.Text == "" {
		middleware.JsonError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	prompt := fmt.Sprintf(`Correct the grammar of the following blog title and return only the corrected title without any additional text, formatting, or symbols: "%s"`, req.Text)

	// Fetch Gemini API Key
	apiKey := r.Context().Value("GeminiAPIKey").(string)

	resultText, err := gemini.CallGemini(r.Context(), apiKey, "gemini-2.5-flash", prompt)
	if err != nil {
		logger.Logger.Error("AI title grammar check failed", zap.Error(err))
		middleware.JsonError(w, http.StatusBadRequest, "Something went wrong!")
		return
	}

	cleaned := CleanAISpanHelper(resultText)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(cleaned)
}

// AIDescription corrects or generates a short description for the blog
func (h *BlogHandler) AIDescription(w http.ResponseWriter, r *http.Request) {
	var req models.AIDescriptionRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.Title == "" {
		middleware.JsonError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	var prompt string
	if req.Description == "" {
		prompt = fmt.Sprintf(`Generate only one short blog description based on this title: "%s". Your response must be only one sentence, strictly under 30 words, with no options, no greetings, and no extra text. Do not explain. Do not say 'here is'. Just return the description only.`, req.Title)
	} else {
		prompt = fmt.Sprintf(`Fix the grammar in the following blog description and return only the corrected sentence. Do not add anything else: "%s"`, req.Description)
	}

	apiKey := r.Context().Value("GeminiAPIKey").(string)

	resultText, err := gemini.CallGemini(r.Context(), apiKey, "gemini-2.5-flash", prompt)
	if err != nil {
		logger.Logger.Error("AI description operation failed", zap.Error(err))
		middleware.JsonError(w, http.StatusBadRequest, "Something went wrong!")
		return
	}

	cleaned := CleanAISpanHelper(resultText)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(cleaned)
}

// AIBlog corrects grammar in rich HTML content
func (h *BlogHandler) AIBlog(w http.ResponseWriter, r *http.Request) {
	var req models.AIBlogRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.Blog == "" {
		middleware.JsonError(w, http.StatusBadRequest, "Please provide blog")
		return
	}

	// Refactored to use AIBlogSystemPrompt from queries constants
	fullMessage := fmt.Sprintf("%s\n\n%s", db.AIBlogSystemPrompt, req.Blog)

	apiKey := r.Context().Value("GeminiAPIKey").(string)

	resultText, err := gemini.CallGemini(r.Context(), apiKey, "gemini-2.5-flash", fullMessage)
	if err != nil {
		logger.Logger.Error("AI blog grammar correction failed", zap.Error(err))
		middleware.JsonError(w, http.StatusBadRequest, "Something went wrong!")
		return
	}

	// Clean markdown wrappers if returned by Gemini (e.g. ```html ... ```)
	cleaned := resultText
	cleaned = strings.ReplaceAll(cleaned, "```html", "")
	cleaned = strings.ReplaceAll(cleaned, "```", "")
	cleaned = strings.TrimSpace(cleaned)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(models.AIBlogResponse{HTML: cleaned})
}
