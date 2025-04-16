package repository

import (
	"context"
	"database/sql"
	"test-tablelink/src/entity"
)

type UserRepository interface {
	GetAll(ctx context.Context) ([]*entity.User, error)
	GetByID(ctx context.Context, id int) (*entity.User, error)
	Create(ctx context.Context, user *entity.User) error
	Update(ctx context.Context, user *entity.User) error
	Delete(ctx context.Context, id int) error
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetAll(ctx context.Context) ([]*entity.User, error) {
	return nil, nil
}

func (r *userRepository) GetByID(ctx context.Context, id int) (*entity.User, error) {
	return nil, nil
}

func (r *userRepository) Create(ctx context.Context, user *entity.User) error {
	return nil
}

func (r *userRepository) Update(ctx context.Context, user *entity.User) error {
	return nil
}

func (r *userRepository) Delete(ctx context.Context, id int) error {
	return nil
}
