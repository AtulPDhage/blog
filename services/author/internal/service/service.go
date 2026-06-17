package service

import (
	"context"
	"fmt"
	"io"

	"author/internal/db"
	"author/internal/models"
	"author/internal/rabbitmq"
	"author/internal/s3"
)

// BlogService encapsulates all the core business logic (ownership checks, cloud uploads, cache invalidations).
type BlogService struct {
	repo db.BlogRepository
}

// NewBlogService instantiates a new instance of BlogService.
func NewBlogService(repo db.BlogRepository) *BlogService {
	return &BlogService{repo: repo}
}

// CreateBlog handles image upload to AWS S3, saves details in the database, and invalidates global cache structures.
func (s *BlogService) CreateBlog(ctx context.Context, authorID string, title, description, blogcontent, category string, file io.Reader, filename, contentType string) (*models.Blog, error) {
	// 1. Upload image to S3 bucket securely using the s3 uploader variable
	imageUrl, err := s3.UploadImageFn(ctx, file, filename, contentType)
	if err != nil {
		return nil, fmt.Errorf("failed to upload image: %w", err)
	}

	// 2. Persist record to PostgreSQL
	blog, err := s.repo.Create(ctx, title, description, imageUrl, blogcontent, category, authorID)
	if err != nil {
		return nil, fmt.Errorf("failed to create blog record: %w", err)
	}

	// 3. Trigger cache invalidation via RabbitMQ to update downstream service view caches
	rabbitmq.InvalidateCacheJob([]string{"blogs:*"})

	return blog, nil
}

// UpdateBlog verifies authorship, processes S3 updates if a new file is uploaded, and saves the changes.
func (s *BlogService) UpdateBlog(ctx context.Context, authorID string, blogID int, title, description, blogcontent, category string, file io.Reader, filename, contentType string) (*models.Blog, error) {
	// 1. Fetch existing blog to inspect authorship and default values
	existing, err := s.repo.GetByID(ctx, blogID)
	if err != nil {
		return nil, fmt.Errorf("blog not found: %w", err)
	}

	// 2. Author check (enforces RBAC/least privilege)
	if existing.Author != authorID {
		return nil, fmt.Errorf("unauthorized to update this blog")
	}

	imageUrl := existing.Image

	// 3. Process image upload if a file was provided
	if file != nil {
		uploadedURL, err := s3.UploadImageFn(ctx, file, filename, contentType)
		if err != nil {
			return nil, fmt.Errorf("failed to upload new image: %w", err)
		}
		imageUrl = uploadedURL
	}

	// 4. Fallback to existing values if parameters are empty
	if title == "" {
		title = existing.Title
	}
	if description == "" {
		description = existing.Description
	}
	if blogcontent == "" {
		blogcontent = existing.BlogContent
	}
	if category == "" {
		category = existing.Category
	}

	// 5. Save updates to PostgreSQL
	updated, err := s.repo.Update(ctx, blogID, title, description, imageUrl, blogcontent, category)
	if err != nil {
		return nil, fmt.Errorf("failed to update blog record: %w", err)
	}

	// 6. Invalidate cache structures for this post and overall list views
	rabbitmq.InvalidateCacheJob([]string{"blogs:*", fmt.Sprintf("blog:%d", blogID)})

	return updated, nil
}

// DeleteBlog checks blog ownership, performs a transactional cascading deletion, and deletes the caches.
func (s *BlogService) DeleteBlog(ctx context.Context, authorID string, blogID int) error {
	// 1. Fetch existing blog to check author
	existing, err := s.repo.GetByID(ctx, blogID)
	if err != nil {
		return fmt.Errorf("blog not found: %w", err)
	}

	// 2. Author ownership check
	if existing.Author != authorID {
		return fmt.Errorf("unauthorized to delete this blog")
	}

	// 3. Delete blog post and relational references (comments, saved links) transactionally
	err = s.repo.DeleteTx(ctx, blogID)
	if err != nil {
		return fmt.Errorf("failed to delete blog transactionally: %w", err)
	}

	// 4. Invalidate caches
	rabbitmq.InvalidateCacheJob([]string{"blogs:*", fmt.Sprintf("blog:%d", blogID)})

	return nil
}
