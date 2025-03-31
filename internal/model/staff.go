package model

import (
	"time"

	"github.com/google/uuid"
)

// Staff представляет информацию о сотруднике
type Staff struct {
	ID           uuid.UUID
	Login        string
	PasswordHash string
	RoleID       int
	RoleName     string
	Permissions  Permissions
	CreatedAt    time.Time
	UpdatedAt    time.Time
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
