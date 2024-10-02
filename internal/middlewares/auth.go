package middlewares

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/maciejas22/conference-manager/api/internal/auth"
	"github.com/maciejas22/conference-manager/api/internal/db"
	"github.com/maciejas22/conference-manager/api/internal/db/repositories"
)

const (
	SessionCookieKey = "session"
)

func getSessionIdFromHeader(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("no session id in header")
	}

	sessionId := strings.TrimPrefix(authHeader, "Bearer ")
	if sessionId == "" {
		return "", errors.New("invalid session id format")
	}

	return sessionId, nil
}

func verifySession(dbClient *db.DB, sessionId string) (*repositories.Session, error) {
	var session *repositories.Session
	err := db.Transaction(context.Background(), dbClient.Conn, func(tx *sqlx.Tx) error {
		var err error
		session, err = repositories.GetSession(tx, sessionId)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, errors.New("could not get session")
	}

	sessionExpiresAt, err := time.Parse(time.RFC3339, session.ExpiresAt)
	if err != nil {
		return nil, errors.New("could not parse session expiration time")
	}

	if time.Now().After(sessionExpiresAt) {
		return nil, errors.New("session expired")
	}

	return session, nil
}

func AuthMiddleware(dbClient *db.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sessionInfo := &auth.SessionInfo{}
			ctx := context.WithValue(r.Context(), auth.SessionInfoKey, sessionInfo)

			sessionId, err := getSessionIdFromHeader(r)
			if err != nil {
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
			session, err := verifySession(dbClient, sessionId)
			if err != nil {
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
			sessionInfo.SessionId = sessionId
			sessionInfo.ExpiresAt = session.ExpiresAt

			err = db.Transaction(context.Background(), dbClient.Conn, func(tx *sqlx.Tx) error {
				user, err := repositories.GetUserBySessionId(tx, sessionId)
				if err != nil {
					return err
				}

				sessionInfo.UserId = user.Id
				sessionInfo.Role = user.Role

				return nil
			})
			if err != nil {
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
