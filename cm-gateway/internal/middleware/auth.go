package middlewares

import (
	"context"
	"errors"
	"net/http"
	"strings"

	authPb "github.com/maciejas22/conference-manager-api/cm-proto/auth"
)

const (
	SessionCookieKey = "session"
)

type SessionInfo struct {
	SessionId       string
	UserId          int32
	Role            authPb.Role
	StripeAccountId string
}

type contextKey string

const SessionInfoKey contextKey = "sessionInfo"

func getSessionIdFromHeader(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("No session id in header")
	}

	sessionId := strings.TrimPrefix(authHeader, "Bearer ")
	if sessionId == "" {
		return "", errors.New("Invalid session id format")
	}

	return sessionId, nil
}

func updateSession(ctx context.Context, authService authPb.AuthServiceClient, sessionId string) (*SessionInfo, error) {
	sOk, err := authService.ValidateSession(ctx, &authPb.ValidateSessionRequest{
		SessionId: sessionId,
	})
	if err != nil || !sOk.IsValid {
		return nil, errors.New("Invalid session")
	}

	u, err := authService.UserProfileBySession(ctx, &authPb.UserProfileBySessionRequest{
		SessionId: sessionId,
	})
	if err != nil {
		return nil, err
	}

	return &SessionInfo{
		SessionId:       sessionId,
		UserId:          u.User.UserId,
		Role:            u.User.Role,
		StripeAccountId: u.User.StripeAccountId,
	}, nil
}

func AuthMiddleware(authService authPb.AuthServiceClient) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sessionInfo := &SessionInfo{}
			ctx := context.WithValue(r.Context(), SessionInfoKey, sessionInfo)

			sessionId, err := getSessionIdFromHeader(r)
			if err != nil {
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
			sessionInfo.SessionId = sessionId

			s, err := updateSession(r.Context(), authService, sessionId)
			if err != nil {
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
			sessionInfo.UserId = s.UserId
			sessionInfo.Role = s.Role
			sessionInfo.StripeAccountId = s.StripeAccountId

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetSessionInfo(ctx context.Context) *SessionInfo {
	if sessionInfo, ok := ctx.Value(SessionInfoKey).(*SessionInfo); ok {
		return sessionInfo
	}
	return nil
}
