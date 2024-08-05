package auth

import (
	"context"
)

var userCtxKey = ctxKey{"user"}

type ctxKey struct {
	keyName string
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
