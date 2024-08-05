package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/maciejas22/conference-manager/api/internal/models"
)

type AmrMethod struct {
	Method    string `json:"method"`
	Timestamp int64  `json:"timestamp"`
}

type AppMetadata struct {
	Provider  string   `json:"provider"`
	Providers []string `json:"providers"`
}

type UserMetadata struct {
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	PhoneVerified bool   `json:"phone_verified"`
	Role          string `json:"role"`
	Sub           string `json:"sub"`
	Username      string `json:"username"`
}

type TokenClaims struct {
	jwt.RegisteredClaims
	Email        string       `json:"email"`
	Role         models.Role  `json:"role"`
	Aal          string       `json:"aal"`
	Amr          []AmrMethod  `json:"amr"`
	SessionID    string       `json:"session_id"`
	IsAnonymous  bool         `json:"is_anonymous"`
	AppMetadata  AppMetadata  `json:"app_metadata"`
	UserMetadata UserMetadata `json:"user_metadata"`
}
