package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"blog/internal/models"
)

type BlogRepository interface {
	GetAllBlogs(ctx context.Context, searchQuery, category string, limit, offset int) ([]models.Blog, error)
	GetSingleBlog(ctx context.Context, id int) (*models.Blog, error)
	AddComment(ctx context.Context, comment string, blogID string, userID, username string) error
	GetAllComments(ctx context.Context, blogID string) ([]models.Comment, error)
	GetCommentByID(ctx context.Context, commentID int) (*models.Comment, error)
	DeleteComment(ctx context.Context, commentID int) error
	GetSavedBlog(ctx context.Context, userID string, blogID string) (*models.SavedBlog, error)
	SaveBlog(ctx context.Context, userID string, blogID string) (bool, error) // Returns true if saved, false if removed
	GetSavedBlogs(ctx context.Context, userID string) ([]models.Blog, error)
	
	GetLikedBlog(ctx context.Context, userID string, blogID string) (*models.LikedBlog, error)
	LikeBlog(ctx context.Context, userID string, blogID string) (bool, error) // Returns true if liked, false if unliked
	GetBlogLikesCount(ctx context.Context, blogID string) (int, error)
	IsBlogLikedByUser(ctx context.Context, userID string, blogID string) (bool, error)
	IncrementBlogViews(ctx context.Context, id int) error
}

type PostgresBlogRepository struct{}

func NewPostgresBlogRepository() BlogRepository {
	return &PostgresBlogRepository{}
}

// GetAllBlogs retrieves all blogs matching optional searchQuery and category, with pagination
func (r *PostgresBlogRepository) GetAllBlogs(ctx context.Context, searchQuery, category string, limit, offset int) ([]models.Blog, error) {
	var rows pgx.Rows
	var err error

	if searchQuery != "" && category != "" {
		wildcardSearch := "%" + searchQuery + "%"
		rows, err = Pool.Query(ctx, SelectBlogsBySearchAndCategoryQuery, wildcardSearch, category, limit, offset)
	} else if searchQuery != "" {
		wildcardSearch := "%" + searchQuery + "%"
		rows, err = Pool.Query(ctx, SelectBlogsBySearchQuery, wildcardSearch, limit, offset)
	} else if category != "" {
		rows, err = Pool.Query(ctx, SelectBlogsByCategoryQuery, category, limit, offset)
	} else {
		rows, err = Pool.Query(ctx, SelectAllBlogsQuery, limit, offset)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to query blogs: %w", err)
	}
	defer rows.Close()

	var blogs []models.Blog
	for rows.Next() {
		var b models.Blog
		err := rows.Scan(
			&b.ID,
			&b.Title,
			&b.Description,
			&b.BlogContent,
			&b.Image,
			&b.Category,
			&b.Author,
			&b.CreatedAt,
			&b.Views,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan blog row: %w", err)
		}
		blogs = append(blogs, b)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during row iteration: %w", err)
	}

	// Ensure we return an empty slice instead of nil for clean JSON responses
	if blogs == nil {
		blogs = []models.Blog{}
	}

	return blogs, nil
}

// GetSingleBlog retrieves a single blog by ID
func (r *PostgresBlogRepository) GetSingleBlog(ctx context.Context, id int) (*models.Blog, error) {
	var b models.Blog
	err := Pool.QueryRow(ctx, SelectBlogByIDQuery, id).Scan(
		&b.ID,
		&b.Title,
		&b.Description,
		&b.BlogContent,
		&b.Image,
		&b.Category,
		&b.Author,
		&b.CreatedAt,
		&b.Views,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // Return nil when not found
		}
		return nil, fmt.Errorf("failed to query single blog: %w", err)
	}
	return &b, nil
}

// AddComment inserts a new comment
func (r *PostgresBlogRepository) AddComment(ctx context.Context, comment string, blogID string, userID, username string) error {
	var c models.Comment
	err := Pool.QueryRow(ctx, InsertCommentQuery, comment, blogID, userID, username).Scan(
		&c.ID,
		&c.Comment,
		&c.UserID,
		&c.Username,
		&c.BlogID,
		&c.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to insert comment: %w", err)
	}
	return nil
}

// GetAllComments retrieves all comments for a blog ID
func (r *PostgresBlogRepository) GetAllComments(ctx context.Context, blogID string) ([]models.Comment, error) {
	rows, err := Pool.Query(ctx, SelectCommentsByBlogIDQuery, blogID)
	if err != nil {
		return nil, fmt.Errorf("failed to query comments: %w", err)
	}
	defer rows.Close()

	var comments []models.Comment
	for rows.Next() {
		var c models.Comment
		err := rows.Scan(
			&c.ID,
			&c.Comment,
			&c.UserID,
			&c.Username,
			&c.BlogID,
			&c.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan comment row: %w", err)
		}
		comments = append(comments, c)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during comment row iteration: %w", err)
	}

	// Ensure empty slice for clean JSON
	if comments == nil {
		comments = []models.Comment{}
	}

	return comments, nil
}

// GetCommentByID retrieves a single comment by its ID
func (r *PostgresBlogRepository) GetCommentByID(ctx context.Context, commentID int) (*models.Comment, error) {
	var c models.Comment
	err := Pool.QueryRow(ctx, SelectCommentByIDQuery, commentID).Scan(
		&c.ID,
		&c.Comment,
		&c.UserID,
		&c.Username,
		&c.BlogID,
		&c.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to query comment by ID: %w", err)
	}
	return &c, nil
}

// DeleteComment deletes a comment by ID
func (r *PostgresBlogRepository) DeleteComment(ctx context.Context, commentID int) error {
	commandTag, err := Pool.Exec(ctx, DeleteCommentByIDQuery, commentID)
	if err != nil {
		return fmt.Errorf("failed to delete comment: %w", err)
	}
	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("comment not found or not deleted")
	}
	return nil
}

// GetSavedBlog checks if a blog has been saved by a user
func (r *PostgresBlogRepository) GetSavedBlog(ctx context.Context, userID string, blogID string) (*models.SavedBlog, error) {
	var sb models.SavedBlog
	err := Pool.QueryRow(ctx, SelectSavedBlogQuery, userID, blogID).Scan(
		&sb.ID,
		&sb.UserID,
		&sb.BlogID,
		&sb.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to query saved blog: %w", err)
	}
	return &sb, nil
}

// SaveBlog saves/unsaves a blog (toggles state) and returns true if saved, false if removed
func (r *PostgresBlogRepository) SaveBlog(ctx context.Context, userID string, blogID string) (bool, error) {
	existing, err := r.GetSavedBlog(ctx, userID, blogID)
	if err != nil {
		return false, err
	}

	if existing == nil {
		// Insert
		var sb models.SavedBlog
		err := Pool.QueryRow(ctx, InsertSavedBlogQuery, userID, blogID).Scan(
			&sb.ID,
			&sb.UserID,
			&sb.BlogID,
			&sb.CreatedAt,
		)
		if err != nil {
			return false, fmt.Errorf("failed to insert saved blog: %w", err)
		}
		return true, nil
	} else {
		// Delete
		_, err := Pool.Exec(ctx, DeleteSavedBlogQuery, userID, blogID)
		if err != nil {
			return false, fmt.Errorf("failed to delete saved blog: %w", err)
		}
		return false, nil
	}
}

// GetSavedBlogs retrieves all saved blogs for a user with full details
func (r *PostgresBlogRepository) GetSavedBlogs(ctx context.Context, userID string) ([]models.Blog, error) {
	rows, err := Pool.Query(ctx, SelectSavedBlogsDetailedByUserIDQuery, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query saved blogs: %w", err)
	}
	defer rows.Close()

	var blogs []models.Blog
	for rows.Next() {
		var b models.Blog
		err := rows.Scan(
			&b.ID,
			&b.Title,
			&b.Description,
			&b.BlogContent,
			&b.Image,
			&b.Category,
			&b.Author,
			&b.CreatedAt,
			&b.Views,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan saved blog row: %w", err)
		}
		blogs = append(blogs, b)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during saved blogs row iteration: %w", err)
	}

	if blogs == nil {
		blogs = []models.Blog{}
	}

	return blogs, nil
}

// GetLikedBlog checks if a user has liked a blog
func (r *PostgresBlogRepository) GetLikedBlog(ctx context.Context, userID string, blogID string) (*models.LikedBlog, error) {
	var lb models.LikedBlog
	err := Pool.QueryRow(ctx, SelectLikedBlogQuery, userID, blogID).Scan(
		&lb.ID,
		&lb.UserID,
		&lb.BlogID,
		&lb.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to query liked blog: %w", err)
	}
	return &lb, nil
}

// LikeBlog toggles the liked status of a blog by a user
func (r *PostgresBlogRepository) LikeBlog(ctx context.Context, userID string, blogID string) (bool, error) {
	existing, err := r.GetLikedBlog(ctx, userID, blogID)
	if err != nil {
		return false, err
	}

	if existing == nil {
		var lb models.LikedBlog
		err := Pool.QueryRow(ctx, InsertLikedBlogQuery, userID, blogID).Scan(
			&lb.ID,
			&lb.UserID,
			&lb.BlogID,
			&lb.CreatedAt,
		)
		if err != nil {
			return false, fmt.Errorf("failed to insert liked blog: %w", err)
		}
		return true, nil
	} else {
		_, err := Pool.Exec(ctx, DeleteLikedBlogQuery, userID, blogID)
		if err != nil {
			return false, fmt.Errorf("failed to delete liked blog: %w", err)
		}
		return false, nil
	}
}

// GetBlogLikesCount returns the total number of likes for a blog post
func (r *PostgresBlogRepository) GetBlogLikesCount(ctx context.Context, blogID string) (int, error) {
	var count int
	err := Pool.QueryRow(ctx, SelectLikesCountByBlogIDQuery, blogID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to query likes count: %w", err)
	}
	return count, nil
}

// IsBlogLikedByUser returns true if a user has liked a blog post
func (r *PostgresBlogRepository) IsBlogLikedByUser(ctx context.Context, userID string, blogID string) (bool, error) {
	existing, err := r.GetLikedBlog(ctx, userID, blogID)
	if err != nil {
		return false, err
	}
	return existing != nil, nil
}

// IncrementBlogViews increments the view count for a blog post
func (r *PostgresBlogRepository) IncrementBlogViews(ctx context.Context, id int) error {
	_, err := Pool.Exec(ctx, IncrementViewsQuery, id)
	if err != nil {
		return fmt.Errorf("failed to increment blog views: %w", err)
	}
	return nil
}
