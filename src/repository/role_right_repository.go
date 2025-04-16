package repository

import (
	"context"
	"log"
	"test-tablelink/src/entity"

	"github.com/jmoiron/sqlx"
)

type RoleRightRepository struct {
	db *sqlx.DB
}

func NewRoleRightRepository(db *sqlx.DB) *RoleRightRepository {
	return &RoleRightRepository{db: db}
}

func (r *RoleRightRepository) GetByRoleIDAndRoute(ctx context.Context, roleID int64, section, route string) (*entity.RoleRight, error) {
	var roleRight entity.RoleRight
	query := `
		SELECT id, role_id, section, route, r_create, r_read, r_update, r_delete
		FROM role_rights
		WHERE role_id = $1 AND section = $2 AND route = $3
	`
	err := r.db.GetContext(ctx, &roleRight, query, roleID, section, route)
	if err != nil {
		log.Printf("Error getting role right: %v", err)
		return nil, err
	}
	return &roleRight, nil
}

func (r *RoleRightRepository) CheckPermission(ctx context.Context, roleID int64, section, route, method string) (bool, error) {
	var count int
	query := `
		SELECT COUNT(*)
		FROM role_rights
		WHERE role_id = $1 
		AND section = $2 
		AND route = $3
		AND (
			($4 = 'POST' AND r_create = true) OR
			($4 = 'GET' AND r_read = true) OR
			($4 = 'PUT' AND r_update = true) OR
			($4 = 'DELETE' AND r_delete = true)
		)
	`
	log.Printf("Executing permission check query with params: roleID=%d, section=%s, route=%s, method=%s", roleID, section, route, method)
	err := r.db.GetContext(ctx, &count, query, roleID, section, route, method)
	if err != nil {
		log.Printf("Error executing permission check query: %v", err)
		return false, err
	}
	log.Printf("Permission check result: count=%d", count)
	return count > 0, nil
}

func (r *RoleRightRepository) Create(ctx context.Context, roleRight *entity.RoleRight) error {
	query := `
		INSERT INTO role_rights (role_id, section, route, r_create, r_read, r_update, r_delete)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id`
	return r.db.QueryRowContext(ctx, query,
		roleRight.RoleID,
		roleRight.Section,
		roleRight.Route,
		roleRight.RCreate,
		roleRight.RRead,
		roleRight.RUpdate,
		roleRight.RDelete,
	).Scan(&roleRight.ID)
}

func (r *RoleRightRepository) Update(ctx context.Context, roleRight *entity.RoleRight) error {
	query := `
		UPDATE role_rights 
		SET r_create = $1, r_read = $2, r_update = $3, r_delete = $4, updated_at = CURRENT_TIMESTAMP
		WHERE id = $5`
	_, err := r.db.ExecContext(ctx, query,
		roleRight.RCreate,
		roleRight.RRead,
		roleRight.RUpdate,
		roleRight.RDelete,
		roleRight.ID,
	)
	return err
}

func (r *RoleRightRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM role_rights WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
