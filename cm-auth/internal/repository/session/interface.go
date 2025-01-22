package repository

type SessionRepoInterface interface {
	GetSessionBySessionId(sessionId string) (*Session, error)
	GetSessionByUserId(userId int) (*Session, error)
	CreateSession(userId int) (*string, error)
	DestroySession(sessionId string) error
}
