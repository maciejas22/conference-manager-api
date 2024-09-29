package auth

import (
	"context"

	"github.com/maciejas22/conference-manager/api/db/repositories"
)

type SessionInfo struct {
	SessionId string
	ExpiresAt string
	UserId    int
	Role      repositories.Role
}

type contextKey string

const SessionInfoKey contextKey = "sessionInfo"

func GetSessionInfo(ctx context.Context) *SessionInfo {
	if sessionInfo, ok := ctx.Value(SessionInfoKey).(*SessionInfo); ok {
		return sessionInfo
	}
	return nil
}
