package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"user/internal/db"
	"user/internal/google"
	"user/internal/models"
	"user/internal/s3"
)

type UserService struct {
	repo         db.UserRepository
	googleClient *google.GoogleClient
	jwtSecret    string
}

// NewUserService instantiates a new UserService layer
func NewUserService(repo db.UserRepository, googleClient *google.GoogleClient, jwtSecret string) *UserService {
	return &UserService{
		repo:         repo,
		googleClient: googleClient,
		jwtSecret:    jwtSecret,
	}
}

// SetGoogleClient overrides the Google API client (useful for unit testing)
func (s *UserService) SetGoogleClient(gc *google.GoogleClient) {
	s.googleClient = gc
}

// GenerateToken signs a new JWT with user payload details expiring in 5 days
func (s *UserService) GenerateToken(user *models.User) (string, error) {
	claims := &models.Claims{
		User: *user,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(5 * 24 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}
	return tokenString, nil
}

// LoginGoogle processes Google authorization code, exchanges it, and handles login/creation
func (s *UserService) LoginGoogle(ctx context.Context, code string) (*models.User, string, error) {
	if code == "" {
		return nil, "", errors.New("authorization code is required")
	}

	tokenResp, err := s.googleClient.ExchangeCodeForToken(ctx, code)
	if err != nil {
		return nil, "", fmt.Errorf("google authentication failed: %w", err)
	}

	userInfo, err := s.googleClient.GetUserInfo(ctx, tokenResp.AccessToken)
	if err != nil {
		return nil, "", fmt.Errorf("failed to fetch user info: %w", err)
	}

	user, err := s.repo.FindByEmail(ctx, userInfo.Email)
	if err != nil {
		return nil, "", err
	}

	if user == nil {
		user, err = s.repo.Create(ctx, userInfo.Name, userInfo.Email, userInfo.Picture)
		if err != nil {
			return nil, "", err
		}
	}

	token, err := s.GenerateToken(user)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

// GetUserProfile retrieves a user profile by string ID
func (s *UserService) GetUserProfile(ctx context.Context, id string) (*models.User, error) {
	if id == "" {
		return nil, errors.New("user ID is required")
	}
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// UpdateUser modifies user name and bio/social details, returning user and new token
func (s *UserService) UpdateUser(ctx context.Context, id string, name, instagram, linkedin, facebook, bio string) (*models.User, string, error) {
	if id == "" {
		return nil, "", errors.New("user ID is required")
	}

	user, err := s.repo.Update(ctx, id, name, instagram, linkedin, facebook, bio)
	if err != nil {
		return nil, "", err
	}

	token, err := s.GenerateToken(user)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

// UpdateProfilePicture uploads a picture to AWS S3 and updates user profile, returning user and new token
func (s *UserService) UpdateProfilePicture(ctx context.Context, id string, file io.Reader, filename, contentType string) (*models.User, string, error) {
	if id == "" {
		return nil, "", errors.New("user ID is required")
	}

	imageUrl, err := s3.UploadImageFn(ctx, file, filename, contentType)
	if err != nil {
		return nil, "", fmt.Errorf("s3 upload failed: %w", err)
	}

	user, err := s.repo.UpdateProfilePicture(ctx, id, imageUrl)
	if err != nil {
		return nil, "", err
	}

	token, err := s.GenerateToken(user)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}
