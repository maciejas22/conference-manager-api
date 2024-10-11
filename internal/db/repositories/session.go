package repositories

import (
	"errors"

	"github.com/jmoiron/sqlx"
)

type Session struct {
	SessionId      string `json:"session_id" db:"session_id"`
	UserId         int    `json:"user_id" db:"user_id"`
	CreatedAt      string `json:"created_at" db:"created_at"`
	ExpiresAt      string `json:"expires_at" db:"expires_at"`
	LastAccessedAt string `json:"last_accessed_at" db:"last_accessed_at"`
}

func (s *Session) TableName() string {
	return "public.sessions"
}

func GetSession(tx *sqlx.Tx, sessionId string) (*Session, error) {
	var session Session
	s := &Session{}
	query := `
    SELECT session_id, user_id, created_at, expires_at, last_accessed_at 
    FROM ` + s.TableName() + ` 
    WHERE session_id = $1 
  `

	err := tx.Get(&session, query, sessionId)
	if err != nil {
		return nil, errors.New("could not get session")
	}
	return &session, nil
}

func CreateSession(tx *sqlx.Tx, sessionId string, userId int) (*string, error) {
	s := &Session{}
	query := `
		INSERT INTO ` + s.TableName() + ` (session_id, user_id, expires_at)
		VALUES ($1, $2, NOW() + INTERVAL '1 hour')
		ON CONFLICT(user_id) DO UPDATE SET
			session_id = excluded.session_id,
			expires_at = NOW() + INTERVAL '1 hour'
	`
	_, err := tx.Exec(
		query,
		sessionId,
		userId,
	)
	if err != nil {
		return nil, errors.New("could not create or update session")
	}
	return &sessionId, nil
}
