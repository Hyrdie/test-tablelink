package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"

	"test-tablelink/src/entity"
	"test-tablelink/src/repository"
)

type AuthService struct {
	userRepo  *repository.UserRepository
	redisRepo *repository.RedisRepository
}

func NewAuthService(userRepo *repository.UserRepository, redisRepo *repository.RedisRepository) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		redisRepo: redisRepo,
	}
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Status      bool   `json:"status"`
	Message     string `json:"message"`
	AccessToken string `json:"access_token,omitempty"`
}

func (s *AuthService) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	// Get user by email and password
	user, err := s.userRepo.GetByEmailAndPassword(ctx, req.Email, req.Password)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Generate access token
	token, err := generateToken()
	if err != nil {
		return nil, err
	}

	// Store user data in Redis
	err = s.redisRepo.SetUser(ctx, user)
	if err != nil {
		return nil, err
	}

	// Store token in Redis
	err = s.redisRepo.SetToken(ctx, token, user.ID)
	if err != nil {
		return nil, err
	}

	// Update last access time
	err = s.userRepo.UpdateLastAccess(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		Status:      true,
		Message:     "Successfully",
		AccessToken: token,
	}, nil
}

func (s *AuthService) Logout(ctx context.Context, token string) error {
	// Get user ID from token
	userID, err := s.redisRepo.GetUserIDByToken(ctx, token)
	if err != nil {
		return err
	}

	// Get user by ID
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	// Remove user data from Redis
	err = s.redisRepo.DeleteUser(ctx, user.ID)
	if err != nil {
		return err
	}

	// Remove token from Redis
	return s.redisRepo.DeleteToken(ctx, token)
}

func (s *AuthService) ValidateToken(ctx context.Context, token string) (*entity.User, error) {
	// Get user ID from token
	userID, err := s.redisRepo.GetUserIDByToken(ctx, token)
	if err != nil {
		return nil, err
	}

	// Get user by ID
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Get user data from Redis
	cachedUser, err := s.redisRepo.GetUser(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	return cachedUser, nil
}

func generateToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
