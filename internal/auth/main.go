package auth

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/maciejas22/conference-manager/api/models"
)

var userCtxKey = ctxKey{"user"}

type ctxKey struct {
	keyName string
}

// var signingKey = []byte("8bea9019597c24e075686b3005517707121e8f7bc998c001b196b3a80d97540d")
var signingKey = []byte("D7WYxga4NrSx+zbX8naY/18NMYCqLopsfh7V4KEwJPpg8rDw5untw3aNv9W99p/AIIcZrNPcFJLvreKg0HMpFw==")

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

type AmrMethod struct {
	Method    string `json:"method"`
	Timestamp int64  `json:"timestamp"`
}

func Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				next.ServeHTTP(w, r)
				return
			}
			splitToken := strings.Split(authHeader, "Bearer ")
			if len(splitToken) != 2 {
				next.ServeHTTP(w, r)
				return
			}
			tokenString := splitToken[1]
			log.Println(tokenString)
			token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
				return signingKey, nil
			})
			if err != nil || !token.Valid {
				next.ServeHTTP(w, r)
				return
			}
			claims, ok := token.Claims.(*TokenClaims)
			if !ok {
				next.ServeHTTP(w, r)
				return
			}
			ctx := NewContext(r.Context(), claims)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

func NewContext(ctx context.Context, claims *TokenClaims) context.Context {
	return context.WithValue(ctx, userCtxKey, claims)
}

func FromContext(ctx context.Context) (*TokenClaims, bool) {
	claims, ok := ctx.Value(userCtxKey).(*TokenClaims)
	if !ok {
		return &TokenClaims{}, false
	}

	return claims, true
}

func GenerateToken(userID string, role models.Role) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, TokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		Email:       "example@example.com",
		Role:        role,
		Aal:         "aal1",
		Amr:         []AmrMethod{{Method: "password", Timestamp: time.Now().Unix()}},
		SessionID:   "session_id_example",
		IsAnonymous: false,
		AppMetadata: AppMetadata{
			Provider:  "email",
			Providers: []string{"email"},
		},
		UserMetadata: UserMetadata{
			Email:         "example@example.com",
			EmailVerified: false,
			PhoneVerified: false,
			Role:          "Organizer",
			Sub:           userID,
			Username:      "username_example",
		},
	})

	s, err := t.SignedString(signingKey)
	if err != nil {
		return "", err
	}

	return s, nil
}
