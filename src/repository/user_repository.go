package repository

import (
	"context"
	"test-tablelink/src/entity"

	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetByID(ctx context.Context, id int64) (*entity.User, error) {
	var user entity.User
	query := `
		SELECT u.id, u.role_id, u.name, u.email, u.password, u.last_access, u.created_at, u.updated_at,
		       r.name as role_name
		FROM users u
		LEFT JOIN roles r ON u.role_id = r.id
		WHERE u.id = $1`
	err := r.db.GetContext(ctx, &user, query, id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User
	query := `
		SELECT u.id, u.role_id, u.name, u.email, u.password, u.last_access, u.created_at, u.updated_at,
		       r.name as role_name
		FROM users u
		LEFT JOIN roles r ON u.role_id = r.id
		WHERE u.email = $1`
	err := r.db.GetContext(ctx, &user, query, email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByEmailAndPassword(ctx context.Context, email, password string) (*entity.User, error) {
	var user entity.User
	query := `
		SELECT u.id, u.role_id, u.name, u.email, u.password, u.last_access, u.created_at, u.updated_at,
		       r.name as role_name
		FROM users u
		LEFT JOIN roles r ON u.role_id = r.id
		WHERE u.email = $1 AND u.password = $2`
	err := r.db.GetContext(ctx, &user, query, email, password)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Create(ctx context.Context, user *entity.User) error {
	query := `
		INSERT INTO users (role_id, name, email, password)
		VALUES ($1, $2, $3, $4)
		RETURNING id`
	return r.db.QueryRowContext(ctx, query,
		user.RoleID,
		user.Name,
		user.Email,
		user.Password,
	).Scan(&user.ID)
}

func (r *UserRepository) Update(ctx context.Context, user *entity.User) error {
	query := `
		UPDATE users 
		SET role_id = $1, name = $2, email = $3, password = $4, updated_at = CURRENT_TIMESTAMP
		WHERE id = $5`
	_, err := r.db.ExecContext(ctx, query,
		user.RoleID,
		user.Name,
		user.Email,
		user.Password,
		user.ID,
	)
	return err
}

func (r *UserRepository) UpdateLastAccess(ctx context.Context, id int64) error {
	query := `UPDATE users SET last_access = CURRENT_TIMESTAMP WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *UserRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *UserRepository) GetAll(ctx context.Context) ([]*entity.User, error) {
	var users []*entity.User
	query := `
		SELECT u.id, u.role_id, u.name, u.email, u.last_access, u.created_at, u.updated_at,
		       r.name as role_name
		FROM users u
		LEFT JOIN roles r ON u.role_id = r.id`
	err := r.db.SelectContext(ctx, &users, query)
	if err != nil {
		return nil, err
	}
	return users, nil
}
