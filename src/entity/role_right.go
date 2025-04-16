package entity

type RoleRight struct {
	ID        int64  `db:"id"`
	RoleID    int64  `db:"role_id"`
	Section   string `db:"section"`
	Route     string `db:"route"`
	RCreate   bool   `db:"r_create"`
	RRead     bool   `db:"r_read"`
	RUpdate   bool   `db:"r_update"`
	RDelete   bool   `db:"r_delete"`
	CreatedAt string `db:"created_at"`
	UpdatedAt string `db:"updated_at"`
}
