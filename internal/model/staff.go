package model

import (
	"time"

	"github.com/google/uuid"
)

// Staff представляет информацию о сотруднике
type Staff struct {
	ID           uuid.UUID `db:"id"`
	Login        string    `db:"login"`
	PasswordHash string    `db:"password_hash"`
	RoleID       int       `db:"role_id"`
	RoleName     string    `db:"role_name"`
	Permissions  Permissions `db:"permissions"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

// Permissions представляет разрешения сотрудника
type Permissions struct {
	Access []string
}

// StaffFilter представляет параметры фильтрации для списка сотрудников
type StaffFilter struct {
	Page       int
	PageSize   int
	SearchTerm string
	RoleID     int
}
