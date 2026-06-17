package models

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Blog represents the PostgreSQL blog schema structure
type Blog struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	BlogContent string    `json:"blogcontent"`
	Image       string    `json:"image"`
	Category    string    `json:"category"`
	Author      string    `json:"author"`
	CreatedAt   time.Time `json:"created_at"`
}

// Comment represents the PostgreSQL comments schema structure
type Comment struct {
	ID        int       `json:"id"`
	Comment   string    `json:"comment"`
	UserID    string    `json:"userid"`
	Username  string    `json:"username"`
	BlogID    string    `json:"blogid"`
	CreatedAt time.Time `json:"created_at"`
}

// SavedBlog represents the PostgreSQL savedblogs schema structure
type SavedBlog struct {
	ID        int       `json:"id"`
	UserID    string    `json:"userid"`
	BlogID    string    `json:"blogid"`
	CreatedAt time.Time `json:"created_at"`
}

// User represents the decoded JWT payload details
type User struct {
	ID        string `json:"_id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Image     string `json:"image"`
	Instagram string `json:"instagram"`
	Linkedin  string `json:"linkedin"`
	Facebook  string `json:"facebook"`
	Bio       string `json:"bio"`
}

// Claims represents the JWT custom claims payload
type Claims struct {
	User User `json:"user"`
	jwt.RegisteredClaims
}
