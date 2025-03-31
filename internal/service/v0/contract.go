package v0

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/s21platform/staff-service/internal/model"
)

// DbRepo определяет все методы для работы с базой данных
type DbRepo interface {
	// Методы для работы со Staff
	StaffGetByID(ctx context.Context, id uuid.UUID) (*model.Staff, error)
	StaffGetByLogin(ctx context.Context, login string) (*model.Staff, error)
	StaffCreate(ctx context.Context, staff *model.Staff) error
	StaffUpdate(ctx context.Context, staff *model.Staff) error
	StaffDelete(ctx context.Context, id uuid.UUID) error
	StaffList(ctx context.Context, filter *model.StaffFilter) ([]*model.Staff, int, error)

	// Методы для работы с Session
	SessionCreate(ctx context.Context, session *model.Session) error
	SessionGetByToken(ctx context.Context, token string) (*model.Session, error)
	SessionGetByRefreshToken(ctx context.Context, refreshToken string) (*model.Session, error)
	SessionDelete(ctx context.Context, token string) error
	SessionDeleteAllForStaff(ctx context.Context, staffID uuid.UUID) error
	SessionUpdateTokens(ctx context.Context, session *model.Session) error

	// Методы для работы с Role
	RoleGetByID(ctx context.Context, id int) (*model.Role, error)
	RoleList(ctx context.Context) ([]*model.Role, error)
}

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

// Role представляет информацию о роли
type Role struct {
	ID   int
	Name string
}

// Session представляет информацию о сессии
type Session struct {
	ID             uuid.UUID
	StaffID        uuid.UUID
	Token          string
	RefreshToken   string
	ExpiresAt      time.Time
	CreatedAt      time.Time
	LastActivityAt time.Time
}

// StaffFilter представляет параметры фильтрации для списка сотрудников
type StaffFilter struct {
	Page       int
	PageSize   int
	SearchTerm string
	RoleID     int
}
