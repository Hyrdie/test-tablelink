package service

import (
	"context"
	"errors"

	"test-tablelink/src/entity"
	"test-tablelink/src/repository"
	"test-tablelink/src/v1/contract"
)

type UserService struct {
	userRepo  *repository.UserRepository
	redisRepo *repository.RedisRepository
}

func NewUserService(userRepo *repository.UserRepository, redisRepo *repository.RedisRepository) *UserService {
	return &UserService{
		userRepo:  userRepo,
		redisRepo: redisRepo,
	}
}

type CreateUserRequest struct {
	RoleID   int64  `json:"role_id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateUserRequest struct {
	Name string `json:"name"`
}

type UserResponse struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func (s *UserService) GetAllUsers(ctx context.Context) (*UserResponse, error) {
	users, err := s.userRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return &UserResponse{
		Status:  true,
		Message: "Successfully",
		Data:    users,
	}, nil
}

func (s *UserService) GetUser(ctx context.Context, id int64) (*contract.UserResponse, error) {
	// Try to get from Redis first
	user, err := s.redisRepo.GetUser(ctx, id)
	if err == nil {
		return &contract.UserResponse{
			Status:  true,
			Message: "Successfully",
			Data:    user,
		}, nil
	}

	// If not in Redis, get from database
	user, err = s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Cache the user in Redis
	if err := s.redisRepo.SetUser(ctx, user); err != nil {
		// Log error but don't fail the request
	}

	return &contract.UserResponse{
		Status:  true,
		Message: "Successfully",
		Data:    user,
	}, nil
}

func (s *UserService) CreateUser(ctx context.Context, req *CreateUserRequest) (*UserResponse, error) {
	user := &entity.User{
		RoleID:   req.RoleID,
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	err := s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	// Cache the new user in Redis
	if err := s.redisRepo.SetUser(ctx, user); err != nil {
		// Log error but don't fail the request
	}

	return &UserResponse{
		Status:  true,
		Message: "Successfully",
	}, nil
}

func (s *UserService) UpdateUser(ctx context.Context, req *UpdateUserRequest) (*UserResponse, error) {
	// Get user from context (set by auth middleware)
	user, ok := ctx.Value("user").(*entity.User)
	if !ok {
		return nil, errors.New("user not found in context")
	}

	// Update user name
	user.Name = req.Name

	err := s.userRepo.Update(ctx, user)
	if err != nil {
		return nil, err
	}

	// Update cache
	if err := s.redisRepo.SetUser(ctx, user); err != nil {
		// Log error but don't fail the request
	}

	return &UserResponse{
		Status:  true,
		Message: "Successfully",
	}, nil
}

func (s *UserService) DeleteUser(ctx context.Context, userID int64) (*UserResponse, error) {
	err := s.userRepo.Delete(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Remove from cache
	if err := s.redisRepo.DeleteUser(ctx, userID); err != nil {
		// Log error but don't fail the request
	}

	return &UserResponse{
		Status:  true,
		Message: "Successfully",
	}, nil
}

func (s *UserService) Login(ctx context.Context, email, password string) (*contract.UserResponse, error) {
	// Get user by email and password
	user, err := s.userRepo.GetByEmailAndPassword(ctx, email, password)
	if err != nil {
		return nil, err
	}

	// Cache the user in Redis
	if err := s.redisRepo.SetUser(ctx, user); err != nil {
		// Log error but don't fail the request
	}

	return &contract.UserResponse{
		Status:  true,
		Message: "Successfully",
		Data:    user,
	}, nil
}

func (s *UserService) Logout(ctx context.Context, id int64) (*contract.UserResponse, error) {
	// Remove from cache
	if err := s.redisRepo.DeleteUser(ctx, id); err != nil {
		// Log error but don't fail the request
	}

	return &contract.UserResponse{
		Status:  true,
		Message: "Successfully",
	}, nil
}
