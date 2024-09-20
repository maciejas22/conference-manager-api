package repositories

import (
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/maciejas22/conference-manager/api/db"
)

type Session struct {
	SessionId      string `json:"session_id" db:"session_id"`
	UserId         int    `json:"user_id" db:"user_id"`
	CreatedAt      string `json:"created_at" db:"created_at"`
	ExpiresAt      string `json:"expires_at" db:"expires_at"`
	LastAccessedAt string `json:"last_accessed_at" db:"last_accessed_at"`
}

func (s *Session) TableName() string {
	return "sessions"
}

func GetSession(qe *db.QueryExecutor, sessionId string) (*Session, error) {
	var session Session
	s := &Session{}
	query := `
    SELECT session_id, user_id, created_at, expires_at, last_accessed_at 
    FROM ` + s.TableName() + ` 
    WHERE session_id = ?
  `

	err := sqlx.Get(qe, &session, query, sessionId)
	if err != nil {
		return nil, errors.New("could not get session")
	}
	return &session, nil
}

func GetUserSession(qe *db.QueryExecutor, userId int) (*Session, error) {
	var session Session
	s := &Session{}
	query := `
    SELECT session_id, user_id, created_at, expires_at, last_accessed_at 
    FROM ` + s.TableName() + ` 
    WHERE user_id = ?
  `
	err := sqlx.Get(
		qe,
		&session,
		query,
		userId,
	)
	if err != nil {
		return nil, errors.New("could not get session")
	}
	return &session, nil
}

func CreateSession(qe *db.QueryExecutor, sessionId string, userId int) (*string, error) {
	s := &Session{}
	query := `
		INSERT INTO ` + s.TableName() + ` (session_id, user_id, expires_at)
		VALUES (?, ?, datetime('now', '+1 hour'))
		ON CONFLICT(user_id) DO UPDATE SET
			session_id = excluded.session_id,
      expires_at = datetime('now', '+1 hour')
	`
	_, err := qe.Exec(
		query,
		sessionId,
		userId,
	)
	if err != nil {
		return nil, errors.New("could not create or update session")
	}
	return &sessionId, nil
}
