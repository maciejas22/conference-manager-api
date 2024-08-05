package middlewares

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/maciejas22/conference-manager/api/internal/auth"
	"github.com/maciejas22/conference-manager/api/internal/config"
)

func AuthMiddleware() func(http.Handler) http.Handler {
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
			token, err := jwt.ParseWithClaims(tokenString, &auth.TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
				return []byte(config.AppConfig.JWTSecret), nil
			})
			if err != nil || !token.Valid {
				next.ServeHTTP(w, r)
				return
			}
			claims, ok := token.Claims.(*auth.TokenClaims)
			if !ok {
				next.ServeHTTP(w, r)
				return
			}
			ctx := auth.NewContext(r.Context(), claims)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}
