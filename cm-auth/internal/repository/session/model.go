package repository

import "time"

type Session struct {
	SessionId      string    `json:"session_id" db:"session_id"`
	UserId         int       `json:"user_id" db:"user_id"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	ExpiresAt      time.Time `json:"expires_at" db:"expires_at"`
	LastAccessedAt time.Time `json:"last_accessed_at" db:"last_accessed_at"`
}

func (s *Session) TableName() string {
	return "public.sessions"
}
