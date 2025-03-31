package model

import (
	"time"

	"github.com/google/uuid"
)

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
