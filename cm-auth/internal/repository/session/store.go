package repository

import (
	"errors"

	"github.com/jmoiron/sqlx"
)

type SessionRepo struct {
	Db *sqlx.DB
}

func NewSessionRepo(db *sqlx.DB) SessionRepoInterface {
	return &SessionRepo{Db: db}
}

func (r *SessionRepo) UpdateLastAccessedAt(sessionId string) error {
	query := `
    UPDATE ` + (new(Session)).TableName() + ` 
    SET last_accessed_at = NOW() 
    WHERE session_id = $1
  `
	_, err := r.Db.Exec(query, sessionId)
	if err != nil {
		return errors.New("Could not update last_accessed_at column")
	}
	return nil
}

func (r *SessionRepo) GetSessionBySessionId(sessionId string) (*Session, error) {
	query := `
    SELECT session_id, user_id, created_at, expires_at, last_accessed_at 
    FROM ` + (new(Session)).TableName() + ` 
    WHERE session_id = $1 
  `

	session := Session{}
	err := r.Db.Get(&session, query, sessionId)
	if err != nil {
		return nil, errors.New("Could not get session")
	}

	err = r.UpdateLastAccessedAt(sessionId)
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *SessionRepo) GetSessionByUserId(userId int) (*Session, error) {
	query := `
    SELECT session_id, user_id, created_at, expires_at, last_accessed_at
    FROM ` + (new(Session)).TableName() + `
    WHERE user_id = $1
  `

	session := Session{}
	err := r.Db.Get(&session, query, userId)
	if err != nil {
		return nil, errors.New("Could not get session")
	}

	err = r.UpdateLastAccessedAt(session.SessionId)
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *SessionRepo) CreateSession(userId int) (*string, error) {
	query := `
		INSERT INTO ` + (new(Session)).TableName() + ` (session_id, user_id, expires_at)
		VALUES (gen_random_uuid(), $1, NOW() + INTERVAL '1 hour')
		ON CONFLICT(user_id) DO UPDATE SET
			session_id = excluded.session_id,
			expires_at = NOW() + INTERVAL '1 hour'
    RETURNING session_id
	`
	var sessionId string
	err := r.Db.QueryRowx(query, userId).Scan(&sessionId)
	if err != nil {
		return nil, errors.New("Could not create or update session")
	}

	err = r.UpdateLastAccessedAt(sessionId)
	if err != nil {
		return nil, err
	}
	return &sessionId, nil
}

func (r *SessionRepo) DestroySession(sessionId string) error {
	query := `
    DELETE FROM ` + (new(Session)).TableName() + `
    WHERE session_id = $1
  `
	_, err := r.Db.Exec(query, sessionId)
	if err != nil {
		return errors.New("Could not destroy session")
	}

	err = r.UpdateLastAccessedAt(sessionId)
	if err != nil {
		return err
	}
	return nil
}
