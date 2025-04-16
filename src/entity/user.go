package entity

import "time"

type User struct {
	ID         int64      `db:"id" json:"id"`
	RoleID     int64      `db:"role_id" json:"role_id"`
	RoleName   string     `db:"role_name" json:"role_name"`
	Name       string     `db:"name" json:"name"`
	Email      string     `db:"email" json:"email"`
	Password   string     `db:"password" json:"-"`
	LastAccess *time.Time `db:"last_access" json:"last_access"`
	CreatedAt  time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time  `db:"updated_at" json:"updated_at"`
}
