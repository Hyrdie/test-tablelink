package repository

import (
	"context"
	"test-tablelink/src/entity"

	"github.com/jmoiron/sqlx"
)

type RoleRepository struct {
	db *sqlx.DB
}

func NewRoleRepository(db *sqlx.DB) *RoleRepository {
	return &RoleRepository{db: db}
}

func (r *RoleRepository) GetByID(ctx context.Context, id int) (*entity.Role, error) {
	var role entity.Role
	query := `SELECT id, name, created_at, updated_at FROM roles WHERE id = $1`
	err := r.db.GetContext(ctx, &role, query, id)
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *RoleRepository) GetByName(ctx context.Context, name string) (*entity.Role, error) {
	var role entity.Role
	query := `SELECT id, name, created_at, updated_at FROM roles WHERE name = $1`
	err := r.db.GetContext(ctx, &role, query, name)
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *RoleRepository) Create(ctx context.Context, role *entity.Role) error {
	query := `INSERT INTO roles (name) VALUES ($1) RETURNING id`
	return r.db.QueryRowContext(ctx, query, role.Name).Scan(&role.ID)
}

func (r *RoleRepository) Update(ctx context.Context, role *entity.Role) error {
	query := `UPDATE roles SET name = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, role.Name, role.ID)
	return err
}

func (r *RoleRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM roles WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
