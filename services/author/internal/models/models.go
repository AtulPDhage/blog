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

// AITitleRequest payload schema
type AITitleRequest struct {
	Text string `json:"text"`
}

// AIDescriptionRequest payload schema
type AIDescriptionRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

// AIBlogRequest payload schema
type AIBlogRequest struct {
	Blog string `json:"blog"`
}

// AIBlogResponse response schema
type AIBlogResponse struct {
	HTML string `json:"html"`
}
