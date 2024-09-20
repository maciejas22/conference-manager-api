package middlewares

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/maciejas22/conference-manager/api/db"
	"github.com/maciejas22/conference-manager/api/db/repositories"
	"github.com/maciejas22/conference-manager/api/internal/auth"
)

type authResponseWriter struct {
	http.ResponseWriter
	sessionInfo *auth.SessionInfo
}

const (
	SessionCookieKey = "session"
)

func (w *authResponseWriter) Write(b []byte) (int, error) {
	if w.sessionInfo.SessionId == "" {
		return w.ResponseWriter.Write(b)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookieKey,
		Value:    w.sessionInfo.SessionId,
		HttpOnly: true,
		Expires:  time.Now().Add(1 * time.Hour),
		SameSite: http.SameSiteLaxMode,
	})
	return w.ResponseWriter.Write(b)
}

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

func verifySession(ctx context.Context, dbClient *db.DB, sessionId string) error {
	err := db.Transaction(ctx, dbClient.QueryExecutor, func(qe *db.QueryExecutor) error {
		session, err := repositories.GetSession(qe, sessionId)
		if err != nil {
			return errors.New("could not get session")
		}

		sessionExpiresAt, err := time.Parse(time.RFC3339, session.ExpiresAt)
		if err != nil {
			return errors.New("could not parse session expiration time")
		}

		if time.Now().After(sessionExpiresAt) {
			return errors.New("session expired")
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func AuthMiddleware(dbClient *db.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sessionInfo := &auth.SessionInfo{}
			arw := &authResponseWriter{
				ResponseWriter: w,
				sessionInfo:    sessionInfo,
			}
			ctx := context.WithValue(r.Context(), auth.SessionInfoKey, sessionInfo)

			sessionId, err := getSessionIdFromHeader(r)
			if err != nil {
				next.ServeHTTP(arw, r.WithContext(ctx))
				return
			}
			err = verifySession(ctx, dbClient, sessionId)
			if err != nil {
				next.ServeHTTP(arw, r.WithContext(ctx))
				return
			}
			arw.sessionInfo.SessionId = sessionId

			err = db.Transaction(ctx, dbClient.QueryExecutor, func(qe *db.QueryExecutor) error {
				user, err := repositories.GetUserBySessionId(qe, sessionId)
				if err != nil {
					newSessionId, err := auth.GenerateSessionId()
					if err != nil {
						return err
					}

					userSession, err := repositories.CreateSession(qe, newSessionId, user.Id)
					if err != nil {
						return err
					}

					arw.sessionInfo.SessionId = *userSession
				} else {
					arw.sessionInfo.UserID = user.Id
					arw.sessionInfo.Role = user.Role
				}
				return err
			})

			if err != nil {
				next.ServeHTTP(arw, r.WithContext(ctx))
				return
			}

			next.ServeHTTP(arw, r.WithContext(ctx))
		})
	}
}
