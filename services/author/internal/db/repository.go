package db

import (
	"context"
	"fmt"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"

	"author/internal/models"
)

// BlogRepository defines the contract for all database interactions involving blogs, comments, and savedblogs.
type BlogRepository interface {
	// Create inserts a new blog post into the database.
	Create(ctx context.Context, title, description, image, blogcontent, category, author string) (*models.Blog, error)
	
	// GetByID retrieves a single blog post using its unique primary key.
	GetByID(ctx context.Context, id int) (*models.Blog, error)
	
	// Update modifies an existing blog post record in the database.
	Update(ctx context.Context, id int, title, description, image, blogcontent, category string) (*models.Blog, error)
	
	// DeleteTx removes a blog post and cascade-deletes comments/savedblogs references in a SQL transaction.
	DeleteTx(ctx context.Context, id int) error
}

// postgresBlogRepository is a concrete implementation of BlogRepository querying PostgreSQL using pgxpool.
type postgresBlogRepository struct {
	pool *pgxpool.Pool
}

// NewPostgresBlogRepository instantiates a new postgres implementation of BlogRepository.
func NewPostgresBlogRepository(pool *pgxpool.Pool) BlogRepository {
	return &postgresBlogRepository{pool: pool}
}

// Create inserts a new record into the blogs table and returns the parsed Blog model.
func (r *postgresBlogRepository) Create(ctx context.Context, title, description, image, blogcontent, category, author string) (*models.Blog, error) {
	var blog models.Blog
	err := r.pool.QueryRow(ctx, InsertBlogQuery, title, description, image, blogcontent, category, author).Scan(
		&blog.ID, &blog.Title, &blog.Description, &blog.Image, &blog.BlogContent, &blog.Category, &blog.Author, &blog.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("repository create failed: %w", err)
	}
	return &blog, nil
}

// GetByID queries the blogs table using the blog primary key ID.
func (r *postgresBlogRepository) GetByID(ctx context.Context, id int) (*models.Blog, error) {
	var blog models.Blog
	err := r.pool.QueryRow(ctx, SelectBlogByIDQuery, id).Scan(
		&blog.ID, &blog.Title, &blog.Description, &blog.Image, &blog.BlogContent, &blog.Category, &blog.Author, &blog.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("repository get by id failed: %w", err)
	}
	return &blog, nil
}

// Update executes an UPDATE statement against the blogs table modifying all mutable fields.
func (r *postgresBlogRepository) Update(ctx context.Context, id int, title, description, image, blogcontent, category string) (*models.Blog, error) {
	var blog models.Blog
	err := r.pool.QueryRow(ctx, UpdateBlogQuery, title, description, image, blogcontent, category, id).Scan(
		&blog.ID, &blog.Title, &blog.Description, &blog.Image, &blog.BlogContent, &blog.Category, &blog.Author, &blog.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("repository update failed: %w", err)
	}
	return &blog, nil
}

// DeleteTx begins a SQL transaction to delete a blog and clean up comments and savedblogs records cascade-style.
func (r *postgresBlogRepository) DeleteTx(ctx context.Context, id int) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("repository begin transaction failed: %w", err)
	}
	defer func() {
		// Rollbacks if commit has not happened, preventing partial data updates
		_ = tx.Rollback(ctx)
	}()

	// 1. Delete main blog post record
	_, err = tx.Exec(ctx, DeleteBlogByIDQuery, id)
	if err != nil {
		return fmt.Errorf("repository delete blog query failed: %w", err)
	}

	// 2. Clean up child references in comments
	_, err = tx.Exec(ctx, DeleteCommentsByBlogIDQuery, strconv.Itoa(id))
	if err != nil {
		return fmt.Errorf("repository delete comments query failed: %w", err)
	}

	// 3. Clean up child references in savedblogs
	_, err = tx.Exec(ctx, DeleteSavedBlogsByBlogIDQuery, strconv.Itoa(id))
	if err != nil {
		return fmt.Errorf("repository delete saved blogs query failed: %w", err)
	}

	// Commit transaction if all steps succeed
	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("repository commit transaction failed: %w", err)
	}

	return nil
}
