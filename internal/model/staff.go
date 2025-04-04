package model

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Staff представляет информацию о сотруднике
type Staff struct {
	ID           uuid.UUID   `db:"id"`
	Login        string      `db:"login"`
	PasswordHash string      `db:"password_hash"`
	RoleID       int         `db:"role_id"`
	RoleName     string      `db:"role_name"`
	Permissions  Permissions `db:"permissions"`
	CreatedAt    time.Time   `db:"created_at"`
	UpdatedAt    time.Time   `db:"updated_at"`
}

// Permissions представляет разрешения сотрудника
type Permissions struct {
	Access []string `json:"access"`
}

// Scan реализует интерфейс sql.Scanner для Permissions
func (p *Permissions) Scan(value interface{}) error {
	if value == nil {
		p.Access = []string{}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		p.Access = []string{}
		return nil
	}

	return json.Unmarshal(bytes, p)
}

// Value реализует интерфейс driver.Valuer для Permissions
func (p Permissions) Value() (driver.Value, error) {
	if p.Access == nil {
		p.Access = []string{}
	}
	return json.Marshal(p)
}

// StaffFilter представляет параметры фильтрации для списка сотрудников
type StaffFilter struct {
	Page       int
	PageSize   int
	SearchTerm string
	RoleID     int
}
